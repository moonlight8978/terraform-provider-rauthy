package rauthy

import (
	"context"
	"fmt"
	"net/http"
)

type Group struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GroupRequest struct {
	Group string `json:"group"`
}

func (c *Client) GetGroups(ctx context.Context) ([]Group, error) {
	var groups []Group
	_, err := c.Request(ctx, http.MethodGet, "/groups", nil, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (c *Client) CreateGroup(ctx context.Context, group *GroupRequest) (*Group, error) {
	var newGroup Group
	_, err := c.Request(ctx, http.MethodPost, "/groups", group, &newGroup)
	if err != nil {
		return nil, err
	}
	return &newGroup, nil
}

func (c *Client) UpdateGroup(ctx context.Context, id string, group *GroupRequest) (*Group, error) {
	var updatedGroup Group
	_, err := c.Request(ctx, http.MethodPut, fmt.Sprintf("/groups/%s", id), group, &updatedGroup)
	if err != nil {
		return nil, err
	}
	return &updatedGroup, nil
}

func (c *Client) DeleteGroup(ctx context.Context, id string) error {
	_, err := c.Request(ctx, http.MethodDelete, fmt.Sprintf("/groups/%s", id), nil, nil)
	return err
}
