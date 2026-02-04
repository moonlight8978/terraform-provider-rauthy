package rauthy_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
	"github.com/stretchr/testify/assert"
)

var roleResponse = `{
	"id": "role-1",
	"name": "Role 1"
}`

var rolesResponse = `[
	{
		"id": "role-1",
		"name": "Role 1"
	},
	{
		"id": "role-2",
		"name": "Role 2"
	}
]`

func TestCreateRole(t *testing.T) {
	ts := CreateServer(roleResponse, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	createdRole, err := client.CreateRole(context.Background(), &rauthy.RoleRequest{Role: "Role 1"})
	assert.NoError(t, err)
	assert.Equal(t, "role-1", createdRole.Id)
	assert.Equal(t, "Role 1", createdRole.Name)
}

func TestGetRoles(t *testing.T) {
	ts := CreateServer(rolesResponse, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	roles, err := client.GetRoles(context.Background())
	assert.NoError(t, err)
	assert.Len(t, roles, 2)
	assert.Equal(t, "role-1", roles[0].Id)
	assert.Equal(t, "role-2", roles[1].Id)
}

func TestGetRole(t *testing.T) {
	ts := CreateServer(rolesResponse, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	role, err := client.GetRole(context.Background(), "role-1")
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "role-1", role.Id)

	role, err = client.GetRole(context.Background(), "role-3")
	assert.Error(t, err)
	assert.Nil(t, role)
	assert.Contains(t, err.Error(), "role role-3 not found")
}

func TestUpdateRole(t *testing.T) {
	ts := CreateServer(roleResponse, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	updatedRole, err := client.UpdateRole(context.Background(), "role-1", &rauthy.RoleRequest{Role: "Role 1 Updated"})
	assert.NoError(t, err)
	assert.Equal(t, "role-1", updatedRole.Id)
	assert.Equal(t, "Role 1", updatedRole.Name)
}

func TestDeleteRole(t *testing.T) {
	ts := CreateServer("", http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	err := client.DeleteRole(context.Background(), "role-1")
	assert.NoError(t, err)
}
