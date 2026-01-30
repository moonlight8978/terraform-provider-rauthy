package rauthy

import (
	"context"
	"fmt"
	"net/http"
)

type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) CreateRole(ctx context.Context, role *Role) (Role, error) {
	var createdRole Role

	if _, err := c.Request(ctx, http.MethodPost, "/roles", role, &createdRole); err != nil {
		return createdRole, err
	}

	return createdRole, nil
}

func (c *Client) GetRoles(ctx context.Context) ([]Role, error) {
	var roles []Role

	if _, err := c.Request(ctx, http.MethodGet, "/roles", nil, &roles); err != nil {
		return roles, err
	}

	return roles, nil
}

func (c *Client) GetRole(ctx context.Context, id string) (*Role, error) {
	roles, err := c.GetRoles(ctx)

	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.Id == id {
			return &role, nil
		}
	}

	return nil, fmt.Errorf("role %s not found", id)
}

func (c *Client) UpdateRole(ctx context.Context, id string, role *Role) (*Role, error) {
	var updatedRole Role

	if _, err := c.Request(ctx, http.MethodPut, fmt.Sprintf("/roles/%s", id), role, &updatedRole); err != nil {
		return nil, err
	}

	return &updatedRole, nil
}

func (c *Client) DeleteRole(ctx context.Context, id string) error {
	if _, err := c.Request(ctx, http.MethodDelete, fmt.Sprintf("/roles/%s", id), nil, nil); err != nil {
		return err
	}

	return nil
}
