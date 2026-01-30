package oidc_client

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/utils"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
)

var _ resource.Resource = &OidcClientSecretResource{}
var _ resource.ResourceWithImportState = &OidcClientSecretResource{}

func NewOidcClientSecretResource() resource.Resource {
	return &OidcClientSecretResource{}
}

type OidcClientSecretResource struct {
	client *rauthy.Client
}

func (r *OidcClientSecretResource) SetClient(c *rauthy.Client) {
	r.client = c
}

type OidcClientSecretResourceModel struct {
	Id                types.String `tfsdk:"id"`
	ClientId          types.String `tfsdk:"client_id"`
	CacheCurrentHours types.Int64  `tfsdk:"cache_current_hours"`
	Secret            types.String `tfsdk:"secret"`
}

func (r *OidcClientSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_secret"
}

func (r *OidcClientSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Client secret resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Client ID",
				Required:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID",
				Required:            true,
			},
			"cache_current_hours": schema.Int64Attribute{
				MarkdownDescription: "Cache current hours",
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(0),
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "Secret",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *OidcClientSecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	utils.ConfigureProvider(ctx, req, resp, r)
}

func (r *OidcClientSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model OidcClientSecretResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.CreateClientSecret(ctx, model.ClientId.ValueString(), &rauthy.ClientSecretRequest{
		CacheCurrentHours: int(model.CacheCurrentHours.ValueInt64()),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create client secret, got error: %s", err))
		return
	}

	model.Secret = types.StringValue(secret.Secret)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *OidcClientSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model OidcClientSecretResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *OidcClientSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model OidcClientSecretResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *OidcClientSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model OidcClientSecretResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *OidcClientSecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	clientId, id, found := strings.Cut(req.ID, "/")

	if !found {
		resp.Diagnostics.AddError("Client Error", "ID format is invalid, expected <client_id>/<secret_id>")
		return
	}

	_, err := r.client.GetOidcClient(ctx, clientId)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import client secret, got error: %s", err))
		return
	}

	model := &OidcClientSecretResourceModel{
		ClientId:          types.StringValue(clientId),
		Id:                types.StringValue(id),
		CacheCurrentHours: types.Int64Value(0),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
