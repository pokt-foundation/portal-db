package postgresdriver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

const connectionString = "postgres://postgres:pgpassword@localhost:5432/postgres?sslmode=disable"

// var (
// 	ctx = context.Background()
// )

func Test_RunPGDriverSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end to end test")
	}

	testSuite := new(PGDriverTestSuite)
	testSuite.connectionString = connectionString
	suite.Run(t, testSuite)
}

func (ts *PGDriverTestSuite) Test_ReadPayPlans() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_ReadApplications() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_ReadLoadBalancers() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_ReadBlockchains() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_WriteLoadBalancer() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_UpdateLoadBalancer() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_RemoveLoadBalancer() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_WriteApplication() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_UpdateApplication() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_UpdateAppFirstDateSurpassed() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_RemoveApp() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_WriteBlockchain() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_WriteRedirect() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}

func (ts *PGDriverTestSuite) Test_ActivateBlockchain() {
	// c := require.New(t)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Should succeed without any errors",
			err:  nil,
		},
	}

	for _, test := range tests {
		fmt.Println("RUNING TEST SUITE", test.name)
	}
}
