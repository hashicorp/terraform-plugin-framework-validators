// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// WarningBool returns a validator which returns a warning diagnostic.
func WarningBool(summary string, detail string) validator.Bool {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningDataSource returns a validator which returns a warning diagnostic.
func WarningDataSource(summary string, detail string) datasource.ConfigValidator {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningFloat64 returns a validator which returns a warning diagnostic.
func WarningFloat64(summary string, detail string) validator.Float64 {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningInt64 returns a validator which returns a warning diagnostic.
func WarningInt64(summary string, detail string) validator.Int64 {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningList returns a validator which returns a warning diagnostic.
func WarningList(summary string, detail string) validator.List {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningMap returns a validator which returns a warning diagnostic.
func WarningMap(summary string, detail string) validator.Map {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningNumber returns a validator which returns a warning diagnostic.
func WarningNumber(summary string, detail string) validator.Number {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningObject returns a validator which returns a warning diagnostic.
func WarningObject(summary string, detail string) validator.Object {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningProvider returns a validator which returns a warning diagnostic.
func WarningProvider(summary string, detail string) provider.ConfigValidator {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningResource returns a validator which returns a warning diagnostic.
func WarningResource(summary string, detail string) resource.ConfigValidator {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningSet returns a validator which returns a warning diagnostic.
func WarningSet(summary string, detail string) validator.Set {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

// WarningString returns a validator which returns a warning diagnostic.
func WarningString(summary string, detail string) validator.String {
	return WarningValidator{
		Summary: summary,
		Detail:  detail,
	}
}

var (
	_ datasource.ConfigValidator = WarningValidator{}
	_ provider.ConfigValidator   = WarningValidator{}
	_ resource.ConfigValidator   = WarningValidator{}
	_ validator.Bool             = WarningValidator{}
	_ validator.Float64          = WarningValidator{}
	_ validator.Int64            = WarningValidator{}
	_ validator.List             = WarningValidator{}
	_ validator.Map              = WarningValidator{}
	_ validator.Number           = WarningValidator{}
	_ validator.Object           = WarningValidator{}
	_ validator.Set              = WarningValidator{}
	_ validator.String           = WarningValidator{}
)

type WarningValidator struct {
	Summary string
	Detail  string
}

func (v WarningValidator) Description(_ context.Context) string {
	return "always returns a warning diagnostic"
}

func (v WarningValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v WarningValidator) ValidateBool(ctx context.Context, request validator.BoolRequest, response *validator.BoolResponse) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateDataSource(ctx context.Context, request datasource.ValidateConfigRequest, response *datasource.ValidateConfigResponse) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateFloat64(ctx context.Context, request validator.Float64Request, response *validator.Float64Response) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateList(ctx context.Context, request validator.ListRequest, response *validator.ListResponse) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateMap(ctx context.Context, request validator.MapRequest, response *validator.MapResponse) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateNumber(ctx context.Context, request validator.NumberRequest, response *validator.NumberResponse) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateObject(ctx context.Context, request validator.ObjectRequest, response *validator.ObjectResponse) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateProvider(ctx context.Context, request provider.ValidateConfigRequest, response *provider.ValidateConfigResponse) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateResource(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateSet(ctx context.Context, request validator.SetRequest, response *validator.SetResponse) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}

func (v WarningValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	response.Diagnostics.AddWarning(v.Summary, v.Detail)
}
