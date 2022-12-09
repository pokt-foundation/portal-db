// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: query.sql

package driver

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const activateBlockchain = `-- name: ActivateBlockchain :exec
UPDATE blockchains
SET active = $2
WHERE blockchain_id = $1
`

type ActivateBlockchainParams struct {
	BlockchainID string       `json:"blockchainID"`
	Active       sql.NullBool `json:"active"`
}

func (q *Queries) ActivateBlockchain(ctx context.Context, arg ActivateBlockchainParams) error {
	_, err := q.db.ExecContext(ctx, activateBlockchain, arg.BlockchainID, arg.Active)
	return err
}

const insertBlockchain = `-- name: InsertBlockchain :exec
INSERT into blockchains (
        blockchain_id,
        active,
        altruist,
        blockchain,
        blockchain_aliases,
        chain_id,
        chain_id_check,
        description,
        enforce_result,
        log_limit_blocks,
        network,
        path,
        request_timeout,
        ticker
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13,
        $14
    )
`

type InsertBlockchainParams struct {
	BlockchainID      string         `json:"blockchainID"`
	Active            sql.NullBool   `json:"active"`
	Altruist          sql.NullString `json:"altruist"`
	Blockchain        sql.NullString `json:"blockchain"`
	BlockchainAliases []string       `json:"blockchainAliases"`
	ChainID           sql.NullString `json:"chainID"`
	ChainIDCheck      sql.NullString `json:"chainIDCheck"`
	Description       sql.NullString `json:"description"`
	EnforceResult     sql.NullString `json:"enforceResult"`
	LogLimitBlocks    sql.NullInt32  `json:"logLimitBlocks"`
	Network           sql.NullString `json:"network"`
	Path              sql.NullString `json:"path"`
	RequestTimeout    sql.NullInt32  `json:"requestTimeout"`
	Ticker            sql.NullString `json:"ticker"`
}

func (q *Queries) InsertBlockchain(ctx context.Context, arg InsertBlockchainParams) error {
	_, err := q.db.ExecContext(ctx, insertBlockchain,
		arg.BlockchainID,
		arg.Active,
		arg.Altruist,
		arg.Blockchain,
		pq.Array(arg.BlockchainAliases),
		arg.ChainID,
		arg.ChainIDCheck,
		arg.Description,
		arg.EnforceResult,
		arg.LogLimitBlocks,
		arg.Network,
		arg.Path,
		arg.RequestTimeout,
		arg.Ticker,
	)
	return err
}

const insertSyncCheckOptions = `-- name: InsertSyncCheckOptions :exec
INSERT into sync_check_options (
        blockchain_id,
        synccheck,
        allowance,
        body,
        path,
        result_key
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
`

type InsertSyncCheckOptionsParams struct {
	BlockchainID string         `json:"blockchainID"`
	Synccheck    sql.NullString `json:"synccheck"`
	Allowance    sql.NullInt32  `json:"allowance"`
	Body         sql.NullString `json:"body"`
	Path         sql.NullString `json:"path"`
	ResultKey    sql.NullString `json:"resultKey"`
}

func (q *Queries) InsertSyncCheckOptions(ctx context.Context, arg InsertSyncCheckOptionsParams) error {
	_, err := q.db.ExecContext(ctx, insertSyncCheckOptions,
		arg.BlockchainID,
		arg.Synccheck,
		arg.Allowance,
		arg.Body,
		arg.Path,
		arg.ResultKey,
	)
	return err
}

const selectBlockchains = `-- name: SelectBlockchains :many
SELECT b.blockchain_id,
    b.altruist,
    b.blockchain,
    b.blockchain_aliases,
    b.chain_id,
    b.chain_id_check,
    b.description,
    b.enforce_result,
    b.log_limit_blocks,
    b.network,
    b.path,
    b.request_timeout,
    b.ticker,
    b.active,
    s.synccheck as s_sync_check,
    s.allowance as s_allowance,
    s.body as s_body,
    s.path as s_path,
    s.result_key as s_result_key
FROM blockchains as b
    LEFT JOIN sync_check_options AS s ON b.blockchain_id = s.blockchain_id
`

type SelectBlockchainsRow struct {
	BlockchainID      string         `json:"blockchainID"`
	Altruist          sql.NullString `json:"altruist"`
	Blockchain        sql.NullString `json:"blockchain"`
	BlockchainAliases []string       `json:"blockchainAliases"`
	ChainID           sql.NullString `json:"chainID"`
	ChainIDCheck      sql.NullString `json:"chainIDCheck"`
	Description       sql.NullString `json:"description"`
	EnforceResult     sql.NullString `json:"enforceResult"`
	LogLimitBlocks    sql.NullInt32  `json:"logLimitBlocks"`
	Network           sql.NullString `json:"network"`
	Path              sql.NullString `json:"path"`
	RequestTimeout    sql.NullInt32  `json:"requestTimeout"`
	Ticker            sql.NullString `json:"ticker"`
	Active            sql.NullBool   `json:"active"`
	SSyncCheck        sql.NullString `json:"sSyncCheck"`
	SAllowance        sql.NullInt32  `json:"sAllowance"`
	SBody             sql.NullString `json:"sBody"`
	SPath             sql.NullString `json:"sPath"`
	SResultKey        sql.NullString `json:"sResultKey"`
}

func (q *Queries) SelectBlockchains(ctx context.Context) ([]SelectBlockchainsRow, error) {
	rows, err := q.db.QueryContext(ctx, selectBlockchains)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectBlockchainsRow
	for rows.Next() {
		var i SelectBlockchainsRow
		if err := rows.Scan(
			&i.BlockchainID,
			&i.Altruist,
			&i.Blockchain,
			pq.Array(&i.BlockchainAliases),
			&i.ChainID,
			&i.ChainIDCheck,
			&i.Description,
			&i.EnforceResult,
			&i.LogLimitBlocks,
			&i.Network,
			&i.Path,
			&i.RequestTimeout,
			&i.Ticker,
			&i.Active,
			&i.SSyncCheck,
			&i.SAllowance,
			&i.SBody,
			&i.SPath,
			&i.SResultKey,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
