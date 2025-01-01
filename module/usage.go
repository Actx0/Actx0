// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"errors"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/ai"
	"github.com/clivern/actx0/pkg/util"
)

const BytesPerGB = 1024 * 1024 * 1024

var ErrWorkspaceCapacityReached = errors.New("workspace capacity reached")

// Usage provides workspace consumption helpers.
type Usage struct{}

// WorkspaceUsageMetrics is current workspace consumption for billing.
type WorkspaceUsageMetrics struct {
	APICalls         int64   `json:"apiCalls"`
	WorkspaceMembers int64   `json:"workspaceMembers"`
	DocumentsCount   int64   `json:"documentsCount"`
	PromptsCount     int64   `json:"promptsCount"`
	AgentsCount      int64   `json:"agentsCount"`
	StorageGB        float64 `json:"storageGB"`
	SessionsCount    int64   `json:"sessionsCount"`
	HistoryRecords   int64   `json:"historyRecords"`
	MemoryRecords    int64   `json:"memoryRecords"`
	AITokens         int64   `json:"aiTokens"`
	AICost           float64 `json:"aiCost"`
}

// UsageSnapshotDeps holds repositories needed to load workspace usage.
type UsageSnapshotDeps struct {
	WorkspaceUserRepository     db.WorkspaceUserRepository
	WorkspaceDocumentRepository db.WorkspaceDocumentRepository
	PromptRepository            db.PromptRepository
	AgentRepository             db.AgentRepository
	UsageRepository             db.UsageRepository
	SubscriptionRepository      db.SubscriptionRepository
	SessionMessageRepository    db.SessionMessageRepository
	SessionMemoryRepository     db.SessionMemoryRepository
	AgentSessionRepository      db.AgentSessionRepository
}

// NewUsage returns a usage module.
func NewUsage() *Usage {
	return &Usage{}
}

// MembersCount returns the number of members in a workspace.
func (u *Usage) MembersCount(workspaceUsers db.WorkspaceUserRepository, workspaceId db.Id) (int64, error) {
	return workspaceUsers.CountByWorkspaceId(workspaceId)
}

// IsMembersCountReached reports whether adding members would exceed the plan limit.
func (u *Usage) IsMembersCountReached(workspaceUsers db.WorkspaceUserRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additional int64) (bool, error) {
	used, err := u.MembersCount(workspaceUsers, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).WorkspaceMembers
	if limit == 0 {
		return false, nil
	}

	return used+additional > limit, nil
}

// DocumentsCount returns the number of documents in a workspace.
func (u *Usage) DocumentsCount(documents db.WorkspaceDocumentRepository, workspaceId db.Id) (int64, error) {
	return documents.CountByWorkspaceId(workspaceId)
}

// IsDocumentsCountReached reports whether adding documents would exceed the plan limit.
func (u *Usage) IsDocumentsCountReached(documents db.WorkspaceDocumentRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additional int64) (bool, error) {
	used, err := u.DocumentsCount(documents, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).DocumentsCount
	if limit == 0 {
		return false, nil
	}

	return used+additional > limit, nil
}

// StorageUsed returns total document storage in bytes for a workspace.
func (u *Usage) StorageUsed(documents db.WorkspaceDocumentRepository, workspaceId db.Id) (int64, error) {
	return documents.SumSizeByWorkspaceId(workspaceId)
}

// IsStorageGBReached reports whether adding storage would exceed the plan limit.
func (u *Usage) IsStorageGBReached(documents db.WorkspaceDocumentRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additionalBytes int64) (bool, error) {
	usedBytes, err := u.StorageUsed(documents, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).StorageGB
	if limit == 0 {
		return false, nil
	}

	usedGB := float64(usedBytes) / BytesPerGB
	additionalGB := float64(additionalBytes) / BytesPerGB

	return usedGB+additionalGB > limit, nil
}

// APICalls returns API call usage for the current calendar month.
func (u *Usage) APICalls(usage db.UsageRepository, workspaceId db.Id) (int64, error) {
	start, _ := util.CurrentMonthPeriod()
	return usage.GetQuantityByPeriod(workspaceId, db.UsageTypeAPICalls, start)
}

// IsAPICallsReached reports whether adding API calls would exceed the plan limit.
func (u *Usage) IsAPICallsReached(usage db.UsageRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additional int64) (bool, error) {
	used, err := u.APICalls(usage, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).APICalls
	if limit == 0 {
		return false, nil
	}

	return used+additional > limit, nil
}

// PromptsCount returns the number of prompts in a workspace.
func (u *Usage) PromptsCount(prompts db.PromptRepository, workspaceId db.Id) (int64, error) {
	return prompts.CountByWorkspaceId(workspaceId)
}

// IsPromptsCountReached reports whether adding prompts would exceed the plan limit.
func (u *Usage) IsPromptsCountReached(prompts db.PromptRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additional int64) (bool, error) {
	used, err := u.PromptsCount(prompts, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).PromptsCount
	if limit == 0 {
		return false, nil
	}

	return used+additional > limit, nil
}

// AgentsCount returns the number of agents in a workspace.
func (u *Usage) AgentsCount(agents db.AgentRepository, workspaceId db.Id) (int64, error) {
	return agents.CountByWorkspaceId(workspaceId)
}

// IsAgentsCountReached reports whether adding agents would exceed the plan limit.
func (u *Usage) IsAgentsCountReached(agents db.AgentRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additional int64) (bool, error) {
	used, err := u.AgentsCount(agents, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).AgentsCount
	if limit == 0 {
		return false, nil
	}

	return used+additional > limit, nil
}

// HistoryRecords returns session message rows stored across a workspace.
func (u *Usage) HistoryRecords(messages db.SessionMessageRepository, workspaceId db.Id) (int64, error) {
	return messages.CountByWorkspaceId(workspaceId)
}

// SessionsCount returns agent session rows stored across a workspace.
func (u *Usage) SessionsCount(sessions db.AgentSessionRepository, workspaceId db.Id) (int64, error) {
	return sessions.CountByWorkspaceId(workspaceId)
}

// IsSessionsCountReached reports whether adding sessions would exceed the plan limit.
func (u *Usage) IsSessionsCountReached(sessions db.AgentSessionRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additional int64) (bool, error) {
	used, err := u.SessionsCount(sessions, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).SessionsCount
	if limit == 0 {
		return false, nil
	}

	return used+additional > limit, nil
}

// IsHistoryRecordsReached reports whether adding messages would exceed the plan limit.
func (u *Usage) IsHistoryRecordsReached(messages db.SessionMessageRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additional int64) (bool, error) {
	used, err := u.HistoryRecords(messages, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).HistoryRecords
	if limit == 0 {
		return false, nil
	}

	return used+additional > limit, nil
}

// MemoryRecords returns session memory rows stored across a workspace.
func (u *Usage) MemoryRecords(memories db.SessionMemoryRepository, workspaceId db.Id) (int64, error) {
	return memories.CountByWorkspaceId(workspaceId)
}

// IsMemoryRecordsReached reports whether adding memories would exceed the plan limit.
func (u *Usage) IsMemoryRecordsReached(memories db.SessionMemoryRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additional int64) (bool, error) {
	used, err := u.MemoryRecords(memories, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).MemoryRecords
	if limit == 0 {
		return false, nil
	}

	return used+additional > limit, nil
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

	agents, err := u.AgentsCount(deps.AgentRepository, workspaceId)
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

	sessionsCount, err := u.SessionsCount(deps.AgentSessionRepository, workspaceId)
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

	aiCost, err := u.AICost(deps.UsageRepository, deps.SubscriptionRepository, workspaceId)
	if err != nil {
		return nil, err
	}

	return &WorkspaceUsageMetrics{
		APICalls:         apiCalls,
		WorkspaceMembers: members,
		DocumentsCount:   documents,
		PromptsCount:     prompts,
		AgentsCount:      agents,
		StorageGB:        float64(storageBytes) / BytesPerGB,
		SessionsCount:    sessionsCount,
		HistoryRecords:   historyRecords,
		MemoryRecords:    memoryRecords,
		AITokens:         aiTokens,
		AICost:           aiCost,
	}, nil
}

// AITokens returns AI token usage for the workspace subscription period.
func (u *Usage) AITokens(usage db.UsageRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id) (int64, error) {
	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return 0, err
	}
	if subscription == nil || subscription.CurrentPeriodStart == nil {
		return 0, nil
	}

	return usage.GetQuantityByPeriod(
		workspaceId,
		db.UsageTypeAITokens,
		subscription.CurrentPeriodStart.UTC(),
	)
}

// IsAITokensReached reports whether adding tokens would exceed the plan limit.
func (u *Usage) IsAITokensReached(usage db.UsageRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additional int64) (bool, error) {
	used, err := u.AITokens(usage, subscriptions, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).AITokens
	if limit == 0 {
		return false, nil
	}

	return used+additional > limit, nil
}

// AICost returns AI spend in USD for the workspace subscription period.
func (u *Usage) AICost(usage db.UsageRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id) (float64, error) {
	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return 0, err
	}
	if subscription == nil || subscription.CurrentPeriodStart == nil {
		return 0, nil
	}

	cost, err := usage.GetQuantityByPeriod(
		workspaceId,
		db.UsageTypeAICost,
		subscription.CurrentPeriodStart.UTC(),
	)
	if err != nil {
		return 0, err
	}

	return float64(cost) / ai.USDToNano, nil
}

// IsAICostReached reports whether adding AI spend would exceed the plan limit.
func (u *Usage) IsAICostReached(usage db.UsageRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, additionalUSD float64) (bool, error) {
	used, err := u.AICost(usage, subscriptions, workspaceId)
	if err != nil {
		return false, err
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	if subscription == nil {
		return false, ErrBillingSubscriptionNotFound
	}

	limit := GetPlanLimits(subscription.Plan).AICost
	if limit == 0 {
		return false, nil
	}

	return used+additionalUSD > limit, nil
}

// IncrementAPICalls records one successful workspace API call for the current month.
func (u *Usage) IncrementAPICalls(usage db.UsageRepository, workspaceId db.Id) error {
	start, end := util.CurrentMonthPeriod()

	return usage.IncrementByPeriod(
		workspaceId,
		db.UsageTypeAPICalls,
		start,
		end,
		1,
		db.UsageUnitCalls,
	)
}

// IncrementAIUsage adds AI token and cost usage for the workspace subscription period.
func (u *Usage) IncrementAIUsage(usage db.UsageRepository, subscriptions db.SubscriptionRepository, workspaceId db.Id, tokens, cost int64) error {
	if tokens == 0 && cost == 0 {
		return nil
	}

	subscription, err := subscriptions.GetByWorkspaceId(workspaceId)
	if err != nil {
		return err
	}
	if subscription == nil || subscription.CurrentPeriodStart == nil || subscription.CurrentPeriodEnd == nil {
		return nil
	}

	pstart := subscription.CurrentPeriodStart.UTC()
	pend := subscription.CurrentPeriodEnd.UTC()

	if tokens != 0 {
		if err := usage.IncrementByPeriod(
			workspaceId,
			db.UsageTypeAITokens,
			pstart,
			pend,
			tokens,
			db.UsageUnitTokens,
		); err != nil {
			return err
		}
	}

	if cost != 0 {
		if err := usage.IncrementByPeriod(
			workspaceId,
			db.UsageTypeAICost,
			pstart,
			pend,
			cost,
			db.UsageUnitNanoUSD,
		); err != nil {
			return err
		}
	}

	return nil
}
