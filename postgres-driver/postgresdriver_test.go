package postgresdriver

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pokt-foundation/portal-db/repository"
	"github.com/stretchr/testify/suite"
)

func PrettyString(label string, thing interface{}) {
	jsonThing, _ := json.Marshal(thing)
	str := string(jsonThing)

	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, []byte(str), "", "    ")
	output := prettyJSON.String()

	fmt.Println(label, output)
}

var ctx = context.Background()

func Test_RunPGDriverSuite(t *testing.T) {
	testSuite := new(PGDriverTestSuite)
	testSuite.connectionString = connectionString

	suite.Run(t, testSuite)
}

func (ts *PGDriverTestSuite) Test_ReadPayPlans() {
	tests := []struct {
		name     string
		payPlans []*repository.PayPlan
		err      error
	}{
		{
			name: "Should return all PayPlans from the database ordered by plan_type",
			payPlans: []*repository.PayPlan{
				{Type: repository.Enterprise, Limit: 0},
				{Type: repository.FreetierV0, Limit: 250000},
				{Type: repository.PayAsYouGoV0, Limit: 0},
				{Type: repository.TestPlan10K, Limit: 10000},
				{Type: repository.TestPlan90k, Limit: 90000},
				{Type: repository.TestPlanV0, Limit: 100},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		payPlans, err := ts.driver.ReadPayPlans(ctx)
		ts.Equal(test.payPlans, payPlans)
		ts.Equal(test.err, err)
	}
}

func (ts *PGDriverTestSuite) Test_ReadApplications() {
	tests := []struct {
		name         string
		applications []*repository.Application
		err          error
	}{
		{
			name: "Should return all Applications from the database ordered by application_id",
			applications: []*repository.Application{
				{
					ID:     "test_app_47hfnths73j2se",
					UserID: "test_user_1dbffbdfeeb225",
					Name:   "pokt_app_123",
					URL:    "https://test.app123.io",
					Dummy:  true,
					Status: repository.InService,
					GatewayAAT: repository.GatewayAAT{
						Address:              "test_34715cae753e67c75fbb340442e7de8e",
						ApplicationPublicKey: "test_11b8d394ca331d7c7a71ca1896d630f6",
						ApplicationSignature: "test_89a3af6a587aec02cfade6f5000424c2",
						ClientPublicKey:      "test_1dc39a2e5a84a35bf030969a0b3231f7",
						PrivateKey:           "test_d2ce53f115f4ecb2208e9188800a85cf",
					},
					GatewaySettings: repository.GatewaySettings{
						SecretKey:         "test_40f482d91a5ef2300ebb4e2308c",
						SecretKeyRequired: true,
					},
					Limit: repository.AppLimit{
						PayPlan: repository.PayPlan{Type: repository.FreetierV0, Limit: 250_000},
					},
					NotificationSettings: repository.NotificationSettings{
						SignedUp:      true,
						Quarter:       false,
						Half:          false,
						ThreeQuarters: true,
						Full:          true,
					},
				},
				{
					ID:     "test_app_5hdf7sh23jd828",
					UserID: "test_user_04228205bd261a",
					Name:   "pokt_app_456",
					URL:    "https://test.app456.io",
					Dummy:  true,
					Status: repository.InService,
					GatewayAAT: repository.GatewayAAT{
						Address:              "test_558c0225c7019e14ccf2e7379ad3eb50",
						ApplicationPublicKey: "test_96c981db344ab6920b7e87853838e285",
						ApplicationSignature: "test_1272a8ab4cbbf636f09bf4fa5395b885",
						ClientPublicKey:      "test_d709871777b89ed3051190f229ea3f01",
						PrivateKey:           "test_53e50765d8bc1fb41b3b0065dd8094de",
					},
					GatewaySettings: repository.GatewaySettings{
						SecretKey:         "test_90210ac4bdd3423e24877d1ff92",
						SecretKeyRequired: false,
					},
					Limit: repository.AppLimit{
						PayPlan:     repository.PayPlan{Type: repository.Enterprise},
						CustomLimit: 2_000_000,
					},
					NotificationSettings: repository.NotificationSettings{
						SignedUp:      true,
						Quarter:       false,
						Half:          false,
						ThreeQuarters: true,
						Full:          true,
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		applications, err := ts.driver.ReadApplications(ctx)
		ts.Equal(test.err, err)
		for i, app := range applications {
			ts.Equal(test.applications[i].ID, app.ID)
			ts.Equal(test.applications[i].UserID, app.UserID)
			ts.Equal(test.applications[i].Name, app.Name)
			ts.Equal(test.applications[i].URL, app.URL)
			ts.Equal(test.applications[i].Dummy, app.Dummy)
			ts.Equal(test.applications[i].Status, app.Status)
			ts.Equal(test.applications[i].GatewayAAT, app.GatewayAAT)
			ts.Equal(test.applications[i].GatewaySettings, app.GatewaySettings)
			ts.Equal(test.applications[i].Limit, app.Limit)
			ts.Equal(test.applications[i].NotificationSettings, app.NotificationSettings)
			ts.NotEmpty(app.CreatedAt)
			ts.NotEmpty(app.UpdatedAt)
		}
	}
}

func (ts *PGDriverTestSuite) Test_ReadLoadBalancers() {
	tests := []struct {
		name          string
		loadBalancers []*repository.LoadBalancer
		err           error
	}{
		{
			name: "Should return all Load Balancers from the database ordered by lb_id",
			loadBalancers: []*repository.LoadBalancer{
				{
					ID:                "test_lb_34987u329rfn23f",
					Name:              "pokt_app_123",
					UserID:            "test_user_1dbffbdfeeb225",
					ApplicationIDs:    []string{"test_app_47hfnths73j2se"},
					RequestTimeout:    5_000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: repository.StickyOptions{
						Duration:      "60",
						StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
						StickyMax:     300,
						Stickiness:    true,
					},
				},
				{
					ID:                "test_lb_34gg4g43g34g5hh",
					Name:              "test_lb_redirect",
					UserID:            "test_user_redirect233344",
					ApplicationIDs:    []string{""},
					RequestTimeout:    5_000,
					Gigastake:         false,
					GigastakeRedirect: false,
					StickyOptions: repository.StickyOptions{
						Duration:      "20",
						StickyOrigins: []string{"test-extension://", "test-extension2://"},
						StickyMax:     600,
						Stickiness:    false,
					},
				},
				{
					ID:                "test_lb_3890ru23jfi32fj",
					Name:              "pokt_app_456",
					UserID:            "test_user_04228205bd261a",
					ApplicationIDs:    []string{"test_app_5hdf7sh23jd828"},
					RequestTimeout:    5_000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: repository.StickyOptions{
						Duration:      "40",
						StickyOrigins: []string{"chrome-extension://"},
						StickyMax:     400,
						Stickiness:    true,
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		loadBalancers, err := ts.driver.ReadLoadBalancers(ctx)
		ts.Equal(test.err, err)
		for i, loadBalancer := range loadBalancers {
			ts.Equal(test.loadBalancers[i].ID, loadBalancer.ID)
			ts.Equal(test.loadBalancers[i].UserID, loadBalancer.UserID)
			ts.Equal(test.loadBalancers[i].Name, loadBalancer.Name)
			ts.Equal(test.loadBalancers[i].UserID, loadBalancer.UserID)
			ts.Equal(test.loadBalancers[i].ApplicationIDs, loadBalancer.ApplicationIDs)
			ts.Equal(test.loadBalancers[i].RequestTimeout, loadBalancer.RequestTimeout)
			ts.Equal(test.loadBalancers[i].Gigastake, loadBalancer.Gigastake)
			ts.Equal(test.loadBalancers[i].GigastakeRedirect, loadBalancer.GigastakeRedirect)
			ts.Equal(test.loadBalancers[i].StickyOptions, loadBalancer.StickyOptions)
			ts.NotEmpty(loadBalancer.CreatedAt)
			ts.NotEmpty(loadBalancer.UpdatedAt)
		}
	}
}

func (ts *PGDriverTestSuite) Test_ReadBlockchains() {
	tests := []struct {
		name        string
		blockchains []*repository.Blockchain
		err         error
	}{
		{
			name: "Should return all Load Balancers from the database ordered by blockchain_id",
			blockchains: []*repository.Blockchain{
				{
					ID:                "0001",
					Altruist:          "https://test:329y293uhfniu23f8@shared-test2.nodes.pokt.network:12345",
					Blockchain:        "pokt-mainnet",
					Description:       "POKT Network Mainnet",
					EnforceResult:     "JSON",
					Network:           "POKT-mainnet",
					Ticker:            "POKT",
					BlockchainAliases: []string{"pokt-mainnet"},
					LogLimitBlocks:    100_000,
					Active:            true,
					Redirects: []repository.Redirect{
						{
							Alias:          "test-mainnet",
							Domain:         "test-rpc1.testnet.pokt.network",
							LoadBalancerID: "test_lb_34gg4g43g34g5hh",
						},
						{
							Alias:          "test-mainnet",
							Domain:         "test-rpc2.testnet.pokt.network",
							LoadBalancerID: "test_lb_34gg4g43g34g5hh",
						},
					},
					SyncCheckOptions: repository.SyncCheckOptions{
						Body:      `{}`,
						Path:      "/v1/query/height",
						ResultKey: "height",
						Allowance: 1,
					},
				},
				{
					ID:                "0021",
					Altruist:          "https://test:2r980u32fh239hf@shared-test2.nodes.eth.network:12345",
					Blockchain:        "eth-mainnet",
					ChainID:           "1",
					ChainIDCheck:      `{\"method\":\"eth_chainId\",\"id\":1,\"jsonrpc\":\"2.0\"}`,
					Description:       "Ethereum Mainnet",
					EnforceResult:     "JSON",
					Network:           "ETH-1",
					Ticker:            "ETH",
					BlockchainAliases: []string{"eth-mainnet"},
					LogLimitBlocks:    100_000,
					Active:            true,
					Redirects: []repository.Redirect{
						{
							Alias:          "eth-mainnet",
							Domain:         "test-rpc.testnet.eth.network",
							LoadBalancerID: "test_lb_34gg4g43g34g5hh",
						},
					},
					SyncCheckOptions: repository.SyncCheckOptions{
						Body:      `{\"method\":\"eth_blockNumber\",\"id\":1,\"jsonrpc\":\"2.0\"}`,
						ResultKey: "result",
						Allowance: 5,
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		blockchains, err := ts.driver.ReadBlockchains(ctx)
		ts.Equal(test.err, err)
		for i, blockchain := range blockchains {
			ts.Equal(test.blockchains[i].ID, blockchain.ID)
			ts.Equal(test.blockchains[i].ID, blockchain.ID)
			ts.Equal(test.blockchains[i].Altruist, blockchain.Altruist)
			ts.Equal(test.blockchains[i].Blockchain, blockchain.Blockchain)
			ts.Equal(test.blockchains[i].ChainID, blockchain.ChainID)
			ts.Equal(test.blockchains[i].ChainIDCheck, blockchain.ChainIDCheck)
			ts.Equal(test.blockchains[i].Description, blockchain.Description)
			ts.Equal(test.blockchains[i].EnforceResult, blockchain.EnforceResult)
			ts.Equal(test.blockchains[i].Network, blockchain.Network)
			ts.Equal(test.blockchains[i].Path, blockchain.Path)
			ts.Equal(test.blockchains[i].SyncCheck, blockchain.SyncCheck)
			ts.Equal(test.blockchains[i].Ticker, blockchain.Ticker)
			ts.Equal(test.blockchains[i].BlockchainAliases, blockchain.BlockchainAliases)
			ts.Equal(test.blockchains[i].LogLimitBlocks, blockchain.LogLimitBlocks)
			ts.Equal(test.blockchains[i].RequestTimeout, blockchain.RequestTimeout)
			ts.Equal(test.blockchains[i].SyncAllowance, blockchain.SyncAllowance)
			ts.Equal(test.blockchains[i].Active, blockchain.Active)
			ts.Equal(test.blockchains[i].Redirects, blockchain.Redirects)
			ts.Equal(test.blockchains[i].SyncCheckOptions, blockchain.SyncCheckOptions)
			ts.NotEmpty(blockchain.CreatedAt)
			ts.NotEmpty(blockchain.UpdatedAt)
		}
	}
}

func (ts *PGDriverTestSuite) Test_WriteLoadBalancer() {
	tests := []struct {
		name               string
		loadBalancerInputs []*repository.LoadBalancer
		expectedNumOfLBs   int
		expectedLB         SelectOneLoadBalancerRow
		err                error
	}{
		{
			name: "Should create a single load balancer successfully with correct input",
			loadBalancerInputs: []*repository.LoadBalancer{
				{
					Name:              "pokt_app_789",
					UserID:            "test_user_47fhsd75jd756sh",
					RequestTimeout:    5000,
					Gigastake:         true,
					GigastakeRedirect: true,
					ApplicationIDs:    []string{"test_app_47hfnths73j2se"},
					StickyOptions: repository.StickyOptions{
						Duration:      "70",
						StickyOrigins: []string{"chrome-extension://"},
						StickyMax:     400,
						Stickiness:    true,
					},
				},
			},
			expectedNumOfLBs: 4,
			expectedLB: SelectOneLoadBalancerRow{
				Name:              sql.NullString{Valid: true, String: "pokt_app_789"},
				UserID:            sql.NullString{Valid: true, String: "test_user_47fhsd75jd756sh"},
				RequestTimeout:    sql.NullInt32{Valid: true, Int32: 5000},
				Gigastake:         sql.NullBool{Valid: true, Bool: true},
				GigastakeRedirect: sql.NullBool{Valid: true, Bool: true},
				Duration:          sql.NullString{Valid: true, String: "70"},
				StickyMax:         sql.NullInt32{Valid: true, Int32: 400},
				Stickiness:        sql.NullBool{Valid: true, Bool: true},
				Origins:           []string{"chrome-extension://"},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		for _, input := range test.loadBalancerInputs {
			createdLB, err := ts.driver.WriteLoadBalancer(ctx, input)
			ts.Equal(test.err, err)
			ts.Len(createdLB.ID, 24)
			ts.Equal(input.Name, createdLB.Name)

			loadBalancers, err := ts.driver.ReadLoadBalancers(ctx)
			ts.Equal(test.err, err)
			ts.Len(loadBalancers, test.expectedNumOfLBs)

			loadBalancer, err := ts.driver.SelectOneLoadBalancer(ctx, createdLB.ID)
			ts.Equal(test.err, err)
			for _, testInput := range test.loadBalancerInputs {
				if testInput.Name == loadBalancer.Name.String {
					ts.Equal(createdLB.ID, loadBalancer.LbID)
					ts.Equal(test.expectedLB.UserID, loadBalancer.UserID)
					ts.Equal(test.expectedLB.Name, loadBalancer.Name)
					ts.Equal(test.expectedLB.UserID, loadBalancer.UserID)
					ts.Equal(test.expectedLB.RequestTimeout, loadBalancer.RequestTimeout)
					ts.Equal(test.expectedLB.Gigastake, loadBalancer.Gigastake)
					ts.Equal(test.expectedLB.GigastakeRedirect, loadBalancer.GigastakeRedirect)
					ts.Equal(test.expectedLB.Duration, loadBalancer.Duration)
					ts.Equal(test.expectedLB.Origins, loadBalancer.Origins)
					ts.Equal(test.expectedLB.StickyMax, loadBalancer.StickyMax)
					ts.Equal(test.expectedLB.Stickiness, loadBalancer.Stickiness)
					ts.NotEmpty(loadBalancer.CreatedAt)
					ts.NotEmpty(loadBalancer.UpdatedAt)
				}

			}
		}
	}
}

func (ts *PGDriverTestSuite) Test_UpdateLoadBalancer() {
	tests := []struct {
		name                string
		loadBalancerID      string
		loadBalancerUpdate  *repository.UpdateLoadBalancer
		expectedAfterUpdate SelectOneLoadBalancerRow
		err                 error
	}{
		{
			name:           "Should update a single load balancer successfully with all fields",
			loadBalancerID: "test_lb_34987u329rfn23f",
			loadBalancerUpdate: &repository.UpdateLoadBalancer{
				Name: "pokt_app_updated",
				StickyOptions: repository.UpdateStickyOptions{
					Duration:      "100",
					StickyOrigins: []string{"chrome-extension://", "test-ext://"},
					StickyMax:     500,
					Stickiness:    boolPointer(false),
				},
			},
			expectedAfterUpdate: SelectOneLoadBalancerRow{
				Name:       sql.NullString{Valid: true, String: "pokt_app_updated"},
				Duration:   sql.NullString{Valid: true, String: "100"},
				StickyMax:  sql.NullInt32{Valid: true, Int32: 500},
				Stickiness: sql.NullBool{Valid: true, Bool: false},
				Origins:    []string{"chrome-extension://", "test-ext://"},
			},
			err: nil,
		},
		{
			name:           "Should update a single load balancer successfully with only some sticky options fields",
			loadBalancerID: "test_lb_3890ru23jfi32fj",
			loadBalancerUpdate: &repository.UpdateLoadBalancer{
				Name: "pokt_app_updated_2",
				StickyOptions: repository.UpdateStickyOptions{
					Duration: "100",
				},
			},
			expectedAfterUpdate: SelectOneLoadBalancerRow{
				Name:       sql.NullString{Valid: true, String: "pokt_app_updated_2"},
				Duration:   sql.NullString{Valid: true, String: "100"},
				StickyMax:  sql.NullInt32{Valid: true, Int32: 400},
				Stickiness: sql.NullBool{Valid: true, Bool: true},
				Origins:    []string{"chrome-extension://"},
			},
			err: nil,
		},
		{
			name:           "Should update a single load balancer successfully with no sticky options fields",
			loadBalancerID: "test_lb_34gg4g43g34g5hh",
			loadBalancerUpdate: &repository.UpdateLoadBalancer{
				Name: "pokt_app_updated_3",
			},
			expectedAfterUpdate: SelectOneLoadBalancerRow{
				Name:       sql.NullString{Valid: true, String: "pokt_app_updated_3"},
				Duration:   sql.NullString{Valid: true, String: "20"},
				StickyMax:  sql.NullInt32{Valid: true, Int32: 600},
				Stickiness: sql.NullBool{Valid: true, Bool: false},
				Origins:    []string{"test-extension://", "test-extension2://"},
			},
			err: nil,
		},
		{
			name:           "Should update a single load balancer successfully with only sticky options origin field",
			loadBalancerID: "test_lb_34gg4g43g34g5hh",
			loadBalancerUpdate: &repository.UpdateLoadBalancer{
				StickyOptions: repository.UpdateStickyOptions{
					StickyOrigins: []string{"chrome-extension://", "test-ext://"},
				},
			},
			expectedAfterUpdate: SelectOneLoadBalancerRow{
				Name:       sql.NullString{Valid: true, String: "pokt_app_updated_3"},
				Duration:   sql.NullString{Valid: true, String: "20"},
				StickyMax:  sql.NullInt32{Valid: true, Int32: 600},
				Stickiness: sql.NullBool{Valid: true, Bool: false},
				Origins:    []string{"chrome-extension://", "test-ext://"},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		_, err := ts.driver.SelectOneLoadBalancer(ctx, test.loadBalancerID)
		ts.Equal(test.err, err)

		err = ts.driver.UpdateLoadBalancer(ctx, test.loadBalancerID, test.loadBalancerUpdate)
		ts.Equal(test.err, err)

		lbAfterUpdate, err := ts.driver.SelectOneLoadBalancer(ctx, test.loadBalancerID)
		ts.Equal(test.err, err)
		ts.Equal(test.expectedAfterUpdate.Name, lbAfterUpdate.Name)
		ts.Equal(test.expectedAfterUpdate.Duration, lbAfterUpdate.Duration)
		ts.Equal(test.expectedAfterUpdate.Origins, lbAfterUpdate.Origins)
		ts.Equal(test.expectedAfterUpdate.StickyMax, lbAfterUpdate.StickyMax)
		ts.Equal(test.expectedAfterUpdate.Stickiness, lbAfterUpdate.Stickiness)
	}
}

func (ts *PGDriverTestSuite) Test_RemoveLoadBalancer() {
	tests := []struct {
		name           string
		loadBalancerID string
		err            error
	}{
		{
			name:           "Should remove a single load balancer successfully with correct input",
			loadBalancerID: "test_lb_34gg4g43g34g5hh",
			err:            nil,
		},
	}

	for _, test := range tests {
		err := ts.driver.RemoveLoadBalancer(ctx, test.loadBalancerID)
		ts.Equal(test.err, err)

		lbAfterRemove, err := ts.driver.SelectOneLoadBalancer(ctx, test.loadBalancerID)
		ts.Equal(test.err, err)
		ts.Empty(lbAfterRemove.UserID.String)
	}
}

func (ts *PGDriverTestSuite) Test_WriteApplication() {
	tests := []struct {
		name              string
		appInputs         []*repository.Application
		expectedNumOfApps int
		expectedApp       SelectOneApplicationRow
		err               error
	}{
		{
			name: "Should create a single load balancer successfully with correct input",
			appInputs: []*repository.Application{
				{
					Name:   "pokt_app_789",
					UserID: "test_user_47fhsd75jd756sh",
					Dummy:  true,
					Status: repository.AppStatus("IN_SERVICE"),
					GatewayAAT: repository.GatewayAAT{
						Address:              "test_e209a2d1f3454ddc69cb9333d547bbcf",
						ApplicationPublicKey: "test_b95c35affacf6df4a5585388490542f0",
						ApplicationSignature: "test_e59760339d9ce02972d1080d73446c90",
						ClientPublicKey:      "test_d591178ab3f48f45b243303fe77dc8c3",
						PrivateKey:           "test_f403700aed7e039c0a8fc2dd22da6fd9",
					},
					GatewaySettings: repository.GatewaySettings{
						SecretKey:         "test_489574398f34uhf4uhjf9328jf23f98j",
						SecretKeyRequired: true,
					},
					Limit: repository.AppLimit{
						PayPlan: repository.PayPlan{Type: repository.PayPlanType("FREETIER_V0")},
					},
					NotificationSettings: repository.NotificationSettings{
						SignedUp:      true,
						Quarter:       false,
						Half:          false,
						ThreeQuarters: true,
						Full:          true,
					},
				},
			},
			expectedNumOfApps: 3,
			expectedApp: SelectOneApplicationRow{
				Name:              sql.NullString{Valid: true, String: "pokt_app_789"},
				UserID:            sql.NullString{Valid: true, String: "test_user_47fhsd75jd756sh"},
				Dummy:             sql.NullBool{Valid: true, Bool: true},
				Status:            sql.NullString{Valid: true, String: "IN_SERVICE"},
				GaAddress:         sql.NullString{Valid: true, String: "test_e209a2d1f3454ddc69cb9333d547bbcf"},
				GaClientPublicKey: sql.NullString{Valid: true, String: "test_d591178ab3f48f45b243303fe77dc8c3"},
				GaPrivateKey:      sql.NullString{Valid: true, String: "test_f403700aed7e039c0a8fc2dd22da6fd9"},
				GaPublicKey:       sql.NullString{Valid: true, String: "test_b95c35affacf6df4a5585388490542f0"},
				GaSignature:       sql.NullString{Valid: true, String: "test_e59760339d9ce02972d1080d73446c90"},
				SecretKey:         sql.NullString{Valid: true, String: "test_489574398f34uhf4uhjf9328jf23f98j"},
				SecretKeyRequired: sql.NullBool{Valid: true, Bool: true},
				SignedUp:          sql.NullBool{Valid: true, Bool: true},
				OnQuarter:         sql.NullBool{Valid: true, Bool: false},
				OnHalf:            sql.NullBool{Valid: true, Bool: false},
				OnThreeQuarters:   sql.NullBool{Valid: true, Bool: true},
				OnFull:            sql.NullBool{Valid: true, Bool: true},
				PayPlan:           sql.NullString{Valid: true, String: "FREETIER_V0"},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		for _, input := range test.appInputs {
			createdApp, err := ts.driver.WriteApplication(ctx, input)
			ts.Equal(test.err, err)
			ts.Len(createdApp.ID, 24)
			ts.Equal(input.Name, createdApp.Name)

			apps, err := ts.driver.ReadApplications(ctx)
			ts.Equal(test.err, err)
			ts.Len(apps, test.expectedNumOfApps)

			app, err := ts.driver.SelectOneApplication(ctx, createdApp.ID)
			ts.Equal(test.err, err)
			for _, testInput := range test.appInputs {
				if testInput.Name == app.Name.String {
					ts.Equal(createdApp.ID, app.ApplicationID)
					ts.Equal(test.expectedApp.Dummy, app.Dummy)
					ts.Equal(test.expectedApp.Status, app.Status)
					ts.Equal(test.expectedApp.GaAddress, app.GaAddress)
					ts.Equal(test.expectedApp.GaClientPublicKey, app.GaClientPublicKey)
					ts.Equal(test.expectedApp.GaPrivateKey, app.GaPrivateKey)
					ts.Equal(test.expectedApp.GaPublicKey, app.GaPublicKey)
					ts.Equal(test.expectedApp.GaSignature, app.GaSignature)
					ts.Equal(test.expectedApp.SecretKey, app.SecretKey)
					ts.Equal(test.expectedApp.SecretKeyRequired, app.SecretKeyRequired)
					ts.Equal(test.expectedApp.SignedUp, app.SignedUp)
					ts.Equal(test.expectedApp.OnQuarter, app.OnQuarter)
					ts.Equal(test.expectedApp.OnHalf, app.OnHalf)
					ts.Equal(test.expectedApp.OnThreeQuarters, app.OnThreeQuarters)
					ts.Equal(test.expectedApp.OnFull, app.OnFull)
					ts.Equal(test.expectedApp.PayPlan, app.PayPlan)
					ts.NotEmpty(app.CreatedAt)
					ts.NotEmpty(app.UpdatedAt)
				}

			}
		}
	}
}

// func (ts *PGDriverTestSuite) Test_UpdateApplication() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{
// 			name: "Should succeed without any errors",
// 			err:  nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		fmt.Println("RUNING TEST SUITE", test.name)
// 	}
// }

// func (ts *PGDriverTestSuite) Test_UpdateAppFirstDateSurpassed() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{
// 			name: "Should succeed without any errors",
// 			err:  nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		fmt.Println("RUNING TEST SUITE", test.name)
// 	}
// }

// func (ts *PGDriverTestSuite) Test_RemoveApp() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{
// 			name: "Should succeed without any errors",
// 			err:  nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		fmt.Println("RUNING TEST SUITE", test.name)
// 	}
// }

// func (ts *PGDriverTestSuite) Test_WriteBlockchain() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{
// 			name: "Should succeed without any errors",
// 			err:  nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		fmt.Println("RUNING TEST SUITE", test.name)
// 	}
// }

// func (ts *PGDriverTestSuite) Test_WriteRedirect() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{
// 			name: "Should succeed without any errors",
// 			err:  nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		fmt.Println("RUNING TEST SUITE", test.name)
// 	}
// }

// func (ts *PGDriverTestSuite) Test_ActivateBlockchain() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{
// 			name: "Should succeed without any errors",
// 			err:  nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		fmt.Println("RUNING TEST SUITE", test.name)
// 	}
// }
