package postgresdriver

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/pokt-foundation/portal-db/repository"
)

/* ReadApplications returns all Applications in the database */
func (q *Queries) ReadApplications(ctx context.Context) ([]*repository.Application, error) {
	dbApplications, err := q.SelectApplications(ctx)
	if err != nil {
		return nil, err
	}

	var applications []*repository.Application

	for _, dbApplication := range dbApplications {
		applications = append(applications, dbApplication.toApplication())
	}

	return applications, nil

}

func (a *SelectApplicationsRow) toApplication() *repository.Application {
	return &repository.Application{
		ID:                 a.ApplicationID,
		UserID:             a.UserID.String,
		Name:               a.Name.String,
		Status:             repository.AppStatus(a.Status.String),
		ContactEmail:       a.ContactEmail.String,
		Description:        a.Description.String,
		Owner:              a.Owner.String,
		URL:                a.Url.String,
		Dummy:              a.Dummy.Bool,
		FirstDateSurpassed: a.FirstDateSurpassed.Time,

		GatewayAAT: repository.GatewayAAT{
			Address:              a.GaAddress.String,
			ApplicationPublicKey: a.GaPublicKey.String,
			ApplicationSignature: a.GaSignature.String,
			ClientPublicKey:      a.GaClientPublicKey.String,
			PrivateKey:           a.GaPrivateKey.String,
			Version:              a.GaVersion.String,
		},
		GatewaySettings: repository.GatewaySettings{
			SecretKey:            a.SecretKey.String,
			SecretKeyRequired:    a.SecretKeyRequired.Bool,
			WhitelistBlockchains: a.WhitelistBlockchains,
			WhitelistContracts:   nullStringToWhitelistContracts(a.WhitelistContracts),
			WhitelistMethods:     nullStringToWhitelistMethods(a.WhitelistMethods),
			WhitelistOrigins:     a.WhitelistOrigins,
			WhitelistUserAgents:  a.WhitelistUserAgents,
		},
		Limit: repository.AppLimit{
			PayPlan: repository.PayPlan{
				Type:  repository.PayPlanType(a.PayPlan.String),
				Limit: int(a.PlanLimit.Int32),
			},
			CustomLimit: int(a.CustomLimit.Int32),
		},
		NotificationSettings: repository.NotificationSettings{
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

func nullStringToWhitelistContracts(rawContracts sql.NullString) []repository.WhitelistContract {
	if !rawContracts.Valid {
		return nil
	}

	return stringToWhitelistContracts(rawContracts.String)
}
func stringToWhitelistContracts(rawContracts string) []repository.WhitelistContract {
	contracts := []repository.WhitelistContract{}

	_ = json.Unmarshal([]byte(rawContracts), &contracts)

	for i, contract := range contracts {
		for j, inContract := range contract.Contracts {
			contracts[i].Contracts[j] = strings.TrimSpace(inContract)
		}
	}

	return contracts
}

func nullStringToWhitelistMethods(rawMethods sql.NullString) []repository.WhitelistMethod {
	if !rawMethods.Valid {
		return nil
	}

	return stringToWhitelistMethods(rawMethods.String)
}
func stringToWhitelistMethods(rawMethods string) []repository.WhitelistMethod {
	methods := []repository.WhitelistMethod{}

	_ = json.Unmarshal([]byte(rawMethods), &methods)

	for i, method := range methods {
		for j, inMethod := range method.Methods {
			methods[i].Methods[j] = strings.TrimSpace(inMethod)
		}
	}

	return methods
}

/* WriteApplication saves input Application to the database */
func (q *Queries) WriteApplication(ctx context.Context, app *repository.Application) (*repository.Application, error) {
	appIsInvalid := app.Validate()
	if appIsInvalid != nil {
		return nil, appIsInvalid
	}

	id, err := generateRandomID()
	if err != nil {
		return nil, err
	}
	app.ID = id

	err = q.InsertApplication(ctx, extractInsertDBApp(app))
	if err != nil {
		return nil, err
	}

	err = q.InsertAppLimit(ctx, extractInsertDBAppLimit(app))
	if err != nil {
		return nil, err
	}
	gatewayAATParams := extractInsertGatewayAAT(app)
	if gatewayAATParams.isNotNull() {
		err = q.InsertGatewayAAT(ctx, gatewayAATParams)
		if err != nil {
			return nil, err
		}
	}
	gatewaySettingsParams := extractInsertGatewaySettings(app)
	if gatewaySettingsParams.isNotNull() {
		err = q.InsertGatewaySettings(ctx, gatewaySettingsParams)
		if err != nil {
			return nil, err
		}
	}
	notificationSettingsParams := extractInsertNotificationSettings(app)
	if notificationSettingsParams.isNotNull() {
		err = q.InsertNotificationSettings(ctx, notificationSettingsParams)
		if err != nil {
			return nil, err
		}
	}

	return app, nil
}

func extractInsertDBApp(app *repository.Application) InsertApplicationParams {
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

func extractInsertDBAppLimit(app *repository.Application) InsertAppLimitParams {
	return InsertAppLimitParams{
		ApplicationID: app.ID,
		PayPlan:       string(app.Limit.PayPlan.Type),
		CustomLimit:   newSQLNullInt32(int32(app.Limit.CustomLimit)),
	}
}

func extractInsertGatewayAAT(app *repository.Application) InsertGatewayAATParams {
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

func extractInsertGatewaySettings(app *repository.Application) InsertGatewaySettingsParams {
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
func marshalWhitelistContractsAndMethods(contracts []repository.WhitelistContract, methods []repository.WhitelistMethod) (string, string) {
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

func extractInsertNotificationSettings(app *repository.Application) InsertNotificationSettingsParams {
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
func (q *Queries) UpdateApplication(ctx context.Context, id string, update *repository.UpdateApplication) error {
	if id == "" {
		return ErrMissingID
	}

	invalidUpdate := update.Validate()
	if invalidUpdate != nil {
		return invalidUpdate
	}

	err := q.UpsertApplication(ctx, extractUpsertApplication(id, update))
	if err != nil {
		return err
	}
	err = q.UpsertAppLimit(ctx, extractUpsertAppLimit(id, update))
	if err != nil {
		return err
	}
	err = q.UpsertGatewaySettings(ctx, extractUpsertGatewaySettings(id, update))
	if err != nil {
		return err
	}
	err = q.UpsertNotificationSettings(ctx, extractUpsertNotificationSettings(id, update))
	if err != nil {
		return err
	}

	return nil
}

func extractUpsertApplication(id string, update *repository.UpdateApplication) UpsertApplicationParams {
	return UpsertApplicationParams{
		ApplicationID:      id,
		Name:               newSQLNullString(update.Name),
		Status:             newSQLNullString(string(update.Status)),
		FirstDateSurpassed: newSQLNullTime(update.FirstDateSurpassed),
	}
}

func extractUpsertAppLimit(id string, update *repository.UpdateApplication) UpsertAppLimitParams {
	return UpsertAppLimitParams{
		ApplicationID: id,
		PayPlan:       string(update.Limit.PayPlan.Type),
		CustomLimit:   newSQLNullInt32(int32(update.Limit.CustomLimit)),
	}
}

func extractUpsertGatewaySettings(id string, update *repository.UpdateApplication) UpsertGatewaySettingsParams {
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

func extractUpsertNotificationSettings(id string, update *repository.UpdateApplication) UpsertNotificationSettingsParams {
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
func (q *Queries) UpdateAppFirstDateSurpassed(ctx context.Context, update *repository.UpdateFirstDateSurpassed) error {
	params := UpdateFirstDateSurpassedParams{
		ApplicationIds:     update.ApplicationIDs,
		FirstDateSurpassed: newSQLNullTime(update.FirstDateSurpassed),
	}

	err := q.UpdateFirstDateSurpassed(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* RemoveApplication updates Application's status field to AwaitingGracePeriod */
func (q *Queries) RemoveApp(ctx context.Context, id string) error {
	if id == "" {
		return ErrMissingID
	}

	params := RemoveApplicationParams{
		ApplicationID: id,
		Status:        newSQLNullString(string(repository.AwaitingGracePeriod)),
	}

	err := q.RemoveApplication(ctx, params)
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
		ApplicationID string                 `json:"application_id"`
		PlanType      repository.PayPlanType `json:"pay_plan"`
		CustomLimit   int                    `json:"custom_limit"`
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

func (j dbAppJSON) toOutput() *repository.Application {
	return &repository.Application{
		ID:                 j.ApplicationID,
		UserID:             j.UserID,
		Name:               j.Name,
		ContactEmail:       j.ContactEmail,
		Description:        j.Description,
		Owner:              j.Owner,
		URL:                j.URL,
		Status:             repository.AppStatus(j.Status),
		CreatedAt:          psqlDateToTime(j.CreatedAt),
		UpdatedAt:          psqlDateToTime(j.UpdatedAt),
		FirstDateSurpassed: psqlDateToTime(j.FirstDateSurpassed),
		Dummy:              j.Dummy,
	}
}
func (j dbAppLimitJSON) toOutput() *repository.AppLimit {
	return &repository.AppLimit{
		ID: j.ApplicationID,
		PayPlan: repository.PayPlan{
			Type: j.PlanType,
		},
		CustomLimit: j.CustomLimit,
	}
}
func (j dbGatewayAATJSON) toOutput() *repository.GatewayAAT {
	return &repository.GatewayAAT{
		ID:                   j.ApplicationID,
		Address:              j.Address,
		ClientPublicKey:      j.ClientPublicKey,
		PrivateKey:           j.PrivateKey,
		ApplicationPublicKey: j.PublicKey,
		ApplicationSignature: j.Signature,
		Version:              j.Version,
	}
}
func (j dbGatewaySettingsJSON) toOutput() *repository.GatewaySettings {
	return &repository.GatewaySettings{
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
func (j dbNotificationSettingsJSON) toOutput() *repository.NotificationSettings {
	return &repository.NotificationSettings{
		ID:            j.ApplicationID,
		SignedUp:      j.SignedUp,
		Quarter:       j.Quarter,
		Half:          j.Half,
		ThreeQuarters: j.ThreeQuarters,
		Full:          j.Full,
	}
}
