// Code generated by tools/tfsdk2fw/main.go. Manual editing is required.

package meta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/fwtypes"
)

func init() {
	registerDataSourceFactory(newDataSourceARN)
}

// newDataSourceARN instantiates a new DataSource for the aws_arn data source.
func newDataSourceARN(ctx context.Context) (datasource.DataSource, error) {
	return &dataSourceARN{}, nil
}

type dataSourceARN struct{}

// Metadata should return the full name of the data source, such as
// examplecloud_thing.
func (d *dataSourceARN) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) { // nosemgrep:ci.meta-in-func-name
	response.TypeName = "aws_arn"
}

// GetSchema returns the schema for this data source.
func (d *dataSourceARN) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	schema := tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"account": {
				Type:     types.StringType,
				Computed: true,
			},
			"arn": {
				Type:     fwtypes.ARNType,
				Required: true,
			},
			"id": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"partition": {
				Type:     types.StringType,
				Computed: true,
			},
			"region": {
				Type:     types.StringType,
				Computed: true,
			},
			"resource": {
				Type:     types.StringType,
				Computed: true,
			},
			"service": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}

	return schema, nil
}

// Configure enables provider-level data or clients to be set in the
// provider-defined DataSource type. It is separately executed for each
// ReadDataSource RPC.
func (d *dataSourceARN) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) { //nolint:unparam
}

// Read is called when the provider must read data source values in order to update state.
// Config values should be read from the ReadRequest and new state values set on the ReadResponse.
func (d *dataSourceARN) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Trace(ctx, "dataSourceARN.Read enter")

	var config dataSourceARNData

	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)

	if response.Diagnostics.HasError() {
		return
	}

	state := config
	arn := &state.ARN.Value
	id := arn.String()

	state.Account = &arn.AccountID
	state.ID = &id
	state.Partition = &arn.Partition
	state.Region = &arn.Region
	state.Resource = &arn.Resource
	state.Service = &arn.Service

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

// TODO: Generate this structure definition.
type dataSourceARNData struct {
	Account   *string     `tfsdk:"account"`
	ARN       fwtypes.ARN `tfsdk:"arn"`
	ID        *string     `tfsdk:"id"`
	Partition *string     `tfsdk:"partition"`
	Region    *string     `tfsdk:"region"`
	Resource  *string     `tfsdk:"resource"`
	Service   *string     `tfsdk:"service"`
}
