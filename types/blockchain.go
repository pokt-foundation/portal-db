package types

import (
	"time"
)

type (
	Blockchain struct {
		ID                string           `json:"id"`
		Altruist          string           `json:"altruist"`
		Blockchain        string           `json:"blockchain"`
		ChainID           string           `json:"chainID"`
		ChainIDCheck      string           `json:"chainIDCheck"`
		Description       string           `json:"description"`
		EnforceResult     string           `json:"enforceResult"`
		Network           string           `json:"network"`
		Path              string           `json:"path"`
		SyncCheck         string           `json:"syncCheck"`
		Ticker            string           `json:"ticker"`
		BlockchainAliases []string         `json:"blockchainAliases"`
		LogLimitBlocks    int              `json:"logLimitBlocks"`
		RequestTimeout    int              `json:"requestTimeout"`
		SyncAllowance     int              `json:"syncAllowance"`
		Active            bool             `json:"active"`
		Redirects         []Redirect       `json:"redirects"`
		SyncCheckOptions  SyncCheckOptions `json:"syncCheckOptions"`
		CreatedAt         time.Time        `json:"createdAt"`
		UpdatedAt         time.Time        `json:"updatedAt"`
	}
	Redirect struct {
		BlockchainID   string    `json:"blockchainID"`
		Alias          string    `json:"alias"`
		Domain         string    `json:"domain"`
		LoadBalancerID string    `json:"loadBalancerID"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}
	SyncCheckOptions struct {
		BlockchainID string `json:"blockchainID"`
		Body         string `json:"body"`
		Path         string `json:"path"`
		ResultKey    string `json:"resultKey"`
		Allowance    int    `json:"allowance"`
	}
	/* Update structs */
	UpdateBlockchain struct {
		BlockchainID      string   `json:"blockchainID"`
		Altruist          string   `json:"altruist"`
		Blockchain        string   `json:"blockchain"`
		BlockchainAliases []string `json:"blockchainAliases"`
		ChainID           string   `json:"chainID"`
		ChainIDCheck      string   `json:"chainIDCheck"`
		Description       string   `json:"description"`
		EnforceResult     string   `json:"enforceResult"`
		LogLimitBlocks    int32    `json:"logLimitBlocks"`
		Network           string   `json:"network"`
		Path              string   `json:"path"`
		RequestTimeout    int32    `json:"requestTimeout"`
		Ticker            string   `json:"ticker"`

		Synccheck     string `json:"synccheck"`
		Allowance     *int32 `json:"allowance"`
		Body          string `json:"body"`
		SyncCheckPath string `json:"sync_check_path"`
		ResultKey     string `json:"resultKey"`

		UpdatedAt time.Time `json:"updatedAt"`
	}
)
