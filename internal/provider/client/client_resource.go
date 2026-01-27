package client

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
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
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	Confidential        types.Bool   `tfsdk:"confidential"`
	RedirectUris        types.List   `tfsdk:"redirect_uris"`
	FlowsEnabled        types.List   `tfsdk:"flows_enabled"`
	AccessTokenAlg      types.String `tfsdk:"access_token_alg"`
	IdTokenAlg          types.String `tfsdk:"id_token_alg"`
	AuthCodeLifetime    types.Int64  `tfsdk:"auth_code_lifetime"`
	AccessTokenLifetime types.Int64  `tfsdk:"access_token_lifetime"`
	Scopes              types.List   `tfsdk:"scopes"`
	DefaultScopes       types.List   `tfsdk:"default_scopes"`
	Challenges          types.List   `tfsdk:"challenges"`
	ForceMfa            types.Bool   `tfsdk:"force_mfa"`
	ClientUri           types.String `tfsdk:"client_uri"`
	Contacts            types.List   `tfsdk:"contacts"`
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
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Client name",
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
			},
			"flows_enabled": schema.ListAttribute{
				MarkdownDescription: "Client flows enabled",
				ElementType:         types.StringType,
			},
			"access_token_alg": schema.StringAttribute{
				MarkdownDescription: "Client access token algorithm",
			},
			"id_token_alg": schema.StringAttribute{
				MarkdownDescription: "Client ID token algorithm",
			},
			"auth_code_lifetime": schema.Int64Attribute{
				MarkdownDescription: "Client auth code lifetime",
			},
			"access_token_lifetime": schema.Int64Attribute{
				MarkdownDescription: "Client access token lifetime",
			},
			"scopes": schema.ListAttribute{
				MarkdownDescription: "Client scopes",
				ElementType:         types.StringType,
			},
			"default_scopes": schema.ListAttribute{
				MarkdownDescription: "Client default scopes",
				ElementType:         types.StringType,
			},
			"challenges": schema.ListAttribute{
				MarkdownDescription: "Client challenges",
				ElementType:         types.StringType,
			},
			"force_mfa": schema.BoolAttribute{
				MarkdownDescription: "Client force MFA",
			},
			"client_uri": schema.StringAttribute{
				MarkdownDescription: "Client URI",
			},
			"contacts": schema.ListAttribute{
				MarkdownDescription: "Client contacts",
				ElementType:         types.StringType,
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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue("example-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClientResourceModel

	// Read Terraform prior state data into the model
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

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := data.ToApi()
	_, err := r.client.UpdateOidcClient(ctx, data.Id.ValueString(), &payload)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update client, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClientResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
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
		RedirectUris:        listToStringSlice(r.RedirectUris),
		FlowsEnabled:        listToStringSlice(r.FlowsEnabled),
		AccessTokenAlg:      r.AccessTokenAlg.ValueString(),
		IdTokenAlg:          r.IdTokenAlg.ValueString(),
		AuthCodeLifetime:    r.AuthCodeLifetime.ValueInt64(),
		AccessTokenLifetime: r.AccessTokenLifetime.ValueInt64(),
		Scopes:              listToStringSlice(r.Scopes),
		DefaultScopes:       listToStringSlice(r.DefaultScopes),
		Challenges:          listToStringSlice(r.Challenges),
		ForceMfa:            r.ForceMfa.ValueBool(),
		ClientUri:           r.ClientUri.ValueString(),
		Contacts:            listToStringSlice(r.Contacts),
	}
}

func listToStringSlice(l types.List) []string {
	if l.IsNull() || l.IsUnknown() {
		return []string{}
	}

	var result []string
	for _, val := range l.Elements() {
		result = append(result, val.(types.String).ValueString())
	}

	return result
}
