package source

import (
	"context"

	"github.com/pokt-foundation/portal-db/repository"
)

// The Source interface represents all database operations required by the Pocket HTTP DB
type Source interface {
	/* Applications Table */
	ReadApplications(ctx context.Context) ([]*repository.Application, error)
	WriteLoadBalancer(ctx context.Context, loadBalancer *repository.LoadBalancer) (*repository.LoadBalancer, error)
	UpdateLoadBalancer(ctx context.Context, id string, options *repository.UpdateLoadBalancer) error
	RemoveLoadBalancer(ctx context.Context, id string) error
	/* Load Balancers Table */
	ReadLoadBalancers(ctx context.Context) ([]*repository.LoadBalancer, error)
	ReadPayPlans(ctx context.Context) ([]*repository.PayPlan, error)
	WriteApplication(ctx context.Context, app *repository.Application) (*repository.Application, error)
	UpdateApplication(ctx context.Context, id string, options *repository.UpdateApplication) error
	UpdateFirstDateSurpassed(ctx context.Context, firstDateSurpassed *repository.UpdateFirstDateSurpassed) error
	RemoveApplication(ctx context.Context, id string) error
	/* Blockchains Table */
	ReadBlockchains(ctx context.Context) ([]*repository.Blockchain, error)
	WriteBlockchain(ctx context.Context, blockchain *repository.Blockchain) (*repository.Blockchain, error)
	WriteRedirect(ctx context.Context, redirect *repository.Redirect) (*repository.Redirect, error)
	ActivateBlockchain(ctx context.Context, id string, active bool) error
	/* Listener Channel */
	NotificationChannel() <-chan *repository.Notification
}
