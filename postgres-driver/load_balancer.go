package postgresdriver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pokt-foundation/portal-db/types"
)

var (
	ErrInvalidUsersJSON        = errors.New("error: users JSON is invalid")
	ErrUserInputIsMissingField = errors.New("error: user access input is missing a required field")
	ErrLBMustHaveUser          = errors.New("error: a new load balancer must have at least one user")
	ErrCannotSetToOwner        = errors.New("error: load balancers may only have one owner and the owner role is already set")
)

/* ReadLoadBalancers returns all LoadBalancers in the database */
func (p *PostgresDriver) ReadLoadBalancers(ctx context.Context) ([]*types.LoadBalancer, error) {
	dbLoadBalancers, err := p.SelectLoadBalancers(ctx)
	if err != nil {
		return nil, err
	}

	var loadbalancers []*types.LoadBalancer
	for _, dbLoadBalancer := range dbLoadBalancers {
		loadBalancer, err := dbLoadBalancer.toLoadBalancer()
		if err != nil {
			return nil, err
		}

		loadbalancers = append(loadbalancers, loadBalancer)
	}

	return loadbalancers, nil
}

func (lb *SelectLoadBalancersRow) toLoadBalancer() (*types.LoadBalancer, error) {
	loadBalancer := types.LoadBalancer{
		ID:                lb.LbID,
		Name:              lb.Name.String,
		UserID:            lb.UserID.String,
		ApplicationIDs:    strings.Split(string(lb.AppIds), ","),
		RequestTimeout:    int(lb.RequestTimeout.Int32),
		Gigastake:         lb.Gigastake.Bool,
		GigastakeRedirect: lb.GigastakeRedirect.Bool,

		StickyOptions: types.StickyOptions{
			Duration:      lb.SDuration.String,
			StickyOrigins: lb.SOrigins,
			StickyMax:     int(lb.SStickyMax.Int32),
			Stickiness:    lb.SStickiness.Bool,
		},

		CreatedAt: lb.CreatedAt.Time,
		UpdatedAt: lb.UpdatedAt.Time,
	}

	// Unmarshal LoadBalancer Users JSON into []types.UserAccess
	err := json.Unmarshal(lb.Users, &loadBalancer.Users)
	if err != nil {
		return &types.LoadBalancer{}, fmt.Errorf("%w: %s", ErrInvalidUsersJSON, err)
	}

	return &loadBalancer, nil
}

/* WriteLoadBalancer saves input LoadBalancer to the database */
func (p *PostgresDriver) WriteLoadBalancer(ctx context.Context, loadBalancer *types.LoadBalancer) (*types.LoadBalancer, error) {
	if len(loadBalancer.Users) < 1 {
		return nil, ErrLBMustHaveUser
	}

	id, err := generateRandomID()
	if err != nil {
		return nil, err
	}
	loadBalancer.ID = id
	time := time.Now()
	loadBalancer.CreatedAt = time
	loadBalancer.UpdatedAt = time

	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := p.WithTx(tx)

	err = qtx.InsertLoadBalancer(ctx, extractInsertLoadBalancer(loadBalancer))
	if err != nil {
		return nil, err
	}

	stickinessParams := extractInsertStickinessOptions(loadBalancer)
	if stickinessParams.isNotNull() {
		err = qtx.InsertStickinessOptions(ctx, stickinessParams)
		if err != nil {
			return nil, err
		}
	}

	// The first User will be the initial creater (owner) of the LoadBalancer
	loadBalancer.Users[0].RoleName = types.RoleOwner
	userAccessParams := extractInsertUserAccess(id, loadBalancer.Users[0], time)
	if userAccessParams.isNotNull() {
		err = qtx.InsertUserAccess(ctx, userAccessParams)
		if err != nil {
			return nil, err
		}
	}

	lbAppParams := InsertLbAppsParams{LbID: loadBalancer.ID}
	lbAppParams.AppIds = append(lbAppParams.AppIds, loadBalancer.ApplicationIDs...)

	err = qtx.InsertLbApps(ctx, lbAppParams)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func extractInsertLoadBalancer(loadBalancer *types.LoadBalancer) InsertLoadBalancerParams {
	return InsertLoadBalancerParams{
		LbID:              loadBalancer.ID,
		Name:              newSQLNullString(loadBalancer.Name),
		UserID:            newSQLNullString(loadBalancer.UserID),
		RequestTimeout:    newSQLNullInt32(int32(loadBalancer.RequestTimeout), false),
		Gigastake:         newSQLNullBool(&loadBalancer.Gigastake),
		GigastakeRedirect: newSQLNullBool(&loadBalancer.GigastakeRedirect),
		CreatedAt:         newSQLNullTime(loadBalancer.CreatedAt),
		UpdatedAt:         newSQLNullTime(loadBalancer.UpdatedAt),
	}
}

func extractInsertStickinessOptions(loadBalancer *types.LoadBalancer) InsertStickinessOptionsParams {
	return InsertStickinessOptionsParams{
		LbID:       loadBalancer.ID,
		Duration:   newSQLNullString(loadBalancer.StickyOptions.Duration),
		Origins:    loadBalancer.StickyOptions.StickyOrigins,
		StickyMax:  newSQLNullInt32(int32(loadBalancer.StickyOptions.StickyMax), false),
		Stickiness: newSQLNullBool(&loadBalancer.StickyOptions.Stickiness),
	}
}
func (i *InsertStickinessOptionsParams) isNotNull() bool {
	return i.Duration.Valid || len(i.Origins) > 0 || i.StickyMax.Valid
}

func extractInsertUserAccess(lbID string, userAccess types.UserAccess, createdAt time.Time) InsertUserAccessParams {
	return InsertUserAccessParams{
		LbID:      newSQLNullString(lbID),
		UserID:    newSQLNullString(userAccess.UserID),
		RoleName:  newSQLNullString(string(userAccess.RoleName)),
		Email:     newSQLNullString(userAccess.Email),
		CreatedAt: newSQLNullTime(createdAt),
		UpdatedAt: newSQLNullTime(createdAt),
	}
}
func (i *InsertUserAccessParams) isNotNull() bool {
	return i.LbID.Valid || i.UserID.Valid || i.RoleName.Valid || i.Email.Valid
}
func (i *InsertUserAccessParams) checkForMissingField() string {
	if !i.UserID.Valid {
		return "UserID"
	}
	if !i.RoleName.Valid {
		return "RoleName"
	}
	if !i.Email.Valid {
		return "Email"
	}
	return ""
}

/* WriteLoadBalancerUser saves input LoadBalancer to the database */
func (p *PostgresDriver) WriteLoadBalancerUser(ctx context.Context, lbID string, userAccess types.UserAccess) error {
	if lbID == "" {
		return ErrMissingID
	}
	if userAccess.RoleName == types.RoleOwner {
		return ErrCannotSetToOwner
	}

	userAccessParams := extractInsertUserAccess(lbID, userAccess, time.Now())

	missingField := userAccessParams.checkForMissingField()
	if missingField != "" {
		return fmt.Errorf("%w: %s", ErrUserInputIsMissingField, missingField)
	}

	err := p.InsertUserAccess(ctx, userAccessParams)
	if err != nil {
		return err
	}

	return nil
}

/* UpdateLoadBalancer updates LoadBalancer and related table rows */
func (p *PostgresDriver) UpdateLoadBalancer(ctx context.Context, id string, update *types.UpdateLoadBalancer) error {
	if id == "" {
		return ErrMissingID
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := p.WithTx(tx)

	err = qtx.UpdateLB(ctx, extractUpsertLoadBalancer(id, update))
	if err != nil {
		return err
	}

	stickinessOptionsParams := extractUpsertStickinessOptions(id, update)
	if stickinessOptionsParams.isNotNull() {
		err = qtx.UpsertStickinessOptions(ctx, *stickinessOptionsParams)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func extractUpsertLoadBalancer(id string, update *types.UpdateLoadBalancer) UpdateLBParams {
	return UpdateLBParams{
		LbID:      id,
		Name:      newSQLNullString(update.Name),
		UpdatedAt: newSQLNullTime(time.Now()),
	}
}

func extractUpsertStickinessOptions(id string, update *types.UpdateLoadBalancer) *UpsertStickinessOptionsParams {
	if update.StickyOptions == nil {
		return nil
	}

	return &UpsertStickinessOptionsParams{
		LbID:       id,
		Duration:   newSQLNullString(update.StickyOptions.Duration),
		StickyMax:  newSQLNullInt32(int32(update.StickyOptions.StickyMax), false),
		Stickiness: newSQLNullBool(update.StickyOptions.Stickiness),
		Origins:    update.StickyOptions.StickyOrigins,
	}
}
func (u *UpsertStickinessOptionsParams) isNotNull() bool {
	return u != nil && (u.Duration.Valid || u.StickyMax.Valid || u.Stickiness.Valid || len(u.Origins) != 0)
}

/* UpdateUserAccessRole updates the RoleName for a UserAccess row */
func (p *PostgresDriver) UpdateUserAccessRole(ctx context.Context, userID, lbID string, roleName types.RoleName) error {
	if userID == "" || lbID == "" {
		return ErrMissingID
	}
	if roleName == types.RoleOwner {
		return ErrCannotSetToOwner
	}

	params := UpdateUserAccessParams{
		UserID:    newSQLNullString(userID),
		LbID:      newSQLNullString(lbID),
		RoleName:  newSQLNullString(string(roleName)),
		UpdatedAt: newSQLNullTime(time.Now()),
	}

	err := p.UpdateUserAccess(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* RemoveLoadBalancer sets the user ID to an empty string (will not appear in Portal API or UI) */
func (p *PostgresDriver) RemoveLoadBalancer(ctx context.Context, id string) error {
	if id == "" {
		return ErrMissingID
	}

	err := p.RemoveLB(ctx, RemoveLBParams{LbID: id, UpdatedAt: newSQLNullTime(time.Now())})
	if err != nil {
		return err
	}

	return nil
}

/* RemoveUserAccess deletes a UserAccess row */
func (p *PostgresDriver) RemoveUserAccess(ctx context.Context, userID, lbID string) error {
	if userID == "" || lbID == "" {
		return ErrMissingID
	}

	params := DeleteUserAccessParams{
		UserID: newSQLNullString(userID),
		LbID:   newSQLNullString(lbID),
	}

	err := p.DeleteUserAccess(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

/* Used by Listener */
type (
	dbLoadBalancerJSON struct {
		LbID              string `json:"lb_id"`
		Name              string `json:"name"`
		UserID            string `json:"user_id"`
		RequestTimeout    int    `json:"request_timeout"`
		Gigastake         bool   `json:"gigastake"`
		GigastakeRedirect bool   `json:"gigastake_redirect"`
		CreatedAt         string `json:"created_at"`
		UpdatedAt         string `json:"updated_at"`
	}

	dbStickinessOptionsJSON struct {
		LbID       string   `json:"lb_id"`
		Duration   string   `json:"duration"`
		Origins    []string `json:"origins"`
		StickyMax  int      `json:"sticky_max"`
		Stickiness bool     `json:"stickiness"`
	}
)

func (j dbLoadBalancerJSON) toOutput() *types.LoadBalancer {
	return &types.LoadBalancer{
		ID:                j.LbID,
		Name:              j.Name,
		UserID:            j.UserID,
		RequestTimeout:    j.RequestTimeout,
		Gigastake:         j.Gigastake,
		GigastakeRedirect: j.GigastakeRedirect,
		CreatedAt:         psqlDateToTime(j.CreatedAt),
		UpdatedAt:         psqlDateToTime(j.UpdatedAt),
	}
}

func (j dbStickinessOptionsJSON) toOutput() *types.StickyOptions {
	return &types.StickyOptions{
		ID:            j.LbID,
		Duration:      j.Duration,
		StickyOrigins: j.Origins,
		StickyMax:     j.StickyMax,
		Stickiness:    j.Stickiness,
	}
}
