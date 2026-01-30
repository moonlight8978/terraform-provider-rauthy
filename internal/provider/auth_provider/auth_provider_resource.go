package auth_provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
)

var _ resource.Resource = &AuthProviderResource{}
var _ resource.ResourceWithImportState = &AuthProviderResource{}

func NewAuthProviderResource() resource.Resource {
	return &AuthProviderResource{}
}

type AuthProviderResource struct {
	client *rauthy.Client
}

type AuthProviderResourceModel struct {
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

func (r *AuthProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_provider"
}

func (r *AuthProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "OIDC Provider resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Provider ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Provider Name",
				Required:            true,
			},
			"issuer": schema.StringAttribute{
				MarkdownDescription: "Provider Issuer",
				Required:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID",
				Required:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client Secret",
				Required:            true,
				Sensitive:           true,
			},
			"authorization_endpoint": schema.StringAttribute{
				MarkdownDescription: "Authorization Endpoint",
				Optional:            true,
			},
			"token_endpoint": schema.StringAttribute{
				MarkdownDescription: "Token Endpoint",
				Optional:            true,
			},
			"userinfo_endpoint": schema.StringAttribute{
				MarkdownDescription: "Userinfo Endpoint",
				Optional:            true,
			},
			"jwks_endpoint": schema.StringAttribute{
				MarkdownDescription: "JWKS Endpoint",
				Optional:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "Scope",
				Optional:            true,
			},
			"admin_claim_path": schema.StringAttribute{
				MarkdownDescription: "Admin Claim Path",
				Optional:            true,
			},
			"admin_claim_value": schema.StringAttribute{
				MarkdownDescription: "Admin Claim Value",
				Optional:            true,
			},
			"mfa_claim_path": schema.StringAttribute{
				MarkdownDescription: "MFA Claim Path",
				Optional:            true,
			},
			"mfa_claim_value": schema.StringAttribute{
				MarkdownDescription: "MFA Claim Value",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enabled",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"auto_link": schema.BoolAttribute{
				MarkdownDescription: "Auto Link",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"auto_onboarding": schema.BoolAttribute{
				MarkdownDescription: "Auto Onboarding",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"client_secret_basic": schema.BoolAttribute{
				MarkdownDescription: "Client Secret Basic",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"client_secret_post": schema.BoolAttribute{
				MarkdownDescription: "Client Secret Post",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"use_pkce": schema.BoolAttribute{
				MarkdownDescription: "Use PKCE",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"typ": schema.StringAttribute{
				MarkdownDescription: "Type",
				Required:            true,
			},
		},
	}
}

func (r *AuthProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AuthProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AuthProviderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiModel := data.ToApi()
	newProvider, err := r.client.CreateAuthProvider(ctx, &apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create OIDC provider, got error: %s", err))
		return
	}

	data.FromApiResource(newProvider)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuthProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AuthProviderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	provider, err := r.client.GetAuthProvider(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read OIDC provider, got error: %s", err))
		return
	}

	data.FromApiResource(provider)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuthProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AuthProviderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiModel := data.ToApi()
	provider, err := r.client.UpdateAuthProvider(ctx, data.Id.ValueString(), &apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update OIDC provider, got error: %s", err))
		return
	}

	data.FromApiResource(provider)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuthProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AuthProviderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAuthProvider(ctx, data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete OIDC provider, got error: %s", err))
		return
	}
}

func (r *AuthProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AuthProviderResourceModel) ToApi() rauthy.AuthProvider {
	return rauthy.AuthProvider{
		Id:                    r.Id.ValueString(),
		Name:                  r.Name.ValueString(),
		Issuer:                r.Issuer.ValueString(),
		ClientId:              r.ClientId.ValueString(),
		ClientSecret:          r.ClientSecret.ValueString(),
		AuthorizationEndpoint: r.AuthorizationEndpoint.ValueString(),
		TokenEndpoint:         r.TokenEndpoint.ValueString(),
		UserinfoEndpoint:      r.UserinfoEndpoint.ValueString(),
		JwksEndpoint:          r.JwksEndpoint.ValueString(),
		Scope:                 r.Scope.ValueString(),
		AdminClaimPath:        r.AdminClaimPath.ValueString(),
		AdminClaimValue:       r.AdminClaimValue.ValueString(),
		MfaClaimPath:          r.MfaClaimPath.ValueString(),
		MfaClaimValue:         r.MfaClaimValue.ValueString(),
		Enabled:               r.Enabled.ValueBool(),
		AutoLink:              r.AutoLink.ValueBool(),
		AutoOnboarding:        r.AutoOnboarding.ValueBool(),
		ClientSecretBasic:     r.ClientSecretBasic.ValueBool(),
		ClientSecretPost:      r.ClientSecretPost.ValueBool(),
		UsePkce:               r.UsePkce.ValueBool(),
		Typ:                   r.Typ.ValueString(),
	}
}

func (r *AuthProviderResourceModel) FromApiResource(provider *rauthy.AuthProvider) {
	r.Id = types.StringValue(provider.Id)
	r.Name = types.StringValue(provider.Name)
	r.Issuer = types.StringValue(provider.Issuer)
	r.ClientId = types.StringValue(provider.ClientId)
	r.ClientSecret = types.StringValue(provider.ClientSecret)
	r.AuthorizationEndpoint = types.StringValue(provider.AuthorizationEndpoint)
	r.TokenEndpoint = types.StringValue(provider.TokenEndpoint)
	r.UserinfoEndpoint = types.StringValue(provider.UserinfoEndpoint)
	r.JwksEndpoint = types.StringValue(provider.JwksEndpoint)
	r.Scope = types.StringValue(provider.Scope)
	r.AdminClaimPath = types.StringValue(provider.AdminClaimPath)
	r.AdminClaimValue = types.StringValue(provider.AdminClaimValue)
	r.MfaClaimPath = types.StringValue(provider.MfaClaimPath)
	r.MfaClaimValue = types.StringValue(provider.MfaClaimValue)
	r.Enabled = types.BoolValue(provider.Enabled)
	r.AutoLink = types.BoolValue(provider.AutoLink)
	r.AutoOnboarding = types.BoolValue(provider.AutoOnboarding)
	r.ClientSecretBasic = types.BoolValue(provider.ClientSecretBasic)
	r.ClientSecretPost = types.BoolValue(provider.ClientSecretPost)
	r.UsePkce = types.BoolValue(provider.UsePkce)
	r.Typ = types.StringValue(provider.Typ)
}
