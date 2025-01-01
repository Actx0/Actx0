// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/ai"
)

const bytesPerGB = 1024 * 1024 * 1024

// Usage provides workspace consumption helpers.
type Usage struct{}

// WorkspaceUsageMetrics is current workspace consumption for billing.
type WorkspaceUsageMetrics struct {
	APICalls         int64   `json:"apiCalls"`
	WorkspaceMembers int64   `json:"workspaceMembers"`
	DocumentsCount   int64   `json:"documentsCount"`
	PromptsCount     int64   `json:"promptsCount"`
	StorageGB        float64 `json:"storageGB"`
	HistoryRecords   int64   `json:"historyRecords"`
	MemoryRecords    int64   `json:"memoryRecords"`
	AITokens         int64   `json:"aiTokens"`
}

// UsageSnapshotDeps holds repositories needed to load workspace usage.
type UsageSnapshotDeps struct {
	WorkspaceUserRepository     db.WorkspaceUserRepository
	WorkspaceDocumentRepository db.WorkspaceDocumentRepository
	PromptRepository            db.PromptRepository
	UsageRepository             db.UsageRepository
	SubscriptionRepository      db.SubscriptionRepository
	SessionMessageRepository    db.SessionMessageRepository
	SessionMemoryRepository     db.SessionMemoryRepository
}

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

// APICalls returns API call usage for the current calendar month.
func (u *Usage) APICalls(ur db.UsageRepository, workspaceId db.Id) (int64, error) {
	start, _ := currentMonthPeriod()
	return ur.GetQuantityByPeriod(workspaceId, db.UsageTypeAPICalls, start)
}

// PromptsCount returns the number of prompts in a workspace.
func (u *Usage) PromptsCount(r db.PromptRepository, workspaceId db.Id) (int64, error) {
	return r.CountByWorkspaceId(workspaceId)
}

// HistoryRecords returns session message rows stored across a workspace.
func (u *Usage) HistoryRecords(r db.SessionMessageRepository, workspaceId db.Id) (int64, error) {
	return r.CountByWorkspaceId(workspaceId)
}

// MemoryRecords returns session memory rows stored across a workspace.
func (u *Usage) MemoryRecords(r db.SessionMemoryRepository, workspaceId db.Id) (int64, error) {
	return r.CountByWorkspaceId(workspaceId)
}

// GetWorkspaceUsage returns all billing usage metrics for a workspace.
func (u *Usage) GetWorkspaceUsage(deps UsageSnapshotDeps, workspaceId db.Id) (*WorkspaceUsageMetrics, error) {
	apiCalls, err := u.APICalls(deps.UsageRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	members, err := u.MembersCount(deps.WorkspaceUserRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	documents, err := u.DocumentsCount(deps.WorkspaceDocumentRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	prompts, err := u.PromptsCount(deps.PromptRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	storageBytes, err := u.StorageUsed(deps.WorkspaceDocumentRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	historyRecords, err := u.HistoryRecords(deps.SessionMessageRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	memoryRecords, err := u.MemoryRecords(deps.SessionMemoryRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	aiTokens, err := u.AITokens(deps.UsageRepository, deps.SubscriptionRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	return &WorkspaceUsageMetrics{
		APICalls:         apiCalls,
		WorkspaceMembers: members,
		DocumentsCount:   documents,
		PromptsCount:     prompts,
		StorageGB:        float64(storageBytes) / bytesPerGB,
		HistoryRecords:   historyRecords,
		MemoryRecords:    memoryRecords,
		AITokens:         aiTokens,
	}, nil
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
