package driver

import (
	"context"

	"github.com/pokt-foundation/portal-db/repository"
)

type (
	// The Driver interface represents all database operations required by the Pocket HTTP DB
	Driver interface {
		Reader
		Writer
	}

	Reader interface {
		ReadPayPlans(ctx context.Context) ([]*repository.PayPlan, error)
		ReadApplications(ctx context.Context) ([]*repository.Application, error)
		ReadLoadBalancers(ctx context.Context) ([]*repository.LoadBalancer, error)
		ReadBlockchains(ctx context.Context) ([]*repository.Blockchain, error)

		NotificationChannel() <-chan *repository.Notification
	}

	Writer interface {
		WriteLoadBalancer(ctx context.Context, loadBalancer *repository.LoadBalancer) (*repository.LoadBalancer, error)
		UpdateLoadBalancer(ctx context.Context, id string, options *repository.UpdateLoadBalancer) error
		RemoveLoadBalancer(ctx context.Context, id string) error

		WriteApplication(ctx context.Context, app *repository.Application) (*repository.Application, error)
		UpdateApplication(ctx context.Context, id string, update *repository.UpdateApplication) error
		UpdateAppFirstDateSurpassed(ctx context.Context, update *repository.UpdateFirstDateSurpassed) error
		RemoveApp(ctx context.Context, id string) error

		WriteBlockchain(ctx context.Context, blockchain *repository.Blockchain) (*repository.Blockchain, error)
		WriteRedirect(ctx context.Context, redirect *repository.Redirect) (*repository.Redirect, error)
		ActivateChain(ctx context.Context, id string, active bool) error
	}
)
