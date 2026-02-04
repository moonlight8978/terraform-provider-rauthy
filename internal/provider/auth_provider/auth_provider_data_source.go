package auth_provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
)

var _ datasource.DataSource = &AuthProviderDataSource{}

func NewAuthProviderDataSource() datasource.DataSource {
	return &AuthProviderDataSource{}
}

type AuthProviderDataSource struct {
	client *rauthy.Client
}

type AuthProviderDataSourceModel struct {
	AdminClaimPath        types.String `tfsdk:"admin_claim_path"`
	AdminClaimValue       types.String `tfsdk:"admin_claim_value"`
	AuthorizationEndpoint types.String `tfsdk:"authorization_endpoint"`
	AutoLink              types.Bool   `tfsdk:"auto_link"`
	AutoOnboarding        types.Bool   `tfsdk:"auto_onboarding"`
	ClientId              types.String `tfsdk:"client_id"`
	ClientSecret          types.String `tfsdk:"client_secret"`
	ClientSecretBasic     types.Bool   `tfsdk:"client_secret_basic"`
	ClientSecretPost      types.Bool   `tfsdk:"client_secret_post"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	Id                    types.String `tfsdk:"id"`
	Issuer                types.String `tfsdk:"issuer"`
	JwksEndpoint          types.String `tfsdk:"jwks_endpoint"`
	MfaClaimPath          types.String `tfsdk:"mfa_claim_path"`
	MfaClaimValue         types.String `tfsdk:"mfa_claim_value"`
	Name                  types.String `tfsdk:"name"`
	Scope                 types.String `tfsdk:"scope"`
	TokenEndpoint         types.String `tfsdk:"token_endpoint"`
	Typ                   types.String `tfsdk:"typ"`
	UsePkce               types.Bool   `tfsdk:"use_pkce"`
	UserinfoEndpoint      types.String `tfsdk:"userinfo_endpoint"`
}

func (d *AuthProviderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_provider"
}

func (d *AuthProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Auth Provider data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Provider ID",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Provider Name",
				Computed:            true,
			},
			"issuer": schema.StringAttribute{
				MarkdownDescription: "Provider Issuer",
				Computed:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID",
				Computed:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client Secret",
				Computed:            true,
				Sensitive:           true,
			},
			"authorization_endpoint": schema.StringAttribute{
				MarkdownDescription: "Authorization Endpoint",
				Computed:            true,
			},
			"token_endpoint": schema.StringAttribute{
				MarkdownDescription: "Token Endpoint",
				Computed:            true,
			},
			"userinfo_endpoint": schema.StringAttribute{
				MarkdownDescription: "Userinfo Endpoint",
				Computed:            true,
			},
			"jwks_endpoint": schema.StringAttribute{
				MarkdownDescription: "JWKS Endpoint",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "Scope",
				Computed:            true,
			},
			"admin_claim_path": schema.StringAttribute{
				MarkdownDescription: "Admin Claim Path",
				Computed:            true,
			},
			"admin_claim_value": schema.StringAttribute{
				MarkdownDescription: "Admin Claim Value",
				Computed:            true,
			},
			"mfa_claim_path": schema.StringAttribute{
				MarkdownDescription: "MFA Claim Path",
				Computed:            true,
			},
			"mfa_claim_value": schema.StringAttribute{
				MarkdownDescription: "MFA Claim Value",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enabled",
				Computed:            true,
			},
			"auto_link": schema.BoolAttribute{
				MarkdownDescription: "Auto Link",
				Computed:            true,
			},
			"auto_onboarding": schema.BoolAttribute{
				MarkdownDescription: "Auto Onboarding",
				Computed:            true,
			},
			"client_secret_basic": schema.BoolAttribute{
				MarkdownDescription: "Client Secret Basic",
				Computed:            true,
			},
			"client_secret_post": schema.BoolAttribute{
				MarkdownDescription: "Client Secret Post",
				Computed:            true,
			},
			"use_pkce": schema.BoolAttribute{
				MarkdownDescription: "Use PKCE",
				Computed:            true,
			},
			"typ": schema.StringAttribute{
				MarkdownDescription: "Type",
				Computed:            true,
			},
		},
	}
}

func (d *AuthProviderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AuthProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AuthProviderDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	provider, err := d.client.GetAuthProvider(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read OIDC provider, got error: %s", err))
		return
	}

	// Map fields
	data.Name = types.StringValue(provider.Name)
	data.Issuer = types.StringValue(provider.Issuer)
	data.ClientId = types.StringValue(provider.ClientId)
	data.ClientSecret = types.StringValue(provider.ClientSecret)
	data.AuthorizationEndpoint = types.StringValue(provider.AuthorizationEndpoint)
	data.TokenEndpoint = types.StringValue(provider.TokenEndpoint)
	data.UserinfoEndpoint = types.StringValue(provider.UserinfoEndpoint)
	data.JwksEndpoint = types.StringValue(provider.JwksEndpoint)
	data.Scope = types.StringValue(provider.Scope)
	data.AdminClaimPath = types.StringValue(provider.AdminClaimPath)
	data.AdminClaimValue = types.StringValue(provider.AdminClaimValue)
	data.MfaClaimPath = types.StringValue(provider.MfaClaimPath)
	data.MfaClaimValue = types.StringValue(provider.MfaClaimValue)
	data.Enabled = types.BoolValue(provider.Enabled)
	data.AutoLink = types.BoolValue(provider.AutoLink)
	data.AutoOnboarding = types.BoolValue(provider.AutoOnboarding)
	data.ClientSecretBasic = types.BoolValue(provider.ClientSecretBasic)
	data.ClientSecretPost = types.BoolValue(provider.ClientSecretPost)
	data.UsePkce = types.BoolValue(provider.UsePkce)
	data.Typ = types.StringValue(provider.Typ)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
