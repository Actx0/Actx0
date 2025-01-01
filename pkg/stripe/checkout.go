// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package stripe

import (
	"context"
	"fmt"

	stripesdk "github.com/stripe/stripe-go/v82"
)

// CheckoutOptions is input for creating a subscription Checkout Session.
type CheckoutOptions struct {
	WorkspaceId   string
	Plan          string
	CustomerId    string
	CustomerEmail string
	SuccessURL    string
	CancelURL     string
}

// CreateCheckoutSession starts a Stripe Checkout subscription flow.
func (c *Client) CreateCheckoutSession(ctx context.Context, opts CheckoutOptions) (*stripesdk.CheckoutSession, error) {
	priceId := c.config.PriceId(opts.Plan)
	if priceId == "" {
		return nil, ErrInvalidPlan
	}

	params := &stripesdk.CheckoutSessionCreateParams{
		Mode:       stripesdk.String(string(stripesdk.CheckoutSessionModeSubscription)),
		SuccessURL: stripesdk.String(opts.SuccessURL),
		CancelURL:  stripesdk.String(opts.CancelURL),
		LineItems: []*stripesdk.CheckoutSessionCreateLineItemParams{
			{
				Price:    stripesdk.String(priceId),
				Quantity: stripesdk.Int64(1),
			},
		},
		Metadata: map[string]string{
			"workspace_id": opts.WorkspaceId,
			"plan":         opts.Plan,
		},
		SubscriptionData: &stripesdk.CheckoutSessionCreateSubscriptionDataParams{
			Metadata: map[string]string{
				"workspace_id": opts.WorkspaceId,
				"plan":         opts.Plan,
			},
		},
	}

	if opts.CustomerId != "" {
		params.Customer = stripesdk.String(opts.CustomerId)
	} else if opts.CustomerEmail != "" {
		params.CustomerEmail = stripesdk.String(opts.CustomerEmail)
	}

	session, err := c.api.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create checkout session: %w", err)
	}
	return session, nil
}
