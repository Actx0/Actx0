// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/stripe"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	stripesdk "github.com/stripe/stripe-go/v82"
)

const (
	BillingProviderStripe = "stripe"
	BillingStatusActive   = "active"
	BillingStatusPastDue  = "past_due"
	BillingStatusCanceled = "canceled"
)

var (
	ErrBillingPortalUnavailable    = errors.New("stripe billing portal is unavailable for this workspace")
	ErrBillingFreePlan             = errors.New("hobby is the default free plan and does not use checkout")
	ErrBillingSubscriptionNotFound = errors.New("workspace subscription not found")
)

// PlanUsageLimits are the consumption caps for a billing plan.
type PlanUsageLimits struct {
	Price            float64 `json:"price"`
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

var PlanLimits = map[string]PlanUsageLimits{
	stripe.PlanHobby: {
		Price:            0,
		APICalls:         10_000,
		WorkspaceMembers: 1,
		DocumentsCount:   10,
		PromptsCount:     10,
		AgentsCount:      2,
		StorageGB:        1,
		SessionsCount:    25,
		HistoryRecords:   1_000,
		MemoryRecords:    500,
		AITokens:         500_000,
		AICost:           2.0,
	},
	stripe.PlanStarter: {
		Price:            19,
		APICalls:         1_000_000,
		WorkspaceMembers: 10,
		DocumentsCount:   500,
		PromptsCount:     100,
		AgentsCount:      10,
		StorageGB:        5,
		SessionsCount:    250,
		HistoryRecords:   20_000,
		MemoryRecords:    10_000,
		AITokens:         2_000_000,
		AICost:           15.2,
	},
	stripe.PlanGrowth: {
		Price:            69,
		APICalls:         3_000_000,
		WorkspaceMembers: 25,
		DocumentsCount:   2_000,
		PromptsCount:     400,
		AgentsCount:      50,
		StorageGB:        25,
		SessionsCount:    1_000,
		HistoryRecords:   75_000,
		MemoryRecords:    40_000,
		AITokens:         7_300_000,
		AICost:           55.2,
	},
	stripe.PlanPro: {
		Price:            219,
		APICalls:         10_000_000,
		WorkspaceMembers: 100,
		DocumentsCount:   6_500,
		PromptsCount:     1_250,
		AgentsCount:      200,
		StorageGB:        100,
		SessionsCount:    5_000,
		HistoryRecords:   250_000,
		MemoryRecords:    125_000,
		AITokens:         23_200_000,
		AICost:           175.2,
	},
}

// Billing coordinates local subscription state with Stripe Checkout and Portal.
type Billing struct {
	WorkspaceRepository    db.WorkspaceRepository
	SubscriptionRepository db.SubscriptionRepository
	Usage                  Usage
}

// BillingStatusResponse is the billing state exposed to the UI.
type BillingStatusResponse struct {
	Id                 db.Id   `json:"id"`
	WorkspaceId        db.Id   `json:"workspaceId"`
	Plan               string  `json:"plan"`
	Status             string  `json:"status"`
	Provider           *string `json:"provider"`
	ProviderCustomerId string  `json:"providerCustomerId"`
	PeriodStart        string  `json:"periodStart"`
	PeriodEnd          string  `json:"periodEnd"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

// BillingCheckoutRequest is the body for creating a Stripe Checkout session.
type BillingCheckoutRequest struct {
	Plan string `json:"plan" validate:"required" label:"Plan"`
}

// BillingSessionResponse contains a Stripe redirect URL.
type BillingSessionResponse struct {
	URL string `json:"url"`
}

// BillingWebhookResponse is returned after a Stripe webhook is handled.
type BillingWebhookResponse struct {
	Received bool `json:"received"`
}

// BillingUsageResponse is workspace usage and limits for the billing page.
type BillingUsageResponse struct {
	Plan        string                `json:"plan"`
	PeriodReset string                `json:"periodReset"`
	Used        WorkspaceUsageMetrics `json:"used"`
	Limits      PlanUsageLimits       `json:"limits"`
}

// NewBilling creates a billing module.
func NewBilling(workspaces db.WorkspaceRepository, subscriptions db.SubscriptionRepository, usage Usage) *Billing {
	return &Billing{
		WorkspaceRepository:    workspaces,
		SubscriptionRepository: subscriptions,
		Usage:                  usage,
	}
}

// PlanLimits returns the default usage limits for a plan.
func GetPlanLimits(plan string) PlanUsageLimits {
	limits, ok := PlanLimits[plan]
	if !ok {
		return PlanLimits[stripe.PlanHobby]
	}

	return limits
}

// GetBillingStatus returns billing status for a workspace.
func (b *Billing) GetBillingStatus(workspaceId db.Id) (*BillingStatusResponse, error) {
	workspace, err := b.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrBillingSubscriptionNotFound
	}

	err = b.UpdateHobbyPlanPeriods(subscription)
	if err != nil {
		return nil, err
	}

	return &BillingStatusResponse{
		Id:                 subscription.Id,
		WorkspaceId:        subscription.WorkspaceId,
		Plan:               subscription.Plan,
		Status:             subscription.Status,
		Provider:           subscription.Provider,
		ProviderCustomerId: lo.FromPtr(subscription.ProviderCustomerId),
		PeriodStart:        subscription.CurrentPeriodStart.UTC().Format(time.RFC3339),
		PeriodEnd:          subscription.CurrentPeriodEnd.UTC().Format(time.RFC3339),
		CreatedAt:          subscription.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:          subscription.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetBillingUsage returns usage and plan limits for a workspace.
func (b *Billing) GetBillingUsage(workspaceId db.Id, deps UsageSnapshotDeps) (*BillingUsageResponse, error) {
	workspace, err := b.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrBillingSubscriptionNotFound
	}

	err = b.UpdateHobbyPlanPeriods(subscription)
	if err != nil {
		return nil, err
	}

	used, err := b.Usage.GetWorkspaceUsage(deps, workspaceId)
	if err != nil {
		return nil, err
	}

	limits := GetPlanLimits(subscription.Plan)

	return &BillingUsageResponse{
		Plan:        subscription.Plan,
		PeriodReset: subscription.CurrentPeriodEnd.UTC().Format(time.RFC3339),
		Used:        *used,
		Limits:      limits,
	}, nil
}

// CreateCheckoutSession starts Stripe Checkout for a paid plan.
func (b *Billing) CreateCheckoutSession(ctx context.Context, workspaceId db.Id, user *db.User, req *BillingCheckoutRequest, successURL string, cancelURL string) (*BillingSessionResponse, error) {
	workspace, err := b.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	if req.Plan == stripe.PlanHobby {
		return nil, ErrBillingFreePlan
	}

	client, err := stripe.New()
	if err != nil {
		return nil, err
	}
	if lo.IsEmpty(client.Config().PriceId(req.Plan)) {
		return nil, stripe.ErrInvalidPlan
	}

	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrBillingSubscriptionNotFound
	}

	session, err := client.CreateCheckoutSession(ctx, stripe.CheckoutOptions{
		WorkspaceId:   workspaceId.String(),
		Plan:          req.Plan,
		CustomerId:    lo.FromPtr(subscription.ProviderCustomerId),
		CustomerEmail: user.Email,
		SuccessURL:    successURL,
		CancelURL:     cancelURL,
	})
	if err != nil {
		return nil, err
	}

	return &BillingSessionResponse{URL: session.URL}, nil
}

// CreatePortalSession opens Stripe's customer portal.
func (b *Billing) CreatePortalSession(ctx context.Context, workspaceId db.Id, returnURL string) (*BillingSessionResponse, error) {
	workspace, err := b.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	client, err := stripe.New()
	if err != nil {
		return nil, err
	}
	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrBillingSubscriptionNotFound
	}
	customerId := lo.FromPtr(subscription.ProviderCustomerId)
	if lo.IsEmpty(customerId) {
		return nil, ErrBillingPortalUnavailable
	}

	session, err := client.CreatePortalSession(ctx, stripe.PortalOptions{
		CustomerId: customerId,
		ReturnURL:  returnURL,
	})
	if err != nil {
		return nil, err
	}

	return &BillingSessionResponse{URL: session.URL}, nil
}

// HandleWebhook updates local billing state from Stripe webhook events.
func (b *Billing) HandleWebhook(payload []byte, signature string) error {
	client, err := stripe.New()
	if err != nil {
		return err
	}

	event, err := client.ConstructWebhookEvent(payload, signature)
	if err != nil {
		return err
	}

	switch event.Type {
	case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted":
		var subscription stripesdk.Subscription
		err = json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return fmt.Errorf("decode subscription: %w", err)
		}

		return b.SyncStripeSubscription(&subscription, string(event.Type))
	case "invoice.paid":
		var invoice stripesdk.Invoice
		err = json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			return fmt.Errorf("decode invoice: %w", err)
		}

		period := invoice.Lines.Data[0].Period
		return b.SyncStripeSubscription(&stripesdk.Subscription{
			ID:       invoice.Parent.SubscriptionDetails.Subscription.ID,
			Customer: invoice.Customer,
			Metadata: invoice.Parent.SubscriptionDetails.Metadata,
			Status:   stripesdk.SubscriptionStatusActive,
			Items: &stripesdk.SubscriptionItemList{
				Data: []*stripesdk.SubscriptionItem{{
					CurrentPeriodStart: period.Start,
					CurrentPeriodEnd:   period.End,
				}},
			},
		}, string(event.Type))
	case "invoice.payment_failed":
		var invoice stripesdk.Invoice
		err = json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			return fmt.Errorf("decode invoice: %w", err)
		}

		return b.SyncStripeSubscription(&stripesdk.Subscription{
			ID:       invoice.Parent.SubscriptionDetails.Subscription.ID,
			Customer: invoice.Customer,
			Metadata: invoice.Parent.SubscriptionDetails.Metadata,
			Status:   stripesdk.SubscriptionStatusPastDue,
		}, string(event.Type))
	default:
		return nil
	}
}

// SyncStripeSubscription syncs a Stripe subscription with the local database.
func (b *Billing) SyncStripeSubscription(payload *stripesdk.Subscription, eventType string) error {
	subscription, err := b.SubscriptionRepository.GetByProviderSubscriptionId(payload.ID)
	if err != nil {
		return err
	}

	if subscription == nil {
		subscription, err = b.SubscriptionRepository.GetByWorkspaceId(db.Id(payload.Metadata["workspaceId"]))
		if err != nil {
			return err
		}
	}

	if subscription == nil {
		return nil
	}

	subscription.Provider = new(BillingProviderStripe)
	subscription.ProviderCustomerId = new(payload.Customer.ID)
	subscription.ProviderSubscriptionId = new(payload.ID)
	if plan := payload.Metadata["plan"]; lo.IsNotEmpty(plan) {
		subscription.Plan = plan
	}
	subscription.Status = lo.Ternary(eventType == "customer.subscription.deleted", BillingStatusCanceled, string(payload.Status))

	if eventType == "invoice.payment_failed" {
		// TODO: notify customers that invoice payment failed
	}

	// Update the subscription period if the event type is invoice.paid
	if eventType == "invoice.paid" {
		item := payload.Items.Data[0]
		subscription.CurrentPeriodStart = new(time.Unix(item.CurrentPeriodStart, 0).UTC())
		subscription.CurrentPeriodEnd = new(time.Unix(item.CurrentPeriodEnd, 0).UTC())
	}

	err = b.SubscriptionRepository.Update(subscription)
	if err != nil {
		return err
	}

	log.Info().
		Str("workspaceId", subscription.WorkspaceId.String()).
		Str("subscriptionId", payload.ID).
		Str("plan", subscription.Plan).
		Str("status", subscription.Status).
		Time("periodStart", lo.FromPtr(subscription.CurrentPeriodStart)).
		Time("periodEnd", lo.FromPtr(subscription.CurrentPeriodEnd)).
		Msg("Stripe subscription synced")

	return nil
}

// UpdateHobbyPlanPeriods updates the hobby plan periods for a subscription.
func (b *Billing) UpdateHobbyPlanPeriods(subscription *db.Subscription) error {
	if subscription == nil || subscription.Plan != stripe.PlanHobby {
		return nil
	}

	now := time.Now().UTC()
	if subscription.CurrentPeriodEnd != nil && !subscription.CurrentPeriodEnd.Before(now) {
		return nil
	}

	periodStart := subscription.CurrentPeriodEnd.UTC()
	periodEnd := periodStart.AddDate(0, 1, 0)

	for !periodEnd.After(now) {
		periodStart = periodEnd
		periodEnd = periodStart.AddDate(0, 1, 0)
	}

	subscription.CurrentPeriodStart = new(periodStart)
	subscription.CurrentPeriodEnd = new(periodEnd)

	err := b.SubscriptionRepository.Update(subscription)

	if err != nil {
		return err
	}
	return nil
}

// GetWorkspaceSubscription gets a workspace subscription from the database.
func (b *Billing) GetWorkspaceSubscription(workspaceId db.Id) (*db.Subscription, error) {
	subscription, err := b.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrBillingSubscriptionNotFound
	}
	return subscription, nil
}
