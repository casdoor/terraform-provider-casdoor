// Copyright 2026 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"os"
	"strings"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type casdoorProvider struct {
	version string
}

type casdoorProviderModel struct {
	Endpoint     types.String `tfsdk:"endpoint"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &casdoorProvider{version: version}
	}
}

func (p *casdoorProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "casdoor"
	resp.Version = p.version
}

func (p *casdoorProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Casdoor provider manages the objects of a [Casdoor](https://casdoor.ai) server: organizations, applications, users, providers, certs, roles, permissions and groups. See the [Terraform guide](https://casdoor.ai/docs/deployment/terraform) for an overview.\n\n" +
			"The provider calls the Casdoor API as an application, using its client ID and client secret. To manage every organization, use an application of the `built-in` organization, e.g. `app-built-in`.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The URL of the Casdoor server, e.g. `https://door.casdoor.com`. Can also be set with the `CASDOOR_ENDPOINT` environment variable.",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client ID of the application. Can also be set with the `CASDOOR_CLIENT_ID` environment variable.",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The client secret of the application. Can also be set with the `CASDOOR_CLIENT_SECRET` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *casdoorProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config casdoorProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := configValue(config.Endpoint, "CASDOOR_ENDPOINT")
	clientId := configValue(config.ClientId, "CASDOOR_CLIENT_ID")
	clientSecret := configValue(config.ClientSecret, "CASDOOR_CLIENT_SECRET")

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(path.Root("endpoint"), "Missing Casdoor endpoint", "Set endpoint in the provider block or the CASDOOR_ENDPOINT environment variable.")
	}
	if clientId == "" {
		resp.Diagnostics.AddAttributeError(path.Root("client_id"), "Missing Casdoor client ID", "Set client_id in the provider block or the CASDOOR_CLIENT_ID environment variable.")
	}
	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(path.Root("client_secret"), "Missing Casdoor client secret", "Set client_secret in the provider block or the CASDOOR_CLIENT_SECRET environment variable.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client := casdoorsdk.NewClient(strings.TrimRight(endpoint, "/"), clientId, clientSecret, "", "", "")
	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *casdoorProvider) Resources(_ context.Context) []func() resource.Resource {
	var res []func() resource.Resource
	for _, kind := range objectKinds {
		res = append(res, newObjectResource(kind))
	}
	return res
}

func (p *casdoorProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func configValue(value types.String, env string) string {
	if !value.IsNull() && !value.IsUnknown() {
		return value.ValueString()
	}
	return os.Getenv(env)
}
