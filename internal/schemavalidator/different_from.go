package schemavalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// This type of validator must satisfy all types.
var (
	_ validator.Bool    = DifferentFromValidator{}
	_ validator.Float64 = DifferentFromValidator{}
	_ validator.Int64   = DifferentFromValidator{}
	_ validator.List    = DifferentFromValidator{}
	_ validator.Map     = DifferentFromValidator{}
	_ validator.Number  = DifferentFromValidator{}
	_ validator.Object  = DifferentFromValidator{}
	_ validator.Set     = DifferentFromValidator{}
	_ validator.String  = DifferentFromValidator{}
)

// DifferentFromValidator is the underlying struct implementing DifferentFrom.
type DifferentFromValidator struct {
	PathExpressions path.Expressions
}

type DifferentFromValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type DifferentFromValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func (dv DifferentFromValidator) Description(ctx context.Context) string {
	return dv.MarkdownDescription(ctx)
}

func (dv DifferentFromValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that if an attribute is set, these don't share the same value: %q", dv.PathExpressions)
}

func (dv DifferentFromValidator) Validate(ctx context.Context, req DifferentFromValidatorRequest, resp *DifferentFromValidatorResponse) {
	// If attribute configuration is null, there is nothing else to validate
	if req.ConfigValue.IsNull() {
		return
	}

	expressions := req.PathExpression.MergeExpressions(dv.PathExpressions...)

	for _, expression := range expressions {
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)

		resp.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		for _, mp := range matchedPaths {
			// If the user specifies the same attribute this validator is applied to,
			// also as part of the input, skip it
			if mp.Equal(req.Path) {
				continue
			}

			var mpVal attr.Value
			diags := req.Config.GetAttribute(ctx, mp, &mpVal)
			resp.Diagnostics.Append(diags...)

			// Collect all errors
			if diags.HasError() {
				continue
			}

			// Delay validation until all involved attribute have a known value
			if mpVal.IsUnknown() {
				return
			}

			// Null attributes can't collide with the attribute undergoing validation
			if mpVal.IsNull() {
				return
			}

			if req.ConfigValue.Equal(mpVal) {
				resp.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					req.Path,
					fmt.Sprintf("Attribute %q cannot have the same value as %q", req.Path, mp),
				))
			}
		}
	}
}

func (dv DifferentFromValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := DifferentFromValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &DifferentFromValidatorResponse{}

	dv.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (dv DifferentFromValidator) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	validateReq := DifferentFromValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &DifferentFromValidatorResponse{}

	dv.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (dv DifferentFromValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := DifferentFromValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &DifferentFromValidatorResponse{}

	dv.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (dv DifferentFromValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := DifferentFromValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &DifferentFromValidatorResponse{}

	dv.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (dv DifferentFromValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := DifferentFromValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &DifferentFromValidatorResponse{}

	dv.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (dv DifferentFromValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	validateReq := DifferentFromValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &DifferentFromValidatorResponse{}

	dv.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (dv DifferentFromValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := DifferentFromValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &DifferentFromValidatorResponse{}

	dv.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (dv DifferentFromValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := DifferentFromValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &DifferentFromValidatorResponse{}

	dv.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (dv DifferentFromValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := DifferentFromValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &DifferentFromValidatorResponse{}

	dv.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
