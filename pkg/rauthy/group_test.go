package rauthy_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
	"github.com/stretchr/testify/assert"
)

var groupResponse = `{
	"id": "group-1",
	"name": "Group 1"
}`

var groupsResponse = `[
	{
		"id": "group-1",
		"name": "Group 1"
	},
	{
		"id": "group-2",
		"name": "Group 2"
	}
]`

func TestCreateGroup(t *testing.T) {
	ts := CreateServer(groupResponse, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	createdGroup, err := client.CreateGroup(context.Background(), &rauthy.GroupRequest{Group: "Group 1"})
	assert.NoError(t, err)
	assert.Equal(t, "group-1", createdGroup.Id)
	assert.Equal(t, "Group 1", createdGroup.Name)
}

func TestGetGroups(t *testing.T) {
	ts := CreateServer(groupsResponse, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	groups, err := client.GetGroups(context.Background())
	assert.NoError(t, err)
	assert.Len(t, groups, 2)
	assert.Equal(t, "group-1", groups[0].Id)
	assert.Equal(t, "group-2", groups[1].Id)
}

func TestUpdateGroup(t *testing.T) {
	ts := CreateServer(groupResponse, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	updatedGroup, err := client.UpdateGroup(context.Background(), "group-1", &rauthy.GroupRequest{Group: "Group 1 Updated"})
	assert.NoError(t, err)
	assert.Equal(t, "group-1", updatedGroup.Id)
	assert.Equal(t, "Group 1", updatedGroup.Name)
}

func TestDeleteGroup(t *testing.T) {
	ts := CreateServer("", http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	err := client.DeleteGroup(context.Background(), "group-1")
	assert.NoError(t, err)
}
