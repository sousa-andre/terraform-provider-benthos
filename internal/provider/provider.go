package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sousa-andre/terraform-provider-benthos/internal/client"
)

var _ provider.Provider = &benthosProvider{}

type benthosProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &benthosProvider{version: version}
	}
}

type benthosProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Protocol types.String `tfsdk:"protocol"`
}

func (p *benthosProvider) Metadata(ctx context.Context, req provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "benthos"
	res.Version = p.version
}

func (p *benthosProvider) Schema(ctx context.Context, req provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required: true,
			},
			"protocol": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *benthosProvider) Configure(ctx context.Context, req provider.ConfigureRequest, res *provider.ConfigureResponse) {
	var config benthosProviderModel
	var endpoint, protocol string

	diag := req.Config.Get(ctx, &config)
	res.Diagnostics.Append(diag...)

	if res.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() {
		res.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Benthos Endpoint",
			"Benthos provider Endpoint not specified. Please use the provider Endpoint attribute to set the API",
		)
		return
	}

	endpoint = config.Endpoint.ValueString()

	protocol = config.Protocol.ValueString()
	if protocol == "" {
		protocol = "http"
	}

	client, err := client.NewClient(fmt.Sprintf("%s://%s", protocol, endpoint))
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to create the client",
			"The Benthos client failed to initialize likely due to wrong credentials "+
				"or failed connection",
		)
		return
	}
	res.ResourceData = client
	res.DataSourceData = client
}

func (p *benthosProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewStreamDataSource,
	}
}

func (p *benthosProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewStreamResource,
	}
}
