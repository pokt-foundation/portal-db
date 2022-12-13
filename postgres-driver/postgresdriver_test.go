package postgresdriver

import (
	"context"
	"testing"

	"github.com/pokt-foundation/portal-db/repository"
	"github.com/stretchr/testify/suite"
)

var ctx = context.Background()

func Test_RunPGDriverSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end to end test")
	}

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
	// c := require.New(t)

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
			},
			err: nil,
		},
	}

	for _, test := range tests {
		loadBalancers, err := ts.driver.ReadLoadBalancers(ctx)
		ts.Equal(test.err, err)
		for i, app := range loadBalancers {
			ts.Equal(test.loadBalancers[i].ID, app.ID)
			ts.Equal(test.loadBalancers[i].UserID, app.UserID)
			ts.Equal(test.loadBalancers[i].Name, app.Name)
			// ts.Equal(test.loadBalancers[i].URL, app.URL)
			// ts.Equal(test.loadBalancers[i].Dummy, app.Dummy)
			// ts.Equal(test.loadBalancers[i].Status, app.Status)
			// ts.Equal(test.loadBalancers[i].GatewayAAT, app.GatewayAAT)
			// ts.Equal(test.loadBalancers[i].GatewaySettings, app.GatewaySettings)
			// ts.Equal(test.loadBalancers[i].Limit, app.Limit)
			// ts.Equal(test.loadBalancers[i].NotificationSettings, app.NotificationSettings)
			// ts.NotEmpty(app.CreatedAt)
			// ts.NotEmpty(app.UpdatedAt)
		}
	}
}

// func (ts *PGDriverTestSuite) Test_ReadBlockchains() {
// 	// c := require.New(t)

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

// func (ts *PGDriverTestSuite) Test_WriteLoadBalancer() {
// 	// c := require.New(t)

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

// func (ts *PGDriverTestSuite) Test_UpdateLoadBalancer() {
// 	// c := require.New(t)

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

// func (ts *PGDriverTestSuite) Test_RemoveLoadBalancer() {
// 	// c := require.New(t)

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

// func (ts *PGDriverTestSuite) Test_WriteApplication() {
// 	// c := require.New(t)

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

// func (ts *PGDriverTestSuite) Test_UpdateApplication() {
// 	// c := require.New(t)

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
// 	// c := require.New(t)

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
// 	// c := require.New(t)

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
// 	// c := require.New(t)

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
// 	// c := require.New(t)

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
// 	// c := require.New(t)

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
