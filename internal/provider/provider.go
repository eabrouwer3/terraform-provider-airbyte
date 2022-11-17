package provider

import (
	"context"
	"github.com/hashicorp/go-retryablehttp"
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
	HostUrl           types.String `tfsdk:"host_url"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	AdditionalHeaders types.Map    `tfsdk:"additional_headers"`
	Timeout           types.Int64  `tfsdk:"timeout"`
}

func (p *AirbyteProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "airbyte"
	resp.Version = p.version
}

func (p *AirbyteProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"host_url": {
				Description: "Airbyte API URL",
				Optional:    true,
				Type:        types.StringType,
			},
			"username": {
				Description: "Airbyte API Username",
				Optional:    true,
				Type:        types.StringType,
			},
			"password": {
				Description: "Airbyte API Password",
				Optional:    true,
				Type:        types.StringType,
				Sensitive:   true,
			},
			"additional_headers": {
				Description: "Additional Headers to pass in requests to Airbyte's API",
				Optional:    true,
				Type:        types.MapType{ElemType: types.StringType},
			},
			"timeout": {
				Description: "HTTP Timeout in Seconds (Default: 600)",
				Optional:    true,
				Type:        types.Int64Type,
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

	hostUrl, ok := os.LookupEnv("AIRBYTE_URL")
	if !ok {
		hostUrl = "http://localhost:8000"
	}
	if !data.HostUrl.IsNull() {
		hostUrl = data.HostUrl.ValueString()
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

	username, ok := os.LookupEnv("AIRBYTE_USERNAME")
	if !ok {
		username = "airbyte"
	}
	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if username == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("host_url"),
			"Missing Airbyte API Username",
			"There is a missing or empty value for the Airbyte API Username. This assumes authentication has been disabled for this Airbyte Instance."+
				"If this is not true, set the username value in the configuration or use the AIRBYTE_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	password, ok := os.LookupEnv("AIRBYTE_PASSWORD")
	if !ok {
		password = "password"
	}
	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	if password == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("password"),
			"Blank Airbyte API Password",
			"There is a missing or empty value for the Airbyte API Password. This assumes authentication has been disabled for this Airbyte Instance."+
				"If this is not true, set the password value in the configuration or use the AIRBYTE_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	additionalHeaders := data.AdditionalHeaders.Elements()
	additionalHeadersVals := make(map[string]string)
	for k, v := range additionalHeaders {
		additionalHeadersVals[k] = v.(types.String).ValueString()
	}

	timeout := time.Duration(600)
	if !data.Timeout.IsUnknown() {
		timeout = time.Duration(data.Timeout.ValueInt64())
	}

	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient = &http.Client{Timeout: timeout * time.Second}
	client := apiclient.ApiClient{
		HostURL:           hostUrl,
		Username:          username,
		Password:          password,
		AdditionalHeaders: additionalHeadersVals,
		HTTPClient:        httpClient,
	}

	err := client.Check()
	if err != nil {
		resp.Diagnostics.AddError(
			"Checking API Status Failed",
			err.Error(),
		)
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
		NewConnectionResource,
		NewOperationResource,
	}
}

func (p *AirbyteProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkspaceDataSource,
		NewWorkspaceIdsDataSource,
		NewSourceSchemaCatalogDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AirbyteProvider{
			version: version,
		}
	}
}
