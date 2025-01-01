// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/ai"
)

// Usage provides workspace consumption helpers.
type Usage struct{}

// NewUsage returns a usage module.
func NewUsage() *Usage {
	return &Usage{}
}

// MembersCount returns the number of members in a workspace.
func (u *Usage) MembersCount(r db.WorkspaceUserRepository, workspaceId db.Id) (int64, error) {
	return r.CountByWorkspaceId(workspaceId)
}

// DocumentsCount returns the number of documents in a workspace.
func (u *Usage) DocumentsCount(r db.WorkspaceDocumentRepository, workspaceId db.Id) (int64, error) {
	return r.CountByWorkspaceId(workspaceId)
}

// StorageUsed returns total document storage in bytes for a workspace.
func (u *Usage) StorageUsed(r db.WorkspaceDocumentRepository, workspaceId db.Id) (int64, error) {
	return r.SumSizeByWorkspaceId(workspaceId)
}

// AITokens returns AI token usage for the workspace subscription period.
func (u *Usage) AITokens(ur db.UsageRepository, sr db.SubscriptionRepository, workspaceId db.Id) (int64, error) {
	subscription, err := sr.GetByWorkspaceId(workspaceId)
	if err != nil {
		return 0, err
	}
	if subscription == nil || subscription.CurrentPeriodStart == nil {
		return 0, nil
	}

	return ur.GetQuantityByPeriod(
		workspaceId,
		db.UsageTypeAITokens,
		subscription.CurrentPeriodStart.UTC(),
	)
}

// AICost returns AI spend in USD for the workspace subscription period.
func (u *Usage) AICost(ur db.UsageRepository, sr db.SubscriptionRepository, workspaceId db.Id) (float64, error) {
	subscription, err := sr.GetByWorkspaceId(workspaceId)
	if err != nil {
		return 0, err
	}
	if subscription == nil || subscription.CurrentPeriodStart == nil {
		return 0, nil
	}

	cost, err := ur.GetQuantityByPeriod(
		workspaceId,
		db.UsageTypeAICost,
		subscription.CurrentPeriodStart.UTC(),
	)
	if err != nil {
		return 0, err
	}

	return float64(cost) / ai.USDToNano, nil
}

func currentMonthPeriod() (time.Time, time.Time) {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	return start, start.AddDate(0, 1, 0)
}

// IncrementAPICalls records one successful workspace API call for the current month.
func (u *Usage) IncrementAPICalls(ur db.UsageRepository, workspaceId db.Id) error {
	start, end := currentMonthPeriod()
	return ur.IncrementByPeriod(
		workspaceId,
		db.UsageTypeAPICalls,
		start,
		end,
		1,
		db.UsageUnitCalls,
	)
}

// IncrementAIUsage adds AI token and cost usage for the workspace subscription period.
func (u *Usage) IncrementAIUsage(ur db.UsageRepository, sr db.SubscriptionRepository, workspaceId db.Id, tokens, cost int64) error {
	if tokens == 0 && cost == 0 {
		return nil
	}

	subscription, err := sr.GetByWorkspaceId(workspaceId)
	if err != nil {
		return err
	}
	if subscription == nil || subscription.CurrentPeriodStart == nil || subscription.CurrentPeriodEnd == nil {
		return nil
	}

	periodStart := subscription.CurrentPeriodStart.UTC()
	periodEnd := subscription.CurrentPeriodEnd.UTC()

	if tokens != 0 {
		if err := ur.IncrementByPeriod(
			workspaceId,
			db.UsageTypeAITokens,
			periodStart,
			periodEnd,
			tokens,
			db.UsageUnitTokens,
		); err != nil {
			return err
		}
	}

	if cost != 0 {
		if err := ur.IncrementByPeriod(
			workspaceId,
			db.UsageTypeAICost,
			periodStart,
			periodEnd,
			cost,
			db.UsageUnitNanoUSD,
		); err != nil {
			return err
		}
	}

	return nil
}
