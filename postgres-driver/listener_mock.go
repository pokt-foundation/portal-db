package postgresdriver

import (
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pokt-foundation/portal-db/repository"
)

type ListenerMock struct {
	Notify chan *pq.Notification
}

func NewListenerMock() *ListenerMock {
	return &ListenerMock{
		Notify: make(chan *pq.Notification, 32),
	}
}

func (l *ListenerMock) NotificationChannel() <-chan *pq.Notification {
	return l.Notify
}

func (l *ListenerMock) Listen(channel string) error {
	return nil
}

func gatewaySettingsIsNull(settings repository.GatewaySettings) bool {
	return settings.SecretKey == "" &&
		len(settings.WhitelistOrigins) == 0 &&
		len(settings.WhitelistUserAgents) == 0 &&
		len(settings.WhitelistContracts) == 0 &&
		len(settings.WhitelistMethods) == 0 &&
		len(settings.WhitelistBlockchains) == 0
}

func applicationInputs(mainTableAction, sideTablesAction repository.Action, content repository.SavedOnDB) []inputStruct {
	app := content.(*repository.Application)

	var inputs []inputStruct

	inputs = append(inputs, inputStruct{
		action: mainTableAction,
		table:  repository.TableApplications,
		input: dbAppJSON{
			ApplicationID:      app.ID,
			UserID:             app.UserID,
			Name:               app.Name,
			ContactEmail:       app.ContactEmail,
			Description:        app.Description,
			Owner:              app.Owner,
			URL:                app.URL,
			Status:             string(app.Status),
			CreatedAt:          app.CreatedAt.Format(psqlDateLayout),
			UpdatedAt:          app.UpdatedAt.Format(psqlDateLayout),
			FirstDateSurpassed: app.FirstDateSurpassed.Format(psqlDateLayout),
			Dummy:              app.Dummy,
		},
	})

	inputs = append(inputs, inputStruct{
		action: sideTablesAction,
		table:  repository.TableAppLimits,
		input: dbAppLimitJSON{
			ApplicationID: app.ID,
			PlanType:      app.Limit.PayPlan.Type,
			CustomLimit:   app.Limit.CustomLimit,
		},
	})

	if app.GatewayAAT != (repository.GatewayAAT{}) {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  repository.TableGatewayAAT,
			input: dbGatewayAATJSON{
				ApplicationID:   app.ID,
				Address:         app.GatewayAAT.Address,
				ClientPublicKey: app.GatewayAAT.ClientPublicKey,
				PrivateKey:      app.GatewayAAT.PrivateKey,
				PublicKey:       app.GatewayAAT.ApplicationPublicKey,
				Signature:       app.GatewayAAT.ApplicationSignature,
				Version:         app.GatewayAAT.Version,
			},
		})
	}

	if !gatewaySettingsIsNull(app.GatewaySettings) {
		contracts, methods := marshalWhitelistContractsAndMethods(app.GatewaySettings.WhitelistContracts,
			app.GatewaySettings.WhitelistMethods)

		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  repository.TableGatewaySettings,
			input: dbGatewaySettingsJSON{
				ApplicationID:        app.ID,
				SecretKey:            app.GatewaySettings.SecretKey,
				SecretKeyRequired:    app.GatewaySettings.SecretKeyRequired,
				WhitelistContracts:   contracts,
				WhitelistMethods:     methods,
				WhitelistOrigins:     app.GatewaySettings.WhitelistOrigins,
				WhitelistUserAgents:  app.GatewaySettings.WhitelistUserAgents,
				WhitelistBlockchains: app.GatewaySettings.WhitelistBlockchains,
			},
		})
	}

	if app.NotificationSettings != (repository.NotificationSettings{}) {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  repository.TableNotificationSettings,
			input: dbNotificationSettingsJSON{
				ApplicationID: app.ID,
				SignedUp:      app.NotificationSettings.SignedUp,
				Quarter:       app.NotificationSettings.Quarter,
				Half:          app.NotificationSettings.Half,
				ThreeQuarters: app.NotificationSettings.ThreeQuarters,
				Full:          app.NotificationSettings.Full,
			},
		})
	}

	return inputs
}

func appLimitInputs(sideTablesAction repository.Action, content repository.SavedOnDB) []inputStruct {
	appLimit := content.(*repository.AppLimit)

	var inputs []inputStruct

	inputs = append(inputs, inputStruct{
		action: sideTablesAction,
		table:  repository.TableAppLimits,
		input: dbAppLimitJSON{
			ApplicationID: appLimit.ID,
			PlanType:      appLimit.PayPlan.Type,
			CustomLimit:   0,
		},
	})

	return inputs
}

func blockchainInputs(mainTableAction, sideTablesAction repository.Action, content repository.SavedOnDB) []inputStruct {
	blockchain := content.(*repository.Blockchain)

	var inputs []inputStruct

	inputs = append(inputs, inputStruct{
		action: mainTableAction,
		table:  repository.TableBlockchains,
		input: dbBlockchainJSON{
			BlockchainID:      blockchain.ID,
			Altruist:          blockchain.Altruist,
			Blockchain:        blockchain.Blockchain,
			ChainID:           blockchain.ChainID,
			ChainIDCheck:      blockchain.ChainIDCheck,
			ChainPath:         blockchain.Path,
			Description:       blockchain.Description,
			EnforceResult:     blockchain.EnforceResult,
			Network:           blockchain.Network,
			Ticker:            blockchain.Ticker,
			BlockchainAliases: blockchain.BlockchainAliases,
			LogLimitBlocks:    blockchain.LogLimitBlocks,
			RequestTimeout:    blockchain.RequestTimeout,
			Active:            blockchain.Active,
			CreatedAt:         blockchain.CreatedAt.Format(psqlDateLayout),
			UpdatedAt:         blockchain.UpdatedAt.Format(psqlDateLayout),
		},
	})

	if blockchain.SyncCheckOptions != (repository.SyncCheckOptions{}) {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  repository.TableSyncCheckOptions,
			input: dbSyncCheckOptionsJSON{
				BlockchainID: blockchain.SyncCheckOptions.BlockchainID,
				Body:         blockchain.SyncCheckOptions.Body,
				Path:         blockchain.SyncCheckOptions.Path,
				ResultKey:    blockchain.SyncCheckOptions.ResultKey,
				Allowance:    blockchain.SyncCheckOptions.Allowance,
			},
		})
	}

	return inputs
}

func loadBalancerInputs(mainTableAction, sideTablesAction repository.Action, content repository.SavedOnDB) []inputStruct {
	lb := content.(*repository.LoadBalancer)

	var inputs []inputStruct

	inputs = append(inputs, inputStruct{
		action: mainTableAction,
		table:  repository.TableLoadBalancers,
		input: dbLoadBalancerJSON{
			LbID:              lb.ID,
			Name:              lb.Name,
			UserID:            lb.UserID,
			RequestTimeout:    lb.RequestTimeout,
			Gigastake:         lb.Gigastake,
			GigastakeRedirect: lb.GigastakeRedirect,
			CreatedAt:         lb.CreatedAt.Format(psqlDateLayout),
			UpdatedAt:         lb.UpdatedAt.Format(psqlDateLayout),
		},
	})

	if !lb.StickyOptions.IsEmpty() {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  repository.TableStickinessOptions,
			input: dbStickinessOptionsJSON{
				LbID:       lb.ID,
				Duration:   lb.StickyOptions.Duration,
				Origins:    lb.StickyOptions.StickyOrigins,
				StickyMax:  lb.StickyOptions.StickyMax,
				Stickiness: lb.StickyOptions.Stickiness,
			},
		})
	}

	for _, appID := range lb.ApplicationIDs {
		inputs = append(inputs, inputStruct{
			action: sideTablesAction,
			table:  repository.TableLbApps,
			input: repository.LbApp{
				LbID:  lb.ID,
				AppID: appID,
			},
		})
	}

	return inputs
}

func redirectInput(action repository.Action, content repository.SavedOnDB) inputStruct {
	redirect := content.(*repository.Redirect)

	return inputStruct{
		action: action,
		table:  repository.TableRedirects,
		input: dbRedirectJSON{
			BlockchainID:   redirect.BlockchainID,
			Alias:          redirect.Alias,
			LoadBalancerID: redirect.LoadBalancerID,
			Domain:         redirect.Domain,
			CreatedAt:      redirect.CreatedAt.Format(psqlDateLayout),
			UpdatedAt:      redirect.UpdatedAt.Format(psqlDateLayout),
		},
	}
}

type inputStruct struct {
	action repository.Action
	table  repository.Table
	input  any
}

func mockInput(inStruct inputStruct) *pq.Notification {
	notification, _ := json.Marshal(notification{
		Table:  inStruct.table,
		Action: inStruct.action,
		Data:   inStruct.input,
	})

	return &pq.Notification{
		Extra: string(notification),
	}
}

func mockContent(mainTableAction, sideTablesAction repository.Action, content repository.SavedOnDB) []*pq.Notification {
	var inputs []inputStruct

	switch content.(type) {
	case *repository.Application:
		inputs = applicationInputs(mainTableAction, sideTablesAction, content)
	case *repository.AppLimit:
		inputs = appLimitInputs(sideTablesAction, content)
	case *repository.Blockchain:
		inputs = blockchainInputs(mainTableAction, sideTablesAction, content)
	case *repository.LoadBalancer:
		inputs = loadBalancerInputs(mainTableAction, sideTablesAction, content)
	case *repository.Redirect:
		inputs = []inputStruct{redirectInput(mainTableAction, content)}
	default:
		panic("type not supported")
	}

	var notifications []*pq.Notification

	for _, input := range inputs {
		notifications = append(notifications, mockInput(input))
	}

	return notifications
}

func (l *ListenerMock) MockEvent(mainTableAction, sideTablesAction repository.Action, content repository.SavedOnDB) {
	notifications := mockContent(mainTableAction, sideTablesAction, content)

	for _, notification := range notifications {
		l.Notify <- notification
	}
}
