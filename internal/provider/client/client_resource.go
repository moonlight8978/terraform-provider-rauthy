package client

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/tfutils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ClientResource{}
var _ resource.ResourceWithImportState = &ClientResource{}

func NewClientResource() resource.Resource {
	return &ClientResource{}
}

// ClientResource defines the resource implementation.
type ClientResource struct {
	client *rauthy.Client
}

// ClientResourceModel describes the resource data model.
type ClientResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Enabled                types.Bool   `tfsdk:"enabled"`
	Confidential           types.Bool   `tfsdk:"confidential"`
	RedirectUris           types.List   `tfsdk:"redirect_uris"`
	PostLogoutRedirectUris types.List   `tfsdk:"post_logout_redirect_uris"`
	FlowsEnabled           types.List   `tfsdk:"flows_enabled"`
	AccessTokenAlg         types.String `tfsdk:"access_token_alg"`
	IdTokenAlg             types.String `tfsdk:"id_token_alg"`
	AuthCodeLifetime       types.Int64  `tfsdk:"auth_code_lifetime"`
	AccessTokenLifetime    types.Int64  `tfsdk:"access_token_lifetime"`
	Scopes                 types.List   `tfsdk:"scopes"`
	DefaultScopes          types.List   `tfsdk:"default_scopes"`
	Challenges             types.List   `tfsdk:"challenges"`
	ForceMfa               types.Bool   `tfsdk:"force_mfa"`
	ClientUri              types.String `tfsdk:"client_uri"`
	Contacts               types.List   `tfsdk:"contacts"`
}

func (r *ClientResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

func (r *ClientResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Client resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Client ID",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Client name",
				Required:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Client enabled",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
			"confidential": schema.BoolAttribute{
				MarkdownDescription: "Client confidential",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"redirect_uris": schema.ListAttribute{
				MarkdownDescription: "Client redirect URIs",
				ElementType:         types.StringType,
				Computed:            true,
				Optional:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"post_logout_redirect_uris": schema.ListAttribute{
				MarkdownDescription: "Client post logout redirect URIs",
				ElementType:         types.StringType,
				Computed:            true,
				Optional:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"flows_enabled": schema.ListAttribute{
				MarkdownDescription: "Client flows enabled",
				ElementType:         types.StringType,
				Computed:            true,
				Optional:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("authorization_code")})),
			},
			"access_token_alg": schema.StringAttribute{
				MarkdownDescription: "Client access token algorithm",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("EdDSA"),
			},
			"id_token_alg": schema.StringAttribute{
				MarkdownDescription: "Client ID token algorithm",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("EdDSA"),
			},
			"auth_code_lifetime": schema.Int64Attribute{
				MarkdownDescription: "Client auth code lifetime",
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(60),
			},
			"access_token_lifetime": schema.Int64Attribute{
				MarkdownDescription: "Client access token lifetime",
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(1800),
			},
			"scopes": schema.ListAttribute{
				MarkdownDescription: "Client scopes",
				ElementType:         types.StringType,
				Computed:            true,
				Optional:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("openid")})),
			},
			"default_scopes": schema.ListAttribute{
				MarkdownDescription: "Client default scopes",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"challenges": schema.ListAttribute{
				MarkdownDescription: "Client challenges",
				ElementType:         types.StringType,
				Computed:            true,
				Optional:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("S256")})),
			},
			"force_mfa": schema.BoolAttribute{
				MarkdownDescription: "Client force MFA",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"client_uri": schema.StringAttribute{
				MarkdownDescription: "Client URI",
				Computed:            true,
			},
			"contacts": schema.ListAttribute{
				MarkdownDescription: "Client contacts",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (r *ClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*rauthy.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *rauthy.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClientResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiModel := data.ToApi()
	createPayload := apiModel.ToCreatePayload()
	newClient, err := r.client.CreateOidcClient(ctx, &createPayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	client, err := r.client.UpdateOidcClient(ctx, newClient.ID, &apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update client, got error: %s", err))
		return
	}

	data.FromApiResource(&client)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClientResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.GetOidcClient(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read client, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClientResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := data.ToApi()
	client, err := r.client.UpdateOidcClient(ctx, data.Id.ValueString(), &payload)
	data.FromApiResource(&client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update client, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClientResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteOidcClient(ctx, data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete client, got error: %s", err))
		return
	}
}

func (r *ClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ClientResourceModel) ToApi() rauthy.OidcClient {
	return rauthy.OidcClient{
		ID:                  r.Id.ValueString(),
		Name:                r.Name.ValueString(),
		Enabled:             r.Enabled.ValueBool(),
		Confidential:        r.Confidential.ValueBool(),
		RedirectUris:        tfutils.ListToStringSlice(r.RedirectUris),
		FlowsEnabled:        tfutils.ListToStringSlice(r.FlowsEnabled),
		AccessTokenAlg:      r.AccessTokenAlg.ValueString(),
		IdTokenAlg:          r.IdTokenAlg.ValueString(),
		AuthCodeLifetime:    r.AuthCodeLifetime.ValueInt64(),
		AccessTokenLifetime: r.AccessTokenLifetime.ValueInt64(),
		Scopes:              tfutils.ListToStringSlice(r.Scopes),
		DefaultScopes:       tfutils.ListToStringSlice(r.DefaultScopes),
		Challenges:          tfutils.ListToStringSlice(r.Challenges),
		ForceMfa:            r.ForceMfa.ValueBool(),
	}
}

func (r *ClientResourceModel) FromApiResource(client *rauthy.OidcClient) {
	r.Scopes = tfutils.StringSliceToList(client.Scopes)
	r.AccessTokenAlg = types.StringValue(client.AccessTokenAlg)
	r.AuthCodeLifetime = types.Int64Value(client.AuthCodeLifetime)
	r.DefaultScopes = tfutils.StringSliceToList(client.DefaultScopes)
	r.ForceMfa = types.BoolValue(client.ForceMfa)
	r.IdTokenAlg = types.StringValue(client.IdTokenAlg)
	r.AccessTokenLifetime = types.Int64Value(client.AccessTokenLifetime)
	r.ClientUri = types.StringValue(client.ClientUri)
	r.Contacts = tfutils.StringSliceToList(client.Contacts)
}
