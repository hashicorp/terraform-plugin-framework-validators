package metavalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/metavalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestAllValidator(t *testing.T) {
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
				stringvalidator.LengthAtLeast(3),
				stringvalidator.LengthAtLeast(5),
			},
			expectError:  true,
			inspectDiags: true,
			expectedValidatorDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value Length",
					"String length must be at least 5, got: 3",
				),
			},
		},
		"String valid": {
			val: types.String{Value: "one"},
			valueValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(2),
				stringvalidator.LengthAtLeast(3),
			},
			expectError:            false,
			inspectDiags:           true,
			expectedValidatorDiags: nil,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := tfsdk.ValidateAttributeRequest{
				AttributePath:   path.Root("test"),
				AttributeConfig: test.val,
			}
			response := tfsdk.ValidateAttributeResponse{}
			metavalidator.All(test.valueValidators...).Validate(context.TODO(), request, &response)

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
