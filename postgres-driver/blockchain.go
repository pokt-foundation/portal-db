package postgresdriver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pokt-foundation/portal-db/repository"
)

var (
	ErrInvalidRedirectJSON = errors.New("error: redirect JSON is invalid")
)

/* ReadBlockchains returns all blockchains on the database and marshals to repository struct */
func (q *Queries) ReadBlockchains(ctx context.Context) ([]*repository.Blockchain, error) {
	dbBlockchains, err := q.SelectBlockchains(ctx)
	if err != nil {
		return nil, err
	}

	var blockchains []*repository.Blockchain

	for _, dbBlockchain := range dbBlockchains {
		blockchain, err := dbBlockchain.toBlockchain()
		if err != nil {
			return nil, err
		}

		blockchains = append(blockchains, blockchain)
	}

	return blockchains, nil
}

func (b *SelectBlockchainsRow) toBlockchain() (*repository.Blockchain, error) {
	blockchain := repository.Blockchain{
		ID:                b.BlockchainID,
		Altruist:          b.Altruist.String,
		Blockchain:        b.Blockchain.String,
		ChainID:           b.ChainID.String,
		ChainIDCheck:      b.ChainIDCheck.String,
		Description:       b.Description.String,
		EnforceResult:     b.EnforceResult.String,
		Network:           b.Network.String,
		Path:              b.Path.String,
		SyncCheck:         b.SSyncCheck.String,
		Ticker:            b.Ticker.String,
		BlockchainAliases: b.BlockchainAliases,
		LogLimitBlocks:    int(b.LogLimitBlocks.Int32),
		RequestTimeout:    int(b.RequestTimeout.Int32),
		Active:            b.Active.Bool,
		SyncCheckOptions: repository.SyncCheckOptions{
			Body:      b.SBody.String,
			ResultKey: b.SResultKey.String,
			Path:      b.Path.String,
			Allowance: int(b.SAllowance.Int32),
		},
	}

	// Unmarshal Blockchain Redirects JSON into []repository.Redirects
	err := json.Unmarshal(b.Redirects, &blockchain.Redirects)
	if err != nil {
		return &repository.Blockchain{}, fmt.Errorf("%w: %s", ErrInvalidRedirectJSON, err)
	}

	return &blockchain, nil
}

/* WriteBlockchain saves input Blockchain struct to the database */
func (q *Queries) WriteBlockchain(ctx context.Context, blockchain *repository.Blockchain) (*repository.Blockchain, error) {
	err := q.InsertBlockchain(ctx, extractInsertDBBlockchain(blockchain))
	if err != nil {
		return nil, err
	}

	synccheckOptionsParams := extractInsertSyncCheckOptions(blockchain)
	if synccheckOptionsParams.isNotNull() {
		err = q.InsertSyncCheckOptions(ctx, extractInsertSyncCheckOptions(blockchain))
		if err != nil {
			return nil, err
		}
	}

	return blockchain, nil
}

func extractInsertDBBlockchain(blockchain *repository.Blockchain) InsertBlockchainParams {
	return InsertBlockchainParams{
		BlockchainID:      blockchain.ID,
		Altruist:          newSQLNullString(blockchain.Altruist),
		Blockchain:        newSQLNullString(blockchain.Blockchain),
		ChainID:           newSQLNullString(blockchain.ChainID),
		ChainIDCheck:      newSQLNullString(blockchain.ChainIDCheck),
		Path:              newSQLNullString(blockchain.Path),
		Description:       newSQLNullString(blockchain.Description),
		EnforceResult:     newSQLNullString(blockchain.EnforceResult),
		Network:           newSQLNullString(blockchain.Network),
		Ticker:            newSQLNullString(blockchain.Ticker),
		BlockchainAliases: blockchain.BlockchainAliases,
		LogLimitBlocks:    newSQLNullInt32(int32(blockchain.LogLimitBlocks)),
		RequestTimeout:    newSQLNullInt32(int32(blockchain.RequestTimeout)),
		Active:            newSQLNullBool(blockchain.Active),
	}
}

func extractInsertSyncCheckOptions(blockchain *repository.Blockchain) InsertSyncCheckOptionsParams {
	return InsertSyncCheckOptionsParams{
		BlockchainID: blockchain.ID,
		Synccheck:    newSQLNullString(blockchain.SyncCheck),
		Body:         newSQLNullString(blockchain.SyncCheckOptions.Body),
		Path:         newSQLNullString(blockchain.SyncCheckOptions.Path),
		ResultKey:    newSQLNullString(blockchain.SyncCheckOptions.ResultKey),
		Allowance:    newSQLNullInt32(int32(blockchain.SyncCheckOptions.Allowance)),
	}
}

func (i *InsertSyncCheckOptionsParams) isNotNull() bool {
	return i.Synccheck.Valid || i.Body.Valid || i.Path.Valid || i.ResultKey.Valid || i.Allowance.Valid
}

/* WriteRedirect saves input Redirect struct to the database.
It must be called separately from WriteBlockchain due to how new chains are added to the dB */
func (q *Queries) WriteRedirect(ctx context.Context, redirect *repository.Redirect) (*repository.Redirect, error) {
	err := q.InsertRedirect(ctx, extractInsertDBRedirect(redirect))
	if err != nil {
		return nil, err
	}

	return redirect, nil
}

func extractInsertDBRedirect(redirect *repository.Redirect) InsertRedirectParams {
	return InsertRedirectParams{
		BlockchainID: redirect.BlockchainID,
		Alias:        redirect.Alias,
		Loadbalancer: redirect.LoadBalancerID,
		Domain:       redirect.Domain,
	}
}

/* Activate chain toggles chain.active field on or off */
func (q *Queries) ActivateChain(ctx context.Context, id string, active bool) error {
	params := ActivateBlockchainParams{BlockchainID: id, Active: newSQLNullBool(active)}

	err := q.ActivateBlockchain(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

// NOTE - temporaily commented out to satisfy linter
/* Listener Structs */
// type (
// 	dbBlockchainJSON struct {
// 		BlockchainID      string   `json:"blockchain_id"`
// 		Altruist          string   `json:"altruist"`
// 		Blockchain        string   `json:"blockchain"`
// 		ChainID           string   `json:"chain_id"`
// 		ChainIDCheck      string   `json:"chain_id_check"`
// 		ChainPath         string   `json:"path"`
// 		Description       string   `json:"description"`
// 		EnforceResult     string   `json:"enforce_result"`
// 		Network           string   `json:"network"`
// 		Ticker            string   `json:"ticker"`
// 		BlockchainAliases []string `json:"blockchain_aliases"`
// 		LogLimitBlocks    int      `json:"log_limit_blocks"`
// 		RequestTimeout    int      `json:"request_timeout"`
// 		Active            bool     `json:"active"`
// 		CreatedAt         string   `json:"created_at"`
// 		UpdatedAt         string   `json:"updated_at"`
// 	}
// 	dbSyncCheckOptionsJSON struct {
// 		BlockchainID string `json:"blockchain_id"`
// 		SyncCheck    string `json:"synccheck"`
// 		Body         string `json:"body"`
// 		Path         string `json:"path"`
// 		ResultKey    string `json:"result_key"`
// 		Allowance    int    `json:"allowance"`
// 	}
// )

// func (j dbBlockchainJSON) toOutput() *repository.Blockchain {
// 	return &repository.Blockchain{
// 		ID:                j.BlockchainID,
// 		Altruist:          j.Altruist,
// 		Blockchain:        j.Blockchain,
// 		ChainID:           j.ChainID,
// 		ChainIDCheck:      j.ChainIDCheck,
// 		Path:              j.ChainPath,
// 		Description:       j.Description,
// 		EnforceResult:     j.EnforceResult,
// 		Network:           j.Network,
// 		Ticker:            j.Ticker,
// 		BlockchainAliases: j.BlockchainAliases,
// 		LogLimitBlocks:    j.LogLimitBlocks,
// 		RequestTimeout:    j.RequestTimeout,
// 		Active:            j.Active,
// 		CreatedAt:         psqlDateToTime(j.CreatedAt),
// 		UpdatedAt:         psqlDateToTime(j.UpdatedAt),
// 	}
// }
// func (j dbSyncCheckOptionsJSON) toOutput() *repository.SyncCheckOptions {
// 	return &repository.SyncCheckOptions{
// 		BlockchainID: j.BlockchainID,
// 		Body:         j.Body,
// 		Path:         j.Path,
// 		ResultKey:    j.ResultKey,
// 		Allowance:    j.Allowance,
// 	}
// }
