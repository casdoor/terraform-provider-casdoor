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
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type fieldType int

const (
	stringField fieldType = iota
	boolField
	intField
	floatField
	stringListField
	stringMapField
	objectListField
)

type field struct {
	name        string
	key         string
	typ         fieldType
	description string

	required  bool
	forceNew  bool
	computed  bool
	sensitive bool
	writeOnly bool
	fallback  string

	fields []field
}

func str(name, key, description string) field {
	return field{name: name, key: key, typ: stringField, description: description}
}

func boolean(name, key, description string) field {
	return field{name: name, key: key, typ: boolField, description: description}
}

func integer(name, key, description string) field {
	return field{name: name, key: key, typ: intField, description: description}
}

func float(name, key, description string) field {
	return field{name: name, key: key, typ: floatField, description: description}
}

func strList(name, key, description string) field {
	return field{name: name, key: key, typ: stringListField, description: description}
}

func strMap(name, key, description string) field {
	return field{name: name, key: key, typ: stringMapField, description: description}
}

func objList(name, key, description string, fields ...field) field {
	return field{name: name, key: key, typ: objectListField, description: description, fields: fields}
}

func (f field) req() field {
	f.required = true
	return f
}

func (f field) replace() field {
	f.forceNew = true
	return f
}

func (f field) readOnly() field {
	f.computed = true
	return f
}

func (f field) secret() field {
	f.sensitive = true
	return f
}

func (f field) password() field {
	f.sensitive = true
	f.writeOnly = true
	return f
}

func (f field) defaults(value string) field {
	f.fallback = value
	return f
}

func (f field) attrType() attr.Type {
	switch f.typ {
	case boolField:
		return types.BoolType
	case intField:
		return types.Int64Type
	case floatField:
		return types.Float64Type
	case stringListField:
		return types.ListType{ElemType: types.StringType}
	case stringMapField:
		return types.MapType{ElemType: types.StringType}
	case objectListField:
		return types.ListType{ElemType: f.objectType()}
	default:
		return types.StringType
	}
}

func (f field) objectType() types.ObjectType {
	attrTypes := map[string]attr.Type{}
	for _, sub := range f.fields {
		attrTypes[sub.name] = sub.attrType()
	}
	return types.ObjectType{AttrTypes: attrTypes}
}

func (f field) schemaAttribute() schema.Attribute {
	required := f.required
	optional := !f.required && !f.computed
	computed := !f.required && !f.writeOnly

	switch f.typ {
	case boolField:
		var modifiers []planmodifier.Bool
		if computed {
			modifiers = append(modifiers, boolplanmodifier.UseStateForUnknown())
		}
		if f.forceNew {
			modifiers = append(modifiers, boolplanmodifier.RequiresReplace())
		}
		return schema.BoolAttribute{MarkdownDescription: f.description, Required: required, Optional: optional, Computed: computed, Sensitive: f.sensitive, PlanModifiers: modifiers}
	case intField:
		var modifiers []planmodifier.Int64
		if computed {
			modifiers = append(modifiers, int64planmodifier.UseStateForUnknown())
		}
		if f.forceNew {
			modifiers = append(modifiers, int64planmodifier.RequiresReplace())
		}
		return schema.Int64Attribute{MarkdownDescription: f.description, Required: required, Optional: optional, Computed: computed, Sensitive: f.sensitive, PlanModifiers: modifiers}
	case floatField:
		var modifiers []planmodifier.Float64
		if computed {
			modifiers = append(modifiers, float64planmodifier.UseStateForUnknown())
		}
		if f.forceNew {
			modifiers = append(modifiers, float64planmodifier.RequiresReplace())
		}
		return schema.Float64Attribute{MarkdownDescription: f.description, Required: required, Optional: optional, Computed: computed, Sensitive: f.sensitive, PlanModifiers: modifiers}
	case stringListField:
		var modifiers []planmodifier.List
		if computed {
			modifiers = append(modifiers, listplanmodifier.UseStateForUnknown())
		}
		if f.forceNew {
			modifiers = append(modifiers, listplanmodifier.RequiresReplace())
		}
		return schema.ListAttribute{MarkdownDescription: f.description, ElementType: types.StringType, Required: required, Optional: optional, Computed: computed, Sensitive: f.sensitive, PlanModifiers: modifiers}
	case stringMapField:
		var modifiers []planmodifier.Map
		if computed {
			modifiers = append(modifiers, mapplanmodifier.UseStateForUnknown())
		}
		if f.forceNew {
			modifiers = append(modifiers, mapplanmodifier.RequiresReplace())
		}
		return schema.MapAttribute{MarkdownDescription: f.description, ElementType: types.StringType, Required: required, Optional: optional, Computed: computed, Sensitive: f.sensitive, PlanModifiers: modifiers}
	case objectListField:
		var modifiers []planmodifier.List
		if computed {
			modifiers = append(modifiers, listplanmodifier.UseStateForUnknown())
		}
		attributes := map[string]schema.Attribute{}
		for _, sub := range f.fields {
			attributes[sub.name] = sub.nestedAttribute()
		}
		return schema.ListNestedAttribute{MarkdownDescription: f.description, NestedObject: schema.NestedAttributeObject{Attributes: attributes}, Required: required, Optional: optional, Computed: computed, PlanModifiers: modifiers}
	default:
		var modifiers []planmodifier.String
		if computed && f.fallback == "" {
			modifiers = append(modifiers, stringplanmodifier.UseStateForUnknown())
		}
		if f.forceNew {
			modifiers = append(modifiers, stringplanmodifier.RequiresReplace())
		}
		attribute := schema.StringAttribute{MarkdownDescription: f.description, Required: required, Optional: optional, Computed: computed, Sensitive: f.sensitive, PlanModifiers: modifiers}
		if f.fallback != "" {
			attribute.Default = stringdefault.StaticString(f.fallback)
		}
		return attribute
	}
}

func (f field) nestedAttribute() schema.Attribute {
	if f.required {
		return schema.StringAttribute{MarkdownDescription: f.description, Required: true}
	}

	switch f.typ {
	case boolField:
		return schema.BoolAttribute{MarkdownDescription: f.description, Optional: true, Computed: true, Default: booldefault.StaticBool(false)}
	default:
		return schema.StringAttribute{MarkdownDescription: f.description, Optional: true, Computed: true, Default: stringdefault.StaticString("")}
	}
}

func (f field) toValue(raw any) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch f.typ {
	case boolField:
		b, _ := raw.(bool)
		return types.BoolValue(b), diags
	case intField:
		n, err := toNumber(raw).Int64()
		if err != nil {
			f64, _ := toNumber(raw).Float64()
			n = int64(f64)
		}
		return types.Int64Value(n), diags
	case floatField:
		f64, _ := toNumber(raw).Float64()
		return types.Float64Value(f64), diags
	case stringListField:
		items, _ := raw.([]any)
		elems := []attr.Value{}
		for _, item := range items {
			elems = append(elems, types.StringValue(toString(item)))
		}
		return types.ListValue(types.StringType, elems)
	case stringMapField:
		items, _ := raw.(map[string]any)
		elems := map[string]attr.Value{}
		for k, v := range items {
			elems[k] = types.StringValue(toString(v))
		}
		return types.MapValue(types.StringType, elems)
	case objectListField:
		items, _ := raw.([]any)
		elems := []attr.Value{}
		for _, item := range items {
			m, _ := item.(map[string]any)
			attrs := map[string]attr.Value{}
			for _, sub := range f.fields {
				value, d := sub.toValue(m[sub.key])
				diags.Append(d...)
				attrs[sub.name] = value
			}
			obj, d := types.ObjectValue(f.objectType().AttrTypes, attrs)
			diags.Append(d...)
			elems = append(elems, obj)
		}
		list, d := types.ListValue(f.objectType(), elems)
		diags.Append(d...)
		return list, diags
	default:
		return types.StringValue(toString(raw)), diags
	}
}

func (f field) fromValue(value attr.Value, existing any) any {
	switch f.typ {
	case boolField:
		return value.(types.Bool).ValueBool()
	case intField:
		return value.(types.Int64).ValueInt64()
	case floatField:
		return value.(types.Float64).ValueFloat64()
	case stringListField:
		res := []string{}
		for _, elem := range value.(types.List).Elements() {
			res = append(res, elem.(types.String).ValueString())
		}
		return res
	case stringMapField:
		res := map[string]string{}
		for k, elem := range value.(types.Map).Elements() {
			res[k] = elem.(types.String).ValueString()
		}
		return res
	case objectListField:
		old := map[string]map[string]any{}
		items, _ := existing.([]any)
		for _, item := range items {
			if m, ok := item.(map[string]any); ok {
				old[toString(m["name"])] = m
			}
		}

		res := []map[string]any{}
		for _, elem := range value.(types.List).Elements() {
			attrs := elem.(types.Object).Attributes()
			name := attrs["name"].(types.String).ValueString()
			m := map[string]any{}
			for k, v := range old[name] {
				m[k] = v
			}
			for _, sub := range f.fields {
				m[sub.key] = sub.fromValue(attrs[sub.name], m[sub.key])
			}
			res = append(res, m)
		}
		return res
	default:
		return value.(types.String).ValueString()
	}
}

func (f field) nullValue() attr.Value {
	switch f.typ {
	case boolField:
		return types.BoolNull()
	case intField:
		return types.Int64Null()
	case floatField:
		return types.Float64Null()
	case stringListField:
		return types.ListNull(types.StringType)
	case stringMapField:
		return types.MapNull(types.StringType)
	case objectListField:
		return types.ListNull(f.objectType())
	default:
		return types.StringNull()
	}
}

func toNumber(raw any) json.Number {
	switch v := raw.(type) {
	case json.Number:
		return v
	case nil:
		return "0"
	default:
		return json.Number(fmt.Sprint(v))
	}
}

func toString(raw any) string {
	switch v := raw.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		return fmt.Sprint(v)
	}
}
