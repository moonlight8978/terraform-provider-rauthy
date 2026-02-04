package group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
)

var _ datasource.DataSource = &GroupDataSource{}

func NewGroupDataSource() datasource.DataSource {
	return &GroupDataSource{}
}

type GroupDataSource struct {
	client *rauthy.Client
}

type GroupDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *GroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *GroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Group data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Group ID",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Group name",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (d *GroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() && data.Id.IsNull() {
		resp.Diagnostics.AddError("Invalid Attribute Combination", "Either 'name' or 'id' must be specified.")
		return
	}

	groups, err := d.client.GetGroups(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read groups, got error: %s", err))
		return
	}

	var foundGroup *rauthy.Group
	for _, g := range groups {
		if !data.Id.IsNull() && g.Id == data.Id.ValueString() {
			foundGroup = &g
			break
		}
		if !data.Name.IsNull() && g.Name == data.Name.ValueString() {
			foundGroup = &g
			break
		}
	}

	if foundGroup == nil {
		resp.Diagnostics.AddError("Group Not Found", "No group found with the specified criteria.")
		return
	}

	data.Id = types.StringValue(foundGroup.Id)
	data.Name = types.StringValue(foundGroup.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
