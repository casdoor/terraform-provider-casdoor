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
	"fmt"
	"strings"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type objectKind struct {
	name        string
	api         string
	description string
	adminOwned  bool
	setPassword bool
	ownerField  field
	fields      []field
}

type objectResource struct {
	kind   objectKind
	client *casdoorsdk.Client
}

var (
	_ resource.ResourceWithConfigure   = &objectResource{}
	_ resource.ResourceWithImportState = &objectResource{}
)

func newObjectResource(kind objectKind) func() resource.Resource {
	return func() resource.Resource {
		return &objectResource{kind: kind}
	}
}

func (k objectKind) allFields() []field {
	var res []field
	if !k.adminOwned {
		res = append(res, k.ownerField)
	}
	res = append(res,
		str("name", "name", fmt.Sprintf("The name of the %s, unique within its owner.", k.name)).req().replace(),
		str("created_time", "createdTime", fmt.Sprintf("The time when the %s was created.", k.name)).readOnly(),
	)
	return append(res, k.fields...)
}

func (k objectKind) attrTypes() map[string]attr.Type {
	res := map[string]attr.Type{"id": types.StringType}
	for _, f := range k.allFields() {
		res[f.name] = f.attrType()
	}
	return res
}

func (k objectKind) id(attrs map[string]attr.Value) string {
	owner := "admin"
	if !k.adminOwned {
		owner = attrs["owner"].(types.String).ValueString()
	}
	return owner + "/" + attrs["name"].(types.String).ValueString()
}

func (r *objectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.kind.name
}

func (r *objectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"id": str("id", "", fmt.Sprintf("The ID of the %s in the `owner/name` format.", r.kind.name)).readOnly().schemaAttribute(),
	}
	for _, f := range r.kind.allFields() {
		attributes[f.name] = f.schemaAttribute()
	}
	resp.Schema = schema.Schema{MarkdownDescription: r.kind.description, Attributes: attributes}
}

func (r *objectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*casdoorsdk.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("Expected *casdoorsdk.Client, got %T.", req.ProviderData))
		return
	}
	r.client = client
}

func (r *objectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := plan.Attributes()
	id := r.kind.id(attrs)
	owner, name, _ := strings.Cut(id, "/")
	obj := map[string]any{
		"owner":       owner,
		"name":        name,
		"createdTime": casdoorsdk.GetCurrentTime(),
	}
	r.overlay(obj, attrs)

	err := modifyObject(r.client, "add-"+r.kind.api, id, obj)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to create %s %s", r.kind.name, id), err.Error())
		return
	}

	r.readInto(ctx, id, attrs, &resp.State, &resp.Diagnostics)
}

func (r *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := state.Attributes()
	id := attrs["id"].(types.String).ValueString()
	if id == "" {
		id = r.kind.id(attrs)
	}

	obj, err := getObject(r.client, r.kind.api, id)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read %s %s", r.kind.name, id), err.Error())
		return
	}
	if obj == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	r.setState(ctx, id, obj, attrs, &resp.State, &resp.Diagnostics)
}

func (r *objectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := plan.Attributes()
	id := r.kind.id(attrs)
	obj, err := getObject(r.client, r.kind.api, id)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read %s %s", r.kind.name, id), err.Error())
		return
	}
	if obj == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to update %s %s", r.kind.name, id), "The object no longer exists.")
		return
	}

	r.overlay(obj, attrs)
	err = modifyObject(r.client, "update-"+r.kind.api, id, obj)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to update %s %s", r.kind.name, id), err.Error())
		return
	}

	if r.kind.setPassword {
		var state types.Object
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}

		password := attrs["password"]
		if !password.IsNull() && !password.IsUnknown() && !password.Equal(state.Attributes()["password"]) {
			owner, name, _ := strings.Cut(id, "/")
			_, err = r.client.SetPassword(owner, name, "", password.(types.String).ValueString())
			if err != nil && !strings.Contains(err.Error(), "different from your current password") {
				resp.Diagnostics.AddError(fmt.Sprintf("Failed to set the password of %s %s", r.kind.name, id), err.Error())
				return
			}
		}
	}

	r.readInto(ctx, id, attrs, &resp.State, &resp.Diagnostics)
}

func (r *objectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := r.kind.id(state.Attributes())
	obj, err := getObject(r.client, r.kind.api, id)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read %s %s", r.kind.name, id), err.Error())
		return
	}
	if obj == nil {
		return
	}

	err = modifyObject(r.client, "delete-"+r.kind.api, id, obj)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to delete %s %s", r.kind.name, id), err.Error())
	}
}

func (r *objectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	owner, name, found := strings.Cut(req.ID, "/")
	if !found {
		owner, name = "admin", req.ID
	}
	if name == "" || (r.kind.adminOwned && owner != "admin") {
		format := "owner/name"
		if r.kind.adminOwned {
			format = "name"
		}
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected %q, got %q.", format, req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), owner+"/"+name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	if !r.kind.adminOwned {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), owner)...)
	}
}

func (r *objectResource) overlay(obj map[string]any, attrs map[string]attr.Value) {
	for _, f := range r.kind.fields {
		value := attrs[f.name]
		if f.computed || value == nil || value.IsNull() || value.IsUnknown() {
			continue
		}
		obj[f.key] = f.fromValue(value, obj[f.key])
	}
}

func (r *objectResource) readInto(ctx context.Context, id string, prior map[string]attr.Value, state stateSetter, diags *diag.Diagnostics) {
	obj, err := getObject(r.client, r.kind.api, id)
	if err != nil {
		diags.AddError(fmt.Sprintf("Failed to read %s %s", r.kind.name, id), err.Error())
		return
	}
	if obj == nil {
		diags.AddError(fmt.Sprintf("Failed to read %s %s", r.kind.name, id), "The object was not found after it was saved.")
		return
	}

	r.setState(ctx, id, obj, prior, state, diags)
}

type stateSetter interface {
	Set(ctx context.Context, val any) diag.Diagnostics
}

func (r *objectResource) setState(ctx context.Context, id string, obj map[string]any, prior map[string]attr.Value, state stateSetter, diags *diag.Diagnostics) {
	attrs := map[string]attr.Value{"id": types.StringValue(id)}
	for _, f := range r.kind.allFields() {
		raw := obj[f.key]
		if f.writeOnly || (f.sensitive && raw == "***") {
			value := prior[f.name]
			if value == nil || value.IsUnknown() {
				value = f.nullValue()
			}
			attrs[f.name] = value
			continue
		}

		value, d := f.toValue(raw)
		diags.Append(d...)
		attrs[f.name] = value
	}
	if diags.HasError() {
		return
	}

	value, d := types.ObjectValue(r.kind.attrTypes(), attrs)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	diags.Append(state.Set(ctx, value)...)
}
