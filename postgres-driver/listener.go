package postgres_driver

import (
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pokt-foundation/portal-db/repository"
)

type Listener interface {
	NotificationChannel() <-chan *pq.Notification
	Listen(channel string) error
}

type notification struct {
	Table  repository.Table  `json:"table"`
	Action repository.Action `json:"action"`
	Data   any               `json:"data"`
}

// func (n notification) parseLoadBalancerNotification() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbLoadBalancer dbLoadBalancerJSON
// 	_ = json.Unmarshal(rawData, &dbLoadBalancer)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbLoadBalancer.toOutput(),
// 	}
// }

// func (n notification) parseStickinessOptionsNotification() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbStickinessOpts dbStickinessOptionsJSON
// 	_ = json.Unmarshal(rawData, &dbStickinessOpts)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbStickinessOpts.toOutput(),
// 	}
// }

// func (n notification) parseLbApps() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var lbApp repository.LbApp
// 	_ = json.Unmarshal(rawData, &lbApp)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   &lbApp,
// 	}
// }

// func (n notification) parseApplicationNotification() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbApp dbAppJSON
// 	_ = json.Unmarshal(rawData, &dbApp)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbApp.toOutput(),
// 	}
// }

// func (n notification) parseAppLimitNotification() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbAppLimit dbAppLimitJSON
// 	_ = json.Unmarshal(rawData, &dbAppLimit)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbAppLimit.toOutput(),
// 	}
// }

// func (n notification) parseGatewayAATNotification() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbGatewayAAT dbGatewayAATJSON
// 	_ = json.Unmarshal(rawData, &dbGatewayAAT)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbGatewayAAT.toOutput(),
// 	}
// }

// func (n notification) parseGatewaySettingsNotification() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbGatewaySettings dbGatewaySettingsJSON
// 	_ = json.Unmarshal(rawData, &dbGatewaySettings)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbGatewaySettings.toOutput(),
// 	}
// }

// func (n notification) parseNotificationSettingsNotification() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbNotificationSettings dbNotificationSettingsJSON
// 	_ = json.Unmarshal(rawData, &dbNotificationSettings)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbNotificationSettings.toOutput(),
// 	}
// }

// func (n notification) parseBlockchainNotification() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbBlockchain dbBlockchainJSON
// 	_ = json.Unmarshal(rawData, &dbBlockchain)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbBlockchain.toOutput(),
// 	}
// }

// func (n notification) parseRedirectNotification() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbRedirect dbRedirectJSON
// 	_ = json.Unmarshal(rawData, &dbRedirect)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbRedirect.toOutput(),
// 	}
// }

// func (n notification) parseSyncOptionsNotification() *repository.Notification {
// 	rawData, _ := json.Marshal(n.Data)
// 	var dbSyncOpts dbSyncCheckOptionsJSON
// 	_ = json.Unmarshal(rawData, &dbSyncOpts)

// 	return &repository.Notification{
// 		Table:  n.Table,
// 		Action: n.Action,
// 		Data:   dbSyncOpts.toOutput(),
// 	}
// }

func (n notification) parseNotification() *repository.Notification {
	switch n.Table {
	// case repository.TableLoadBalancers:
	// 	return n.parseLoadBalancerNotification()
	// case repository.TableStickinessOptions:
	// 	return n.parseStickinessOptionsNotification()

	// case repository.TableLbApps:
	// 	return n.parseLbApps()

	// case repository.TableApplications:
	// 	return n.parseApplicationNotification()
	// case repository.TableAppLimits:
	// 	return n.parseAppLimitNotification()
	// case repository.TableGatewayAAT:
	// 	return n.parseGatewayAATNotification()
	// case repository.TableGatewaySettings:
	// 	return n.parseGatewaySettingsNotification()
	// case repository.TableNotificationSettings:
	// 	return n.parseNotificationSettingsNotification()

	// case repository.TableBlockchains:
	// 	return n.parseBlockchainNotification()
	// case repository.TableRedirects:
	// 	return n.parseRedirectNotification()
	// case repository.TableSyncCheckOptions:
	// 	return n.parseSyncOptionsNotification()
	}

	return nil
}

func parsePQNotification(n *pq.Notification, outCh chan *repository.Notification) {
	if n != nil {
		var notification notification
		_ = json.Unmarshal([]byte(n.Extra), &notification)
		outCh <- notification.parseNotification()
	}
}

func Listen(inCh <-chan *pq.Notification, outCh chan *repository.Notification) {
	for {
		n := <-inCh
		go parsePQNotification(n, outCh)
	}
}

func (d *PostgresDriver) CloseListener() {
	close(d.notification)
}
