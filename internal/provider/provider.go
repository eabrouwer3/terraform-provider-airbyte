package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure AirbyteProvider satisfies various provider interfaces.
var _ provider.Provider = &AirbyteProvider{}
var _ provider.ProviderWithMetadata = &AirbyteProvider{}

// AirbyteProvider defines the provider implementation.
type AirbyteProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AirbyteProviderModel describes the provider data model.
type AirbyteProviderModel struct {
	HostUrl types.String `tfsdk:"host_url"`
}

func (p *AirbyteProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "airbyte"
	resp.Version = p.version
}

func (p *AirbyteProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"host_url": {
				MarkdownDescription: "Airbyte API URL",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (p *AirbyteProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AirbyteProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.HostUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host_url"),
			"Unknown Airbyte API URL",
			"The provider cannot create the Airbyte API client as there is an unknown configuration value for the Airbyte API URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AIRBYTE_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	hostUrl, ok := os.LookupEnv("AIRBYTE_URL")
	if !ok {
		hostUrl = "http://localhost:8000"
	}

	if !data.HostUrl.IsNull() {
		hostUrl = data.HostUrl.Value
	}

	if hostUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host_url"),
			"Missing Airbyte API URL",
			"The provider cannot create the Airbyte API client as there is a missing or empty value for the Airbyte API URL. "+
				"Set the host value in the configuration or use the AIRBYTE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := apiclient.ApiClient{
		HostURL:    hostUrl,
		HTTPClient: &http.Client{Timeout: 120 * time.Second},
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *AirbyteProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWorkspaceResource,
		NewSourceDefinitionResource,
		NewSourceResource,
		NewDestinationDefinitionResource,
		NewDestinationResource,
	}
}

func (p *AirbyteProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkspaceDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AirbyteProvider{
			version: version,
		}
	}
}
