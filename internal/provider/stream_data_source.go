package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sousa-andre/terraform-provider-benthos/internal/client"
)

var (
	_ datasource.DataSource              = &streamDatasource{}
	_ datasource.DataSourceWithConfigure = &streamDatasource{}
)

type streamDatasource struct {
	client *client.BenthosClient
}

func NewStreamDataSource() datasource.DataSource {
	return &streamDatasource{}
}

type streamDatasourceModel struct {
	Id             types.String  `tfsdk:"id"`
	Active         types.Bool    `tfsdk:"active"`
	Uptime         types.Float64 `tfsdk:"uptime"`
	ReadableUptime types.String  `tfsdk:"readable_uptime"`
	Config         types.String  `tfsdk:"config"`
}

func (d *streamDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_stream"
}

func (d *streamDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"active": schema.BoolAttribute{
				Computed: true,
			},
			"uptime": schema.Float64Attribute{
				Computed: true,
			},
			"readable_uptime": schema.StringAttribute{
				Computed: true,
			},
			"config": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *streamDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.BenthosClient)
}

func (d *streamDatasource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	var config streamDatasourceModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if res.Diagnostics.HasError() {
		return
	}
	streamId := config.Id.ValueString()

	detailedStream, err := d.client.GetStream(streamId)

	if err != nil {
		res.Diagnostics.AddError(
			fmt.Sprintf("Could not retrieve the stream with id %s", streamId),
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("requested stream for id %s", streamId))

	streamConfig, err := json.Marshal(detailedStream.Configuration)
	if err != nil {
		res.Diagnostics.AddError(
			fmt.Sprintf("Failed to marshall stream configuration"),
			err.Error(),
		)
	}
	streamModel := streamDatasourceModel{
		Id:             types.StringValue(streamId),
		Active:         types.BoolValue(detailedStream.Active),
		Uptime:         types.Float64Value(detailedStream.Uptime),
		ReadableUptime: types.StringValue(detailedStream.UptimeStr),
		Config:         types.StringValue(string(streamConfig)),
	}

	res.Diagnostics.Append(res.State.Set(ctx, &streamModel)...)
	if res.Diagnostics.HasError() {
		return
	}
}
