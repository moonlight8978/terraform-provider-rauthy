// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package passwordpolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PasswordPolicyResource{}
var _ resource.ResourceWithImportState = &PasswordPolicyResource{}

func NewPasswordPolicyResource() resource.Resource {
	return &PasswordPolicyResource{}
}

// ExampleResource defines the resource implementation.
type PasswordPolicyResource struct {
	client *rauthy.Client
}

// ExampleResourceModel describes the resource data model.
type PasswordPolicyResourceModel struct {
	IncludeUpperCase types.Int64 `tfsdk:"include_upper_case"`
	IncludeLowerCase types.Int64 `tfsdk:"include_lower_case"`
	IncludeDigits    types.Int64 `tfsdk:"include_digits"`
	IncludeSpecial   types.Int64 `tfsdk:"include_special"`
	LengthMin        types.Int64 `tfsdk:"length_min"`
	LengthMax        types.Int64 `tfsdk:"length_max"`
	NotRecentlyUsed  types.Int64 `tfsdk:"not_recently_used"`
	ValidDays        types.Int64 `tfsdk:"valid_days"`
}

func (r *PasswordPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_policy"
}

func (r *PasswordPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Password policy resource",

		Attributes: map[string]schema.Attribute{
			"include_digits": schema.Int64Attribute{
				MarkdownDescription: "Number of digits required in password",
				Optional:            true,
				Default:             int64default.StaticInt64(1),
				Computed:            true,
			},
			"include_lower_case": schema.Int64Attribute{
				MarkdownDescription: "Number of lowercase letters required in password",
				Optional:            true,
				Default:             int64default.StaticInt64(1),
				Computed:            true,
			},
			"include_upper_case": schema.Int64Attribute{
				MarkdownDescription: "Number of uppercase letters required in password",
				Optional:            true,
				Default:             int64default.StaticInt64(1),
				Computed:            true,
			},
			"include_special": schema.Int64Attribute{
				MarkdownDescription: "Number of special characters required in password",
				Optional:            true,
				Default:             int64default.StaticInt64(0),
				Computed:            true,
			},
			"length_max": schema.Int64Attribute{
				MarkdownDescription: "Maximum length of password",
				Required:            true,
			},
			"length_min": schema.Int64Attribute{
				MarkdownDescription: "Minimum length of password",
				Required:            true,
			},
			"not_recently_used": schema.Int64Attribute{
				MarkdownDescription: "Number of previous passwords that cannot be reused",
				Optional:            true,
				Default:             int64default.StaticInt64(3),
				Computed:            true,
			},
			"valid_days": schema.Int64Attribute{
				MarkdownDescription: "Number of days before password expires",
				Optional:            true,
				Default:             int64default.StaticInt64(180),
				Computed:            true,
			},
		},
	}
}

func (r *PasswordPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PasswordPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model PasswordPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdatePasswordPolicy(ctx, model.ToPayload())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create password policy, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *PasswordPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PasswordPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PasswordPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PasswordPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdatePasswordPolicy(ctx, data.ToPayload())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update password policy, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PasswordPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PasswordPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PasswordPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	panic("not supported")
}

func (model *PasswordPolicyResourceModel) ToPayload() *rauthy.PasswordPolicy {
	return &rauthy.PasswordPolicy{
		IncludeUpperCase: int(model.IncludeUpperCase.ValueInt64()),
		IncludeLowerCase: int(model.IncludeLowerCase.ValueInt64()),
		IncludeDigits:    int(model.IncludeDigits.ValueInt64()),
		IncludeSpecial:   int(model.IncludeSpecial.ValueInt64()),
		LengthMin:        int(model.LengthMin.ValueInt64()),
		LengthMax:        int(model.LengthMax.ValueInt64()),
		NotRecentlyUsed:  int(model.NotRecentlyUsed.ValueInt64()),
		ValidDays:        int(model.ValidDays.ValueInt64()),
	}
}
