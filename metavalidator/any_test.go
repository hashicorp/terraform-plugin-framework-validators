package metavalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/metavalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ tfsdk.AttributeValidator = warningValidator{}

type warningValidator struct {
	summary string
	detail  string
}

func (validator warningValidator) Description(_ context.Context) string {
	return "description"
}

func (validator warningValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator warningValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	response.Diagnostics.Append(diag.NewWarningDiagnostic(validator.summary, validator.detail))
}

func TestAnyValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val                    attr.Value
		valueValidators        []tfsdk.AttributeValidator
		expectError            bool
		inspectDiags           bool
		expectedValidatorDiags diag.Diagnostics
	}
	tests := map[string]testCase{
		"Type mismatch": {
			val: types.Int64{Value: 12},
			valueValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(3),
				stringvalidator.LengthAtLeast(5),
			},
			expectError: true,
			// We can't test the diags returned as they are in the /internal/reflect pkg.
			inspectDiags: false,
		},
		"String invalid": {
			val: types.String{Value: "one"},
			valueValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(4),
				stringvalidator.LengthAtLeast(5),
			},
			expectError:  true,
			inspectDiags: true,
			expectedValidatorDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					tftypes.NewAttributePath().WithAttributeName("test"),
					"Invalid Attribute Value Length",
					"String length must be at least 4, got: 3",
				),
				diag.NewAttributeErrorDiagnostic(
					tftypes.NewAttributePath().WithAttributeName("test"),
					"Invalid Attribute Value Length",
					"String length must be at least 5, got: 3",
				),
			},
		},
		"String valid": {
			val: types.String{Value: "one"},
			valueValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(5),
				stringvalidator.LengthAtLeast(3),
			},
			expectError:            false,
			inspectDiags:           true,
			expectedValidatorDiags: diag.Diagnostics{},
		},
		"String invalid in all nested validators": {
			val: types.String{Value: "one"},
			valueValidators: []tfsdk.AttributeValidator{
				metavalidator.All(stringvalidator.LengthAtLeast(6), stringvalidator.LengthAtLeast(3)),
				metavalidator.All(stringvalidator.LengthAtLeast(5), stringvalidator.LengthAtLeast(3)),
			},
			expectError:  true,
			inspectDiags: true,
			expectedValidatorDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					tftypes.NewAttributePath().WithAttributeName("test"),
					"Invalid Attribute Value Length",
					"String length must be at least 6, got: 3",
				),
				diag.NewAttributeErrorDiagnostic(
					tftypes.NewAttributePath().WithAttributeName("test"),
					"Invalid Attribute Value Length",
					"String length must be at least 5, got: 3",
				),
			},
		},
		"String valid in one of the nested validators": {
			val: types.String{Value: "one"},
			valueValidators: []tfsdk.AttributeValidator{
				metavalidator.All(stringvalidator.LengthAtLeast(6), stringvalidator.LengthAtLeast(3)),
				metavalidator.All(stringvalidator.LengthAtLeast(2), stringvalidator.LengthAtLeast(3)),
			},
			expectError:            false,
			inspectDiags:           true,
			expectedValidatorDiags: diag.Diagnostics{},
		},
		"String valid in one of the nested validators with warning": {
			val: types.String{Value: "one"},
			valueValidators: []tfsdk.AttributeValidator{
				metavalidator.All(stringvalidator.LengthAtLeast(6), stringvalidator.LengthAtLeast(3)),
				metavalidator.All(stringvalidator.LengthAtLeast(2), warningValidator{
					summary: "Warning",
					detail:  "Warning",
				}),
			},
			expectError:  false,
			inspectDiags: true,
			expectedValidatorDiags: diag.Diagnostics{
				diag.NewWarningDiagnostic("Warning", "Warning")},
		},
		"String valid in one of the nested validators with warning and warning from failed validation": {
			val: types.String{Value: "one"},
			valueValidators: []tfsdk.AttributeValidator{
				metavalidator.All(stringvalidator.LengthAtLeast(6), warningValidator{
					summary: "Warning",
					detail:  "Warning from failed validation",
				}),
				metavalidator.All(stringvalidator.LengthAtLeast(2), warningValidator{
					summary: "Warning",
					detail:  "Warning from successful validation",
				}),
			},
			expectError:  false,
			inspectDiags: true,
			expectedValidatorDiags: diag.Diagnostics{
				diag.NewWarningDiagnostic("Warning", "Warning from failed validation"),
				diag.NewWarningDiagnostic("Warning", "Warning from successful validation"),
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := tfsdk.ValidateAttributeRequest{
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
				AttributeConfig: test.val,
			}
			response := tfsdk.ValidateAttributeResponse{}
			metavalidator.Any(test.valueValidators...).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}

			if test.inspectDiags {
				if diff := cmp.Diff(response.Diagnostics, test.expectedValidatorDiags); diff != "" {
					t.Errorf("unexpected diags difference: %s", diff)
				}
			}
		})
	}
}
