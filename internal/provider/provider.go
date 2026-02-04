// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/auth_provider"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/group"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/oidc_client"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/passwordpolicy"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/role"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
)

// Ensure RauthyProvider satisfies various provider interfaces.
var _ provider.Provider = &RauthyProvider{}
var _ provider.ProviderWithFunctions = &RauthyProvider{}
var _ provider.ProviderWithEphemeralResources = &RauthyProvider{}
var _ provider.ProviderWithActions = &RauthyProvider{}

// RauthyProvider defines the provider implementation.
type RauthyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RauthyProviderModel describes the provider data model.
type RauthyProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func (p *RauthyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "rauthy"
	resp.Version = p.version
}

func (p *RauthyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *RauthyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var model RauthyProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	config := ProviderConfig{}
	config.FromEnv()
	config.Override(model)

	if err := config.Validate(); err != nil {
		resp.Diagnostics.AddError("Invalid provider configuration", err.Error())
		return
	}

	client := rauthy.NewClient(config.Endpoint, config.Insecure, rauthy.NewApiKeyAuthenticator(config.APIKey))

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *RauthyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		oidc_client.NewOidcClientResource,
		oidc_client.NewOidcClientSecretResource,
		auth_provider.NewAuthProviderResource,
		passwordpolicy.NewPasswordPolicyResource,
		role.NewRoleResource,
		group.NewGroupResource,
	}
}

func (p *RauthyProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		// NewExampleEphemeralResource,
	}
}

func (p *RauthyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// NewExampleDataSource,
	}
}

func (p *RauthyProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// NewExampleFunction,
	}
}

func (p *RauthyProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{
		// NewExampleAction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RauthyProvider{
			version: version,
		}
	}
}

type ProviderConfig struct {
	Endpoint string
	APIKey   string
	Insecure bool
}

func (c *ProviderConfig) FromEnv() {
	c.Endpoint = os.Getenv("RAUTHY_ENDPOINT")
	c.APIKey = os.Getenv("RAUTHY_API_KEY")
	c.Insecure = os.Getenv("RAUTHY_INSECURE") == "true"
}

func (c *ProviderConfig) Override(model RauthyProviderModel) {
	if !model.Endpoint.IsNull() {
		c.Endpoint = model.Endpoint.ValueString()
	}

	if !model.APIKey.IsNull() {
		c.APIKey = model.APIKey.ValueString()
	}

	if !model.Insecure.IsNull() {
		c.Insecure = model.Insecure.ValueBool()
	}
}

func (c *ProviderConfig) Validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("`endpoint` or `RAUTHY_ENDPOINT` is required")
	}

	if c.APIKey == "" {
		return fmt.Errorf("api_key` or `RAUTHY_API_KEY` is required")
	}

	return nil
}
