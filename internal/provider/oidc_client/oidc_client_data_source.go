package oidc_client

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
)

var _ datasource.DataSource = &OidcClientDataSource{}

func NewOidcClientDataSource() datasource.DataSource {
	return &OidcClientDataSource{}
}

type OidcClientDataSource struct {
	client *rauthy.Client
}

type OidcClientDataSourceModel struct {
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

func (d *OidcClientDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

func (d *OidcClientDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "OIDC Client data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the client",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the client",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the client is enabled",
				Computed:            true,
			},
			"confidential": schema.BoolAttribute{
				MarkdownDescription: "Whether the client is confidential",
				Computed:            true,
			},
			"redirect_uris": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Redirect URIs",
				Computed:            true,
			},
			"post_logout_redirect_uris": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Post logout redirect URIs",
				Computed:            true,
			},
			"flows_enabled": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Enabled flows",
				Computed:            true,
			},
			"access_token_alg": schema.StringAttribute{
				MarkdownDescription: "Access token algorithm",
				Computed:            true,
			},
			"id_token_alg": schema.StringAttribute{
				MarkdownDescription: "ID token algorithm",
				Computed:            true,
			},
			"auth_code_lifetime": schema.Int64Attribute{
				MarkdownDescription: "Auth code lifetime",
				Computed:            true,
			},
			"access_token_lifetime": schema.Int64Attribute{
				MarkdownDescription: "Access token lifetime",
				Computed:            true,
			},
			"scopes": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Scopes",
				Computed:            true,
			},
			"default_scopes": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Default scopes",
				Computed:            true,
			},
			"challenges": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Challenges",
				Computed:            true,
			},
			"force_mfa": schema.BoolAttribute{
				MarkdownDescription: "Force MFA",
				Computed:            true,
			},
			"client_uri": schema.StringAttribute{
				MarkdownDescription: "Client URI",
				Computed:            true,
			},
			"contacts": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Contacts",
				Computed:            true,
			},
		},
	}
}

func (d *OidcClientDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*rauthy.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *rauthy.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *OidcClientDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OidcClientDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	oidcClient, err := d.client.GetOidcClient(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read OIDC client, got error: %s", err))
		return
	}

	data.Name = types.StringValue(oidcClient.Name)
	data.Enabled = types.BoolValue(oidcClient.Enabled)
	data.Confidential = types.BoolValue(oidcClient.Confidential)
	data.AccessTokenAlg = types.StringValue(oidcClient.AccessTokenAlg)
	data.IdTokenAlg = types.StringValue(oidcClient.IdTokenAlg)
	data.AuthCodeLifetime = types.Int64Value(oidcClient.AuthCodeLifetime)
	data.AccessTokenLifetime = types.Int64Value(oidcClient.AccessTokenLifetime)
	data.ForceMfa = types.BoolValue(oidcClient.ForceMfa)

	if oidcClient.ClientUri != "" {
		data.ClientUri = types.StringValue(oidcClient.ClientUri)
	} else {
		data.ClientUri = types.StringNull()
	}

	var diags diag.Diagnostics

	data.RedirectUris, diags = types.ListValueFrom(ctx, types.StringType, oidcClient.RedirectUris)
	resp.Diagnostics.Append(diags...)

	data.PostLogoutRedirectUris, diags = types.ListValueFrom(ctx, types.StringType, oidcClient.PostLogoutUri)
	resp.Diagnostics.Append(diags...)

	data.FlowsEnabled, diags = types.ListValueFrom(ctx, types.StringType, oidcClient.FlowsEnabled)
	resp.Diagnostics.Append(diags...)

	data.Scopes, diags = types.ListValueFrom(ctx, types.StringType, oidcClient.Scopes)
	resp.Diagnostics.Append(diags...)

	data.DefaultScopes, diags = types.ListValueFrom(ctx, types.StringType, oidcClient.DefaultScopes)
	resp.Diagnostics.Append(diags...)

	data.Challenges, diags = types.ListValueFrom(ctx, types.StringType, oidcClient.Challenges)
	resp.Diagnostics.Append(diags...)

	if len(oidcClient.Contacts) > 0 {
		data.Contacts, diags = types.ListValueFrom(ctx, types.StringType, oidcClient.Contacts)
		resp.Diagnostics.Append(diags...)
	} else {
		data.Contacts = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
