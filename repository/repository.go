package repository

import (
	"errors"
	"time"
)

var (
	ErrNoFieldsToUpdate               = errors.New("no fields to update")
	ErrInvalidAppStatus               = errors.New("invalid app status")
	ErrInvalidPayPlanType             = errors.New("invalid pay plan type")
	ErrNotEnterprisePlan              = errors.New("custom limits may only be set on enterprise plans")
	ErrEnterprisePlanNeedsCustomLimit = errors.New("enterprise plans must have a custom limit set")
)

type (
	Table  string
	Action string

	Notification struct {
		Table  Table
		Action Action
		Data   SavedOnDB
	}
)

const (
	TableLoadBalancers        Table = "loadbalancers"
	TableStickinessOptions    Table = "stickiness_options"
	TableLbApps               Table = "lb_apps"
	TableApplications         Table = "applications"
	TableAppLimits            Table = "app_limits"
	TableGatewayAAT           Table = "gateway_aat"
	TableGatewaySettings      Table = "gateway_settings"
	TableNotificationSettings Table = "notification_settings"
	TableBlockchains          Table = "blockchains"
	TableRedirects            Table = "redirects"
	TableSyncCheckOptions     Table = "sync_check_options"

	ActionInsert Action = "INSERT"
	ActionUpdate Action = "UPDATE"
)

type SavedOnDB interface {
	Table() Table
}

/* Pay Plans Table */
type PayPlan struct {
	Type  PayPlanType `json:"planType"`
	Limit int         `json:"dailyLimit"`
}

func (p *PayPlan) Validate() error {
	if !ValidPayPlanTypes[p.Type] {
		return ErrInvalidPayPlanType
	}

	return nil
}

/* Applications Table */
type (
	Application struct {
		ID                   string               `json:"id"`
		UserID               string               `json:"userID"`
		Name                 string               `json:"name"`
		ContactEmail         string               `json:"contactEmail"`
		Description          string               `json:"description"`
		Owner                string               `json:"owner"`
		URL                  string               `json:"url"`
		Dummy                bool                 `json:"dummy"`
		Status               AppStatus            `json:"status"`
		FirstDateSurpassed   time.Time            `json:"firstDateSurpassed"`
		GatewayAAT           GatewayAAT           `json:"gatewayAAT"`
		GatewaySettings      GatewaySettings      `json:"gatewaySettings"`
		Limit                AppLimit             `json:"limit"`
		NotificationSettings NotificationSettings `json:"notificationSettings"`
		CreatedAt            time.Time            `json:"createdAt"`
		UpdatedAt            time.Time            `json:"updatedAt"`
	}
	GatewayAAT struct {
		ID                   string `json:"id,omitempty"`
		Address              string `json:"address"`
		ApplicationPublicKey string `json:"applicationPublicKey"`
		ApplicationSignature string `json:"applicationSignature"`
		ClientPublicKey      string `json:"clientPublicKey"`
		PrivateKey           string `json:"privateKey"`
		Version              string `json:"version"`
	}
	GatewaySettings struct {
		ID                   string              `json:"id,omitempty"`
		SecretKey            string              `json:"secretKey"`
		SecretKeyRequired    bool                `json:"secretKeyRequired"`
		WhitelistOrigins     []string            `json:"whitelistOrigins,omitempty"`
		WhitelistUserAgents  []string            `json:"whitelistUserAgents,omitempty"`
		WhitelistContracts   []WhitelistContract `json:"whitelistContracts,omitempty"`
		WhitelistMethods     []WhitelistMethod   `json:"whitelistMethods,omitempty"`
		WhitelistBlockchains []string            `json:"whitelistBlockchains,omitempty"`
	}
	WhitelistContract struct {
		BlockchainID string   `json:"blockchainID"`
		Contracts    []string `json:"contracts"`
	}
	WhitelistMethod struct {
		BlockchainID string   `json:"blockchainID"`
		Methods      []string `json:"methods"`
	}
	AppLimit struct {
		ID          string  `json:"id,omitempty"`
		PayPlan     PayPlan `json:"payPlan"`
		CustomLimit int     `json:"customLimit"`
	}
	NotificationSettings struct {
		ID            string `json:"id,omitempty"`
		SignedUp      bool   `json:"signedUp"`
		Quarter       bool   `json:"quarter"`
		Half          bool   `json:"half"`
		ThreeQuarters bool   `json:"threeQuarters"`
		Full          bool   `json:"full"`
	}

	UpdateApplication struct {
		Name                 string                     `json:"name,omitempty"`
		Status               AppStatus                  `json:"status,omitempty"`
		FirstDateSurpassed   time.Time                  `json:"firstDateSurpassed,omitempty"`
		GatewaySettings      UpdateGatewaySettings      `json:"gatewaySettings,omitempty"`
		NotificationSettings UpdateNotificationSettings `json:"notificationSettings,omitempty"`
		Limit                *AppLimit                  `json:"appLimit,omitempty"`
		Remove               bool                       `json:"remove,omitempty"`
	}
	UpdateGatewaySettings struct {
		ID                   string              `json:"id,omitempty"`
		SecretKey            string              `json:"secretKey"`
		SecretKeyRequired    *bool               `json:"secretKeyRequired"`
		WhitelistOrigins     []string            `json:"whitelistOrigins,omitempty"`
		WhitelistUserAgents  []string            `json:"whitelistUserAgents,omitempty"`
		WhitelistContracts   []WhitelistContract `json:"whitelistContracts,omitempty"`
		WhitelistMethods     []WhitelistMethod   `json:"whitelistMethods,omitempty"`
		WhitelistBlockchains []string            `json:"whitelistBlockchains,omitempty"`
	}
	UpdateFirstDateSurpassed struct {
		ApplicationIDs     []string  `json:"applicationIDs"`
		FirstDateSurpassed time.Time `json:"firstDateSurpassed"`
	}
	UpdateNotificationSettings struct {
		ID            string `json:"id,omitempty"`
		SignedUp      *bool  `json:"signedUp"`
		Quarter       *bool  `json:"quarter"`
		Half          *bool  `json:"half"`
		ThreeQuarters *bool  `json:"threeQuarters"`
		Full          *bool  `json:"full"`
	}

	AppStatus   string
	PayPlanType string
)

const (
	AwaitingFreetierFunds   AppStatus = "AWAITING_FREETIER_FUNDS"
	AwaitingFreetierStaking AppStatus = "AWAITING_FREETIER_STAKING"
	AwaitingFunds           AppStatus = "AWAITING_FUNDS"
	AwaitingFundsRemoval    AppStatus = "AWAITING_FUNDS_REMOVAL"
	AwaitingGracePeriod     AppStatus = "AWAITING_GRACE_PERIOD"
	AwaitingSlotFunds       AppStatus = "AWAITING_SLOT_FUNDS"
	AwaitingSlotStaking     AppStatus = "AWAITING_SLOT_STAKING"
	AwaitingStaking         AppStatus = "AWAITING_STAKING"
	AwaitingUnstaking       AppStatus = "AWAITING_UNSTAKING"
	Decomissioned           AppStatus = "DECOMISSIONED"
	InService               AppStatus = "IN_SERVICE"
	Orphaned                AppStatus = "ORPHANED"
	Ready                   AppStatus = "READY"
	Swappable               AppStatus = "SWAPPABLE"

	TestPlanV0   PayPlanType = "TEST_PLAN_V0"
	TestPlan10K  PayPlanType = "TEST_PLAN_10K"
	TestPlan90k  PayPlanType = "TEST_PLAN_90K"
	FreetierV0   PayPlanType = "FREETIER_V0"
	PayAsYouGoV0 PayPlanType = "PAY_AS_YOU_GO_V0"
	Enterprise   PayPlanType = "ENTERPRISE"
)

var (
	ValidAppStatuses = map[AppStatus]bool{
		"":                      true, // needed since it can be empty too
		AwaitingFreetierFunds:   true,
		AwaitingFreetierStaking: true,
		AwaitingFunds:           true,
		AwaitingFundsRemoval:    true,
		AwaitingGracePeriod:     true,
		AwaitingSlotFunds:       true,
		AwaitingSlotStaking:     true,
		AwaitingStaking:         true,
		AwaitingUnstaking:       true,
		Decomissioned:           true,
		InService:               true,
		Orphaned:                true,
		Ready:                   true,
		Swappable:               true,
	}

	ValidPayPlanTypes = map[PayPlanType]bool{
		"":           true, // needs to be allowed while the change for all apps to have plans is done
		TestPlanV0:   true,
		TestPlan10K:  true,
		TestPlan90k:  true,
		FreetierV0:   true,
		PayAsYouGoV0: true,
		Enterprise:   true,
	}
)

func (a *Application) Table() Table {
	return TableApplications
}
func (a *GatewayAAT) Table() Table {
	return TableGatewayAAT
}
func (s *GatewaySettings) Table() Table {
	return TableGatewaySettings
}
func (a *AppLimit) Table() Table {
	return TableAppLimits
}
func (s *NotificationSettings) Table() Table {
	return TableNotificationSettings
}

func (a *Application) DailyLimit() int {
	if a.Limit.PayPlan.Type == Enterprise {
		return a.Limit.CustomLimit
	}

	return a.Limit.PayPlan.Limit
}
func (a *Application) Validate() error {
	if !ValidAppStatuses[a.Status] {
		return ErrInvalidAppStatus
	}

	if !ValidPayPlanTypes[a.Limit.PayPlan.Type] {
		return ErrInvalidPayPlanType
	}

	if a.Limit.PayPlan.Type != Enterprise && a.Limit.CustomLimit != 0 {
		return ErrNotEnterprisePlan
	}
	return nil
}

func (u *UpdateApplication) Validate() error {
	if u == nil {
		return ErrNoFieldsToUpdate
	}
	if !ValidAppStatuses[u.Status] {
		return ErrInvalidAppStatus
	}
	if u.Limit != nil && !ValidPayPlanTypes[u.Limit.PayPlan.Type] {
		return ErrInvalidPayPlanType
	}
	if u.Limit != nil && u.Limit.PayPlan.Type != Enterprise && u.Limit.CustomLimit != 0 {
		return ErrNotEnterprisePlan
	}
	if u.Limit != nil && u.Limit.PayPlan.Type == Enterprise && u.Limit.CustomLimit == 0 {
		return ErrEnterprisePlanNeedsCustomLimit
	}
	return nil
}

/* LB Apps Table represents DB relationship of LBs and apps */
// do not change the tags, they're snake_case on purpose
type LbApp struct {
	LbID  string `json:"lb_id"`
	AppID string `json:"app_id"`
}

func (l *LbApp) Table() Table {
	return TableLbApps
}

/* Load Balancers Table */
type (
	LoadBalancer struct {
		ID                string        `json:"id"`
		Name              string        `json:"name"`
		UserID            string        `json:"userID"`
		ApplicationIDs    []string      `json:"applicationIDs,omitempty"`
		RequestTimeout    int           `json:"requestTimeout"`
		Gigastake         bool          `json:"gigastake"`
		GigastakeRedirect bool          `json:"gigastakeRedirect"`
		StickyOptions     StickyOptions `json:"stickinessOptions"`
		Applications      []*Application
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
	}
	StickyOptions struct {
		ID            string   `json:"id,omitempty"`
		Duration      string   `json:"duration"`
		StickyOrigins []string `json:"stickyOrigins"`
		StickyMax     int      `json:"stickyMax"`
		Stickiness    bool     `json:"stickiness"`
	}

	UpdateLoadBalancer struct {
		Name          string              `json:"name,omitempty"`
		StickyOptions UpdateStickyOptions `json:"stickinessOptions,omitempty"`
		Remove        bool                `json:"remove,omitempty"`
	}
	UpdateStickyOptions struct {
		ID            string   `json:"id,omitempty"`
		Duration      string   `json:"duration"`
		StickyOrigins []string `json:"stickyOrigins"`
		StickyMax     int      `json:"stickyMax"`
		Stickiness    *bool    `json:"stickiness"`
	}
)

func (l *LoadBalancer) Table() Table {
	return TableLoadBalancers
}
func (s *StickyOptions) Table() Table {
	return TableStickinessOptions
}
func (s *StickyOptions) IsEmpty() bool {
	if !s.Stickiness {
		return true
	}
	return len(s.StickyOrigins) == 0
}

/* Blockchains Table */
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
)

func (b *Blockchain) Table() Table {
	return TableBlockchains
}
func (r *Redirect) Table() Table {
	return TableRedirects
}

func (o *SyncCheckOptions) Table() Table {
	return TableSyncCheckOptions
}
