package rauthy

import (
	"context"
	"fmt"
	"net/http"
)

type PasswordPolicy struct {
	LengthMin        int `json:"length_min"`
	LengthMax        int `json:"length_max"`
	IncludeDigits    int `json:"include_digits"`
	IncludeLowerCase int `json:"include_lower_case"`
	IncludeUpperCase int `json:"include_upper_case"`
	IncludeSpecial   int `json:"include_special"`
	// NotRecentlyUsed  int `json:"not_recently_used,omitempty"`
	ValidDays int `json:"valid_days,omitempty"`
}

func (c *Client) GetPasswordPolicy(ctx context.Context) (*PasswordPolicy, error) {
	var passwordPolicy PasswordPolicy
	if _, err := c.Request(ctx, http.MethodGet, "/password_policy", nil, &passwordPolicy); err != nil {
		return nil, fmt.Errorf("Failed to get password policy - Reason: %w", err)
	}

	return &passwordPolicy, nil
}

func (c *Client) UpdatePasswordPolicy(ctx context.Context, passwordPolicy *PasswordPolicy) (*PasswordPolicy, error) {
	var updatedPasswordPolicy PasswordPolicy
	if _, err := c.Request(ctx, http.MethodPut, "/password_policy", passwordPolicy, &updatedPasswordPolicy); err != nil {
		return nil, fmt.Errorf("Failed to update password policy - Reason: %w", err)
	}

	return &updatedPasswordPolicy, nil
}
