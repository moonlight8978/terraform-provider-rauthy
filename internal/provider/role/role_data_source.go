package role

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
)

var _ datasource.DataSource = &RoleDataSource{}

func NewRoleDataSource() datasource.DataSource {
	return &RoleDataSource{}
}

type RoleDataSource struct {
	client *rauthy.Client
}

type RoleDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *RoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (d *RoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Role data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Role ID",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Role name",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (d *RoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RoleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() && data.Id.IsNull() {
		resp.Diagnostics.AddError("Invalid Attribute Combination", "Either 'name' or 'id' must be specified.")
		return
	}

	roles, err := d.client.GetRoles(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read roles, got error: %s", err))
		return
	}

	var foundRole *rauthy.Role
	for _, r := range roles {
		if !data.Id.IsNull() && r.Id == data.Id.ValueString() {
			foundRole = &r
			break
		}
		if !data.Name.IsNull() && r.Name == data.Name.ValueString() {
			foundRole = &r
			break
		}
	}

	if foundRole == nil {
		resp.Diagnostics.AddError("Role Not Found", "No role found with the specified criteria.")
		return
	}

	data.Id = types.StringValue(foundRole.Id)
	data.Name = types.StringValue(foundRole.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
