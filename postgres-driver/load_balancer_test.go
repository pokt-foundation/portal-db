package postgresdriver

import (
	"database/sql"

	"github.com/pokt-foundation/portal-db/types"
)

func (ts *PGDriverTestSuite) Test_ReadLoadBalancers() {
	tests := []struct {
		name          string
		loadBalancers []*types.LoadBalancer
		err           error
	}{
		{
			name: "Should return all Load Balancers from the database ordered by lb_id",
			loadBalancers: []*types.LoadBalancer{
				{
					ID:                "test_lb_34987u329rfn23f",
					Name:              "pokt_app_123",
					UserID:            "test_user_1dbffbdfeeb225",
					ApplicationIDs:    []string{"test_app_47hfnths73j2se"},
					RequestTimeout:    5_000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: types.StickyOptions{
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
					StickyOptions: types.StickyOptions{
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
					StickyOptions: types.StickyOptions{
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
		loadBalancers, err := ts.driver.ReadLoadBalancers(testCtx)
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

func (ts *PGDriverTestSuite) Test_WriteLoadBalancer() {
	tests := []struct {
		name               string
		loadBalancerInputs []*types.LoadBalancer
		expectedNumOfLBs   int
		expectedLB         SelectOneLoadBalancerRow
		err                error
	}{
		{
			name: "Should create a single load balancer successfully with correct input",
			loadBalancerInputs: []*types.LoadBalancer{
				{
					Name:              "pokt_app_789",
					UserID:            "test_user_47fhsd75jd756sh",
					RequestTimeout:    5000,
					Gigastake:         true,
					GigastakeRedirect: true,
					ApplicationIDs:    []string{"test_app_47hfnths73j2se"},
					StickyOptions: types.StickyOptions{
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
			createdLB, err := ts.driver.WriteLoadBalancer(testCtx, input)
			ts.Equal(test.err, err)
			ts.Len(createdLB.ID, 24)
			ts.Equal(input.Name, createdLB.Name)

			loadBalancers, err := ts.driver.ReadLoadBalancers(testCtx)
			ts.Equal(test.err, err)
			ts.Len(loadBalancers, test.expectedNumOfLBs)

			loadBalancer, err := ts.driver.SelectOneLoadBalancer(testCtx, createdLB.ID)
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
		loadBalancerUpdate  *types.UpdateLoadBalancer
		expectedAfterUpdate SelectOneLoadBalancerRow
		err                 error
	}{
		{
			name:           "Should update a single load balancer successfully with all fields",
			loadBalancerID: "test_lb_34987u329rfn23f",
			loadBalancerUpdate: &types.UpdateLoadBalancer{
				Name: "pokt_app_updated",
				StickyOptions: types.UpdateStickyOptions{
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
			loadBalancerUpdate: &types.UpdateLoadBalancer{
				Name: "pokt_app_updated_2",
				StickyOptions: types.UpdateStickyOptions{
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
			loadBalancerUpdate: &types.UpdateLoadBalancer{
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
			loadBalancerUpdate: &types.UpdateLoadBalancer{
				StickyOptions: types.UpdateStickyOptions{
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
		_, err := ts.driver.SelectOneLoadBalancer(testCtx, test.loadBalancerID)
		ts.Equal(test.err, err)

		err = ts.driver.UpdateLoadBalancer(testCtx, test.loadBalancerID, test.loadBalancerUpdate)
		ts.Equal(test.err, err)

		lbAfterUpdate, err := ts.driver.SelectOneLoadBalancer(testCtx, test.loadBalancerID)
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
		err := ts.driver.RemoveLoadBalancer(testCtx, test.loadBalancerID)
		ts.Equal(test.err, err)

		lbAfterRemove, err := ts.driver.SelectOneLoadBalancer(testCtx, test.loadBalancerID)
		ts.Equal(test.err, err)
		ts.Empty(lbAfterRemove.UserID.String)
	}
}
