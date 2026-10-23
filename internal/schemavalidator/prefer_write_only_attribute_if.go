// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package schemavalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ validator.Bool    = PreferWriteOnlyAttributeIf{}
	_ validator.Dynamic = PreferWriteOnlyAttributeIf{}
	_ validator.Float32 = PreferWriteOnlyAttributeIf{}
	_ validator.Float64 = PreferWriteOnlyAttributeIf{}
	_ validator.Int32   = PreferWriteOnlyAttributeIf{}
	_ validator.Int64   = PreferWriteOnlyAttributeIf{}
	_ validator.List    = PreferWriteOnlyAttributeIf{}
	_ validator.Map     = PreferWriteOnlyAttributeIf{}
	_ validator.Number  = PreferWriteOnlyAttributeIf{}
	_ validator.Object  = PreferWriteOnlyAttributeIf{}
	_ validator.String  = PreferWriteOnlyAttributeIf{}
)

// PreferWriteOnlyAttributeIfFuncResponse is the response type for a
// PreferWriteOnlyAttributeIfFunc.
type PreferWriteOnlyAttributeIfFuncResponse struct {
	// Diagnostics report errors or warnings related to the conditional logic.
	Diagnostics diag.Diagnostics

	// Warn enables the PreferWriteOnlyAttribute validation.
	Warn bool
}

// PreferWriteOnlyAttributeIf implements conditional PreferWriteOnlyAttribute validation.
type PreferWriteOnlyAttributeIf struct {
	WriteOnlyAttribute path.Expression
	BoolFunc           func(context.Context, validator.BoolRequest, *PreferWriteOnlyAttributeIfFuncResponse)
	DynamicFunc        func(context.Context, validator.DynamicRequest, *PreferWriteOnlyAttributeIfFuncResponse)
	Float32Func        func(context.Context, validator.Float32Request, *PreferWriteOnlyAttributeIfFuncResponse)
	Float64Func        func(context.Context, validator.Float64Request, *PreferWriteOnlyAttributeIfFuncResponse)
	Int32Func          func(context.Context, validator.Int32Request, *PreferWriteOnlyAttributeIfFuncResponse)
	Int64Func          func(context.Context, validator.Int64Request, *PreferWriteOnlyAttributeIfFuncResponse)
	ListFunc           func(context.Context, validator.ListRequest, *PreferWriteOnlyAttributeIfFuncResponse)
	MapFunc            func(context.Context, validator.MapRequest, *PreferWriteOnlyAttributeIfFuncResponse)
	NumberFunc         func(context.Context, validator.NumberRequest, *PreferWriteOnlyAttributeIfFuncResponse)
	ObjectFunc         func(context.Context, validator.ObjectRequest, *PreferWriteOnlyAttributeIfFuncResponse)
	StringFunc         func(context.Context, validator.StringRequest, *PreferWriteOnlyAttributeIfFuncResponse)
}

func (v PreferWriteOnlyAttributeIf) Description(ctx context.Context) string {
	return PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.Description(ctx)
}

func (v PreferWriteOnlyAttributeIf) MarkdownDescription(ctx context.Context) string {
	return PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.MarkdownDescription(ctx)
}

func preferWriteOnlyAttributeIf[Request any](ctx context.Context, req Request, f func(context.Context, Request, *PreferWriteOnlyAttributeIfFuncResponse)) PreferWriteOnlyAttributeIfFuncResponse {
	resp := PreferWriteOnlyAttributeIfFuncResponse{}
	f(ctx, req, &resp)
	return resp
}

func (v PreferWriteOnlyAttributeIf) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.BoolFunc)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateBool(ctx, req, resp)
	}
}

func (v PreferWriteOnlyAttributeIf) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.DynamicFunc)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateDynamic(ctx, req, resp)
	}
}

func (v PreferWriteOnlyAttributeIf) ValidateFloat32(ctx context.Context, req validator.Float32Request, resp *validator.Float32Response) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.Float32Func)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateFloat32(ctx, req, resp)
	}
}

func (v PreferWriteOnlyAttributeIf) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.Float64Func)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateFloat64(ctx, req, resp)
	}
}

func (v PreferWriteOnlyAttributeIf) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.Int32Func)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateInt32(ctx, req, resp)
	}
}

func (v PreferWriteOnlyAttributeIf) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.Int64Func)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateInt64(ctx, req, resp)
	}
}

func (v PreferWriteOnlyAttributeIf) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.ListFunc)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateList(ctx, req, resp)
	}
}

func (v PreferWriteOnlyAttributeIf) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.MapFunc)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateMap(ctx, req, resp)
	}
}

func (v PreferWriteOnlyAttributeIf) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.NumberFunc)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateNumber(ctx, req, resp)
	}
}

func (v PreferWriteOnlyAttributeIf) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.ObjectFunc)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateObject(ctx, req, resp)
	}
}

func (v PreferWriteOnlyAttributeIf) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	ifResp := preferWriteOnlyAttributeIf(ctx, req, v.StringFunc)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		PreferWriteOnlyAttribute{WriteOnlyAttribute: v.WriteOnlyAttribute}.ValidateString(ctx, req, resp)
	}
}
