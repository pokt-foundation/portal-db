package postgresdriver

import (
	"context"

	"github.com/pokt-foundation/portal-db/types"
)

/* ReadPayPlans returns all pay plans in the database and marshals to repository struct */
func (q *Queries) ReadPayPlans(ctx context.Context) ([]*types.PayPlan, error) {
	dbPayPlans, err := q.SelectPayPlans(ctx)
	if err != nil {
		return nil, err
	}

	var payPlans []*types.PayPlan

	for _, dbPayPlan := range dbPayPlans {
		payPlan, err := dbPayPlan.toPayPlan()
		if err != nil {
			return nil, err
		}

		payPlans = append(payPlans, payPlan)
	}

	return payPlans, nil
}

func (p *SelectPayPlansRow) toPayPlan() (*types.PayPlan, error) {
	payPlan := types.PayPlan{
		Type:  types.PayPlanType(p.PlanType),
		Limit: int(p.DailyLimit),
	}

	err := payPlan.Validate()
	if err != nil {
		return nil, err
	}

	return &payPlan, nil
}
