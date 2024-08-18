package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sousa-andre/terraform-provider-benthos/internal/client"
)

var (
	_ resource.Resource              = &streamResource{}
	_ resource.ResourceWithConfigure = &streamResource{}
)

type streamResource struct {
	client *client.BenthosClient
}

func NewStreamResource() resource.Resource {
	return &streamResource{}
}

type streamResourceModel struct {
	Id     types.String `tfsdk:"id"`
	Config types.String `tfsdk:"config"`
}

func (r *streamResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_stream"
}

func (r *streamResource) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"config": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *streamResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.BenthosClient)
}

func (r *streamResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan streamResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if res.Diagnostics.HasError() {
		return
	}

	streamId := plan.Id.ValueString()
	streamConfig := plan.Config.ValueString()
	// TODO: validate configuration

	err := r.client.CreateStream(streamId, streamConfig)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to create new stream",
			err.Error(),
		)
		return
	}

	state := streamResourceModel{
		Id:     types.StringValue(streamId),
		Config: types.StringValue(streamConfig),
	}
	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

func (r *streamResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state streamResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)

	streamId := state.Id.ValueString()
	stream, err := r.client.GetStream(streamId)
	if err != nil {
		res.Diagnostics.AddError(
			fmt.Sprintf("Failed to query stream %s", streamId),
			err.Error(),
		)
		return
	}

	newConfiguration, err := json.Marshal(stream.Configuration)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to marshall configuration",
			err.Error(),
		)
		return
	}

	newState := streamResourceModel{
		Id:     types.StringValue(streamId),
		Config: types.StringValue(string(newConfiguration)),
	}
	res.Diagnostics.Append(res.State.Set(ctx, &newState)...)
}

func (r *streamResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan streamResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	streamId := plan.Id.ValueString()
	streamConfig := plan.Config.ValueString()

	err := r.client.UpdateStream(streamId, streamConfig)
	if err != nil {
		res.Diagnostics.AddError(
			fmt.Sprintf("Failed to update stream %s", streamId),
			err.Error(),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

func (r *streamResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state streamResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	streamId := state.Id.ValueString()
	if err := r.client.DeleteStream(streamId); err != nil {
		res.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete stream %s", streamId),
			err.Error(),
		)
		return
	}
}
