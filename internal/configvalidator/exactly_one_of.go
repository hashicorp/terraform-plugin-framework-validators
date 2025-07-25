// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ datasource.ConfigValidator = &ExactlyOneOfValidator{}
var _ provider.ConfigValidator = &ExactlyOneOfValidator{}
var _ resource.ConfigValidator = &ExactlyOneOfValidator{}

// ExactlyOneOfValidator is the underlying struct implementing ExactlyOneOf.
type ExactlyOneOfValidator struct {
	PathExpressions path.Expressions
}

func (v ExactlyOneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v ExactlyOneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Exactly one of these attributes must be configured: %s", v.PathExpressions)
}

func (v ExactlyOneOfValidator) ValidateDataSource(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ExactlyOneOfValidator) ValidateEphemeralResource(ctx context.Context, req ephemeral.ValidateConfigRequest, resp *ephemeral.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ExactlyOneOfValidator) ValidateListResourceConfig(ctx context.Context, req list.ValidateConfigRequest, resp *list.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ExactlyOneOfValidator) ValidateProvider(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ExactlyOneOfValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ExactlyOneOfValidator) Validate(ctx context.Context, config tfsdk.Config) diag.Diagnostics {
	var configuredPaths, unknownPaths path.Paths
	var diags diag.Diagnostics

	for _, expression := range v.PathExpressions {
		matchedPaths, matchedPathsDiags := config.PathMatches(ctx, expression)

		diags.Append(matchedPathsDiags...)

		// Collect all errors
		if matchedPathsDiags.HasError() {
			continue
		}

		for _, matchedPath := range matchedPaths {
			var value attr.Value
			getAttributeDiags := config.GetAttribute(ctx, matchedPath, &value)

			diags.Append(getAttributeDiags...)

			// Collect all errors
			if getAttributeDiags.HasError() {
				continue
			}

			// If value is unknown, it may be null or a value, so we cannot
			// know if the validator should succeed or not. Collect the path
			// path so we use it to skip the validation later and continue to
			// collect all path matching diagnostics.
			if value.IsUnknown() {
				unknownPaths.Append(matchedPath)
				continue
			}

			// If value is null, move onto the next one.
			if value.IsNull() {
				continue
			}

			// Value is known and not null, it is configured.
			configuredPaths.Append(matchedPath)
		}
	}

	// We can always return an error if more than one path was configured.
	if len(configuredPaths) > 1 {
		diags.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			configuredPaths[0],
			v.Description(ctx),
		))
	}

	// If there are unknown values, we cannot know if the validator should
	// succeed or not.
	if len(unknownPaths) > 0 {
		return diags
	}

	// Only return missing attribute configuration when error diagnostics are
	// not present, since they likely represent a provider developer mistake,
	// such as an invalid path expression.
	if len(configuredPaths) == 0 && !diags.HasError() {
		diags.Append(diag.NewErrorDiagnostic(
			"Missing Attribute Configuration",
			v.Description(ctx),
		))
	}

	return diags
}
