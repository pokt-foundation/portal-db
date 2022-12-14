package postgresdriver

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/pokt-foundation/portal-db/types"
)

/* ReadApplications returns all Applications in the database */
func (p *PostgresDriver) ReadApplications(ctx context.Context) ([]*types.Application, error) {
	dbApplications, err := p.SelectApplications(ctx)
	if err != nil {
		return nil, err
	}

	var applications []*types.Application
	for _, dbApplication := range dbApplications {
		applications = append(applications, dbApplication.toApplication())
	}

	return applications, nil

}

func (a *SelectApplicationsRow) toApplication() *types.Application {
	return &types.Application{
		ID:                 a.ApplicationID,
		UserID:             a.UserID.String,
		Name:               a.Name.String,
		Status:             types.AppStatus(a.Status.String),
		ContactEmail:       a.ContactEmail.String,
		Description:        a.Description.String,
		Owner:              a.Owner.String,
		URL:                a.Url.String,
		Dummy:              a.Dummy.Bool,
		FirstDateSurpassed: a.FirstDateSurpassed.Time,

		GatewayAAT: types.GatewayAAT{
			Address:              a.GaAddress.String,
			ApplicationPublicKey: a.GaPublicKey.String,
			ApplicationSignature: a.GaSignature.String,
			ClientPublicKey:      a.GaClientPublicKey.String,
			PrivateKey:           a.GaPrivateKey.String,
			Version:              a.GaVersion.String,
		},
		GatewaySettings: types.GatewaySettings{
			SecretKey:            a.SecretKey.String,
			SecretKeyRequired:    a.SecretKeyRequired.Bool,
			WhitelistBlockchains: a.WhitelistBlockchains,
			WhitelistContracts:   nullStringToWhitelistContracts(a.WhitelistContracts),
			WhitelistMethods:     nullStringToWhitelistMethods(a.WhitelistMethods),
			WhitelistOrigins:     a.WhitelistOrigins,
			WhitelistUserAgents:  a.WhitelistUserAgents,
		},
		Limit: types.AppLimit{
			PayPlan: types.PayPlan{
				Type:  types.PayPlanType(a.PayPlan.String),
				Limit: int(a.PlanLimit.Int32),
			},
			CustomLimit: int(a.CustomLimit.Int32),
		},
		NotificationSettings: types.NotificationSettings{
			SignedUp:      a.SignedUp.Bool,
			Quarter:       a.OnQuarter.Bool,
			Half:          a.OnHalf.Bool,
			ThreeQuarters: a.OnThreeQuarters.Bool,
			Full:          a.OnFull.Bool,
		},

		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func nullStringToWhitelistContracts(rawContracts sql.NullString) []types.WhitelistContract {
	if !rawContracts.Valid {
		return nil
	}

	return stringToWhitelistContracts(rawContracts.String)
}
func stringToWhitelistContracts(rawContracts string) []types.WhitelistContract {
	contracts := []types.WhitelistContract{}

	_ = json.Unmarshal([]byte(rawContracts), &contracts)

	for i, contract := range contracts {
		for j, inContract := range contract.Contracts {
			contracts[i].Contracts[j] = strings.TrimSpace(inContract)
		}
	}

	return contracts
}

func nullStringToWhitelistMethods(rawMethods sql.NullString) []types.WhitelistMethod {
	if !rawMethods.Valid {
		return nil
	}

	return stringToWhitelistMethods(rawMethods.String)
}
func stringToWhitelistMethods(rawMethods string) []types.WhitelistMethod {
	methods := []types.WhitelistMethod{}

	_ = json.Unmarshal([]byte(rawMethods), &methods)

	for i, method := range methods {
		for j, inMethod := range method.Methods {
			methods[i].Methods[j] = strings.TrimSpace(inMethod)
		}
	}

	return methods
}

/* ReadPayPlans returns all pay plans in the database and marshals to types struct */
func (p *PostgresDriver) ReadPayPlans(ctx context.Context) ([]*types.PayPlan, error) {
	dbPayPlans, err := p.SelectPayPlans(ctx)
	if err != nil {
		return nil, err
	}

	var payPlans []*types.PayPlan

	for _, dbPayPlan := range dbPayPlans {
		payPlan, err := dbPayPlan.toPayPlan()
		if err != nil {
			return nil, err
		}

		payPlans = append(payPlans, payPlan)
	}

	return payPlans, nil
}

func (p *SelectPayPlansRow) toPayPlan() (*types.PayPlan, error) {
	payPlan := types.PayPlan{
		Type:  types.PayPlanType(p.PlanType),
		Limit: int(p.DailyLimit),
	}

	err := payPlan.Validate()
	if err != nil {
		return nil, err
	}

	return &payPlan, nil
}

/* WriteApplication saves input Application to the database */
func (p *PostgresDriver) WriteApplication(ctx context.Context, app *types.Application) (*types.Application, error) {
	appIsInvalid := app.Validate()
	if appIsInvalid != nil {
		return nil, appIsInvalid
	}

	id, err := generateRandomID()
	if err != nil {
		return nil, err
	}
	app.ID = id

	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := p.WithTx(tx)

	err = qtx.InsertApplication(ctx, extractInsertDBApp(app))
	if err != nil {
		return nil, err
	}

	err = qtx.InsertAppLimit(ctx, extractInsertDBAppLimit(app))
	if err != nil {
		return nil, err
	}
	gatewayAATParams := extractInsertDBGatewayAAT(app)
	if gatewayAATParams.isNotNull() {
		err = qtx.InsertGatewayAAT(ctx, gatewayAATParams)
		if err != nil {
			return nil, err
		}
	}
	gatewaySettingsParams := extractInsertDBGatewaySettings(app)
	if gatewaySettingsParams.isNotNull() {
		err = qtx.InsertGatewaySettings(ctx, gatewaySettingsParams)
		if err != nil {
			return nil, err
		}
	}
	notificationSettingsParams := extractInsertDBNotificationSettings(app)
	if notificationSettingsParams.isNotNull() {
		err = qtx.InsertNotificationSettings(ctx, notificationSettingsParams)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func extractInsertDBApp(app *types.Application) InsertApplicationParams {
	return InsertApplicationParams{
		ApplicationID: app.ID,
		UserID:        newSQLNullString(app.UserID),
		Name:          newSQLNullString(app.Name),
		ContactEmail:  newSQLNullString(app.ContactEmail),
		Description:   newSQLNullString(app.Description),
		Owner:         newSQLNullString(app.Owner),
		Url:           newSQLNullString(app.URL),
		Status:        newSQLNullString(string(app.Status)),
		Dummy:         newSQLNullBool(&app.Dummy),
	}
}

func extractInsertDBAppLimit(app *types.Application) InsertAppLimitParams {
	return InsertAppLimitParams{
		ApplicationID: app.ID,
		PayPlan:       string(app.Limit.PayPlan.Type),
		CustomLimit:   newSQLNullInt32(int32(app.Limit.CustomLimit), false),
	}
}

func extractInsertDBGatewayAAT(app *types.Application) InsertGatewayAATParams {
	return InsertGatewayAATParams{
		ApplicationID:   app.ID,
		Address:         app.GatewayAAT.Address,
		ClientPublicKey: app.GatewayAAT.ClientPublicKey,
		PrivateKey:      newSQLNullString(app.GatewayAAT.PrivateKey),
		PublicKey:       app.GatewayAAT.ApplicationPublicKey,
		Signature:       app.GatewayAAT.ApplicationSignature,
		Version:         newSQLNullString(app.GatewayAAT.Version),
	}
}
func (i *InsertGatewayAATParams) isNotNull() bool {
	return i.Version.Valid || i.PrivateKey.Valid
}

func extractInsertDBGatewaySettings(app *types.Application) InsertGatewaySettingsParams {
	marshaledWhitelistContracts, marshaledWhitelistMethods :=
		marshalWhitelistContractsAndMethods(app.GatewaySettings.WhitelistContracts, app.GatewaySettings.WhitelistMethods)

	return InsertGatewaySettingsParams{
		ApplicationID:        app.ID,
		SecretKey:            newSQLNullString(app.GatewaySettings.SecretKey),
		SecretKeyRequired:    newSQLNullBool(&app.GatewaySettings.SecretKeyRequired),
		WhitelistContracts:   newSQLNullString(marshaledWhitelistContracts),
		WhitelistMethods:     newSQLNullString(marshaledWhitelistMethods),
		WhitelistOrigins:     app.GatewaySettings.WhitelistOrigins,
		WhitelistUserAgents:  app.GatewaySettings.WhitelistUserAgents,
		WhitelistBlockchains: app.GatewaySettings.WhitelistBlockchains,
	}
}
func marshalWhitelistContractsAndMethods(contracts []types.WhitelistContract, methods []types.WhitelistMethod) (string, string) {
	var marshaledWhitelistContracts []byte
	if len(contracts) > 0 {
		marshaledWhitelistContracts, _ = json.Marshal(contracts)
	}

	var marshaledWhitelistMethods []byte
	if len(methods) > 0 {
		marshaledWhitelistMethods, _ = json.Marshal(methods)
	}

	return string(marshaledWhitelistContracts), string(marshaledWhitelistMethods)
}
func (i *InsertGatewaySettingsParams) isNotNull() bool {
	return i.SecretKey.Valid || i.WhitelistContracts.Valid || i.WhitelistMethods.Valid ||
		len(i.WhitelistOrigins) != 0 || len(i.WhitelistUserAgents) != 0 || len(i.WhitelistBlockchains) != 0
}

func extractInsertDBNotificationSettings(app *types.Application) InsertNotificationSettingsParams {
	return InsertNotificationSettingsParams{
		ApplicationID:   app.ID,
		SignedUp:        newSQLNullBool(&app.NotificationSettings.SignedUp),
		OnQuarter:       newSQLNullBool(&app.NotificationSettings.Quarter),
		OnHalf:          newSQLNullBool(&app.NotificationSettings.Half),
		OnThreeQuarters: newSQLNullBool(&app.NotificationSettings.ThreeQuarters),
		OnFull:          newSQLNullBool(&app.NotificationSettings.Full),
	}
}
func (i *InsertNotificationSettingsParams) isNotNull() bool {
	return true
}

/* UpdateApplication updates Application and related table rows */
func (p *PostgresDriver) UpdateApplication(ctx context.Context, id string, update *types.UpdateApplication) error {
	if id == "" {
		return ErrMissingID
	}

	invalidUpdate := update.Validate()
	if invalidUpdate != nil {
		return invalidUpdate
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := p.WithTx(tx)

	err = qtx.UpsertApplication(ctx, extractUpsertApplication(id, update))
	if err != nil {
		return err
	}
	err = qtx.UpsertAppLimit(ctx, extractUpsertAppLimit(id, update))
	if err != nil {
		return err
	}
	err = qtx.UpsertGatewaySettings(ctx, extractUpsertGatewaySettings(id, update))
	if err != nil {
		return err
	}
	err = qtx.UpsertNotificationSettings(ctx, extractUpsertNotificationSettings(id, update))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func extractUpsertApplication(id string, update *types.UpdateApplication) UpsertApplicationParams {
	return UpsertApplicationParams{
		ApplicationID:      id,
		Name:               newSQLNullString(update.Name),
		Status:             newSQLNullString(string(update.Status)),
		FirstDateSurpassed: newSQLNullTime(update.FirstDateSurpassed),
	}
}

func extractUpsertAppLimit(id string, update *types.UpdateApplication) UpsertAppLimitParams {
	customLimit := int32(update.Limit.CustomLimit)
	if update.Limit.PayPlan.Type != types.Enterprise {
		customLimit = 0
	}

	return UpsertAppLimitParams{
		ApplicationID: id,
		PayPlan:       string(update.Limit.PayPlan.Type),
		CustomLimit:   newSQLNullInt32(customLimit, true),
	}
}

func extractUpsertGatewaySettings(id string, update *types.UpdateApplication) UpsertGatewaySettingsParams {
	marshaledWhitelistContracts, marshaledWhitelistMethods :=
		marshalWhitelistContractsAndMethods(update.GatewaySettings.WhitelistContracts, update.GatewaySettings.WhitelistMethods)

	return UpsertGatewaySettingsParams{
		ApplicationID:        id,
		SecretKey:            newSQLNullString(update.GatewaySettings.SecretKey),
		SecretKeyRequired:    newSQLNullBool(update.GatewaySettings.SecretKeyRequired),
		WhitelistContracts:   newSQLNullString(marshaledWhitelistContracts),
		WhitelistMethods:     newSQLNullString(marshaledWhitelistMethods),
		WhitelistOrigins:     update.GatewaySettings.WhitelistOrigins,
		WhitelistUserAgents:  update.GatewaySettings.WhitelistUserAgents,
		WhitelistBlockchains: update.GatewaySettings.WhitelistBlockchains,
	}
}

func extractUpsertNotificationSettings(id string, update *types.UpdateApplication) UpsertNotificationSettingsParams {
	return UpsertNotificationSettingsParams{
		ApplicationID:   id,
		SignedUp:        newSQLNullBool(update.NotificationSettings.SignedUp),
		OnQuarter:       newSQLNullBool(update.NotificationSettings.Quarter),
		OnHalf:          newSQLNullBool(update.NotificationSettings.Half),
		OnThreeQuarters: newSQLNullBool(update.NotificationSettings.ThreeQuarters),
		OnFull:          newSQLNullBool(update.NotificationSettings.Full),
	}
}

/* UpdateAppFirstDateSurpassed updates Application's firstDateSurpassed field */
func (p *PostgresDriver) UpdateAppFirstDateSurpassed(ctx context.Context, update *types.UpdateFirstDateSurpassed) error {
	params := UpdateFirstDateSurpassedParams{
		ApplicationIds:     update.ApplicationIDs,
		FirstDateSurpassed: newSQLNullTime(update.FirstDateSurpassed),
	}

	err := p.UpdateFirstDateSurpassed(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* RemoveApplication updates Application's status field to AwaitingGracePeriod */
func (p *PostgresDriver) RemoveApplication(ctx context.Context, id string) error {
	if id == "" {
		return ErrMissingID
	}

	params := RemoveAppParams{
		ApplicationID: id,
		Status:        newSQLNullString(string(types.AwaitingGracePeriod)),
	}

	err := p.RemoveApp(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* Used by Listener */
type (
	dbAppJSON struct {
		ApplicationID      string `json:"application_id"`
		UserID             string `json:"user_id"`
		Name               string `json:"name"`
		ContactEmail       string `json:"contact_email"`
		Description        string `json:"description"`
		Owner              string `json:"owner"`
		URL                string `json:"url"`
		Status             string `json:"status"`
		CreatedAt          string `json:"created_at"`
		UpdatedAt          string `json:"updated_at"`
		FirstDateSurpassed string `json:"first_date_surpassed"`
		Dummy              bool   `json:"dummy"`
	}
	dbAppLimitJSON struct {
		ApplicationID string            `json:"application_id"`
		PlanType      types.PayPlanType `json:"pay_plan"`
		CustomLimit   int               `json:"custom_limit"`
	}
	dbGatewayAATJSON struct {
		ApplicationID   string `json:"application_id"`
		Address         string `json:"address"`
		ClientPublicKey string `json:"client_public_key"`
		PrivateKey      string `json:"private_key"`
		PublicKey       string `json:"public_key"`
		Signature       string `json:"signature"`
		Version         string `json:"version"`
	}
	dbGatewaySettingsJSON struct {
		ApplicationID        string   `json:"application_id"`
		SecretKey            string   `json:"secret_key"`
		SecretKeyRequired    bool     `json:"secret_key_required"`
		WhitelistContracts   string   `json:"whitelist_contracts"`
		WhitelistMethods     string   `json:"whitelist_methods"`
		WhitelistOrigins     []string `json:"whitelist_origins"`
		WhitelistUserAgents  []string `json:"whitelist_user_agents"`
		WhitelistBlockchains []string `json:"whitelist_blockchains"`
	}
	dbNotificationSettingsJSON struct {
		ApplicationID string `json:"application_id"`
		SignedUp      bool   `json:"signed_up"`
		Quarter       bool   `json:"on_quarter"`
		Half          bool   `json:"on_half"`
		ThreeQuarters bool   `json:"on_three_quarters"`
		Full          bool   `json:"on_full"`
	}
)

func (j dbAppJSON) toOutput() *types.Application {
	return &types.Application{
		ID:                 j.ApplicationID,
		UserID:             j.UserID,
		Name:               j.Name,
		ContactEmail:       j.ContactEmail,
		Description:        j.Description,
		Owner:              j.Owner,
		URL:                j.URL,
		Status:             types.AppStatus(j.Status),
		CreatedAt:          psqlDateToTime(j.CreatedAt),
		UpdatedAt:          psqlDateToTime(j.UpdatedAt),
		FirstDateSurpassed: psqlDateToTime(j.FirstDateSurpassed),
		Dummy:              j.Dummy,
	}
}
func (j dbAppLimitJSON) toOutput() *types.AppLimit {
	return &types.AppLimit{
		ID: j.ApplicationID,
		PayPlan: types.PayPlan{
			Type: j.PlanType,
		},
		CustomLimit: j.CustomLimit,
	}
}
func (j dbGatewayAATJSON) toOutput() *types.GatewayAAT {
	return &types.GatewayAAT{
		ID:                   j.ApplicationID,
		Address:              j.Address,
		ClientPublicKey:      j.ClientPublicKey,
		PrivateKey:           j.PrivateKey,
		ApplicationPublicKey: j.PublicKey,
		ApplicationSignature: j.Signature,
		Version:              j.Version,
	}
}
func (j dbGatewaySettingsJSON) toOutput() *types.GatewaySettings {
	return &types.GatewaySettings{
		ID:                   j.ApplicationID,
		SecretKey:            j.SecretKey,
		SecretKeyRequired:    j.SecretKeyRequired,
		WhitelistContracts:   stringToWhitelistContracts(j.WhitelistContracts),
		WhitelistMethods:     stringToWhitelistMethods(j.WhitelistMethods),
		WhitelistOrigins:     j.WhitelistOrigins,
		WhitelistUserAgents:  j.WhitelistUserAgents,
		WhitelistBlockchains: j.WhitelistBlockchains,
	}
}
func (j dbNotificationSettingsJSON) toOutput() *types.NotificationSettings {
	return &types.NotificationSettings{
		ID:            j.ApplicationID,
		SignedUp:      j.SignedUp,
		Quarter:       j.Quarter,
		Half:          j.Half,
		ThreeQuarters: j.ThreeQuarters,
		Full:          j.Full,
	}
}
