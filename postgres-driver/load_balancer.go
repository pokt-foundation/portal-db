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
	ErrInvalidUsersJSON = errors.New("error: users JSON is invalid")
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

/* WriteLoadBalancerUser saves input LoadBalancer to the database */
func (p *PostgresDriver) WriteLoadBalancerUser(ctx context.Context, insert types.InsertUserAccess) error {
	if insert.ID == "" {
		return ErrMissingID
	}

	time := time.Now()
	params := InsertUserAccessParams{
		LbID:      newSQLNullString(insert.ID),
		UserID:    newSQLNullString(insert.UserID),
		RoleName:  newSQLNullString(insert.RoleName),
		Email:     newSQLNullString(insert.Email),
		CreatedAt: newSQLNullTime(time),
		UpdatedAt: newSQLNullTime(time),
	}

	err := p.InsertUserAccess(ctx, params)
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

// UpdateLoadBalancer updates fields available in options in db
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
