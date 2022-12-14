package types

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
	TableApplications         Table = "applications"
	TableAppLimits            Table = "app_limits"
	TableGatewayAAT           Table = "gateway_aat"
	TableGatewaySettings      Table = "gateway_settings"
	TableNotificationSettings Table = "notification_settings"
	TableLbApps               Table = "lb_apps"
	TableLoadBalancers        Table = "loadbalancers"
	TableStickinessOptions    Table = "stickiness_options"
	TableBlockchains          Table = "blockchains"
	TableRedirects            Table = "redirects"
	TableSyncCheckOptions     Table = "sync_check_options"

	ActionInsert Action = "INSERT"
	ActionUpdate Action = "UPDATE"
)

type SavedOnDB interface {
	Table() Table
}

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
func (l *LbApp) Table() Table {
	return TableLbApps
}
func (l *LoadBalancer) Table() Table {
	return TableLoadBalancers
}
func (s *StickyOptions) Table() Table {
	return TableStickinessOptions
}
func (b *Blockchain) Table() Table {
	return TableBlockchains
}
func (r *Redirect) Table() Table {
	return TableRedirects
}
func (o *SyncCheckOptions) Table() Table {
	return TableSyncCheckOptions
}
