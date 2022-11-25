package listvalidator_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
)

func TestValueNumbersAreValidatorValidateList(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.List
		elementValidators   []validator.Number
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.ListValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(big.NewFloat(1.2)),
					types.NumberValue(big.NewFloat(2.4)),
				},
			),
		},
		"List unknown": {
			val: types.ListUnknown(
				types.NumberType,
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(1.2)),
			},
		},
		"List null": {
			val: types.ListNull(
				types.NumberType,
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(1.2)),
			},
		},
		"List elements invalid": {
			val: types.ListValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(big.NewFloat(1.2)),
					types.NumberValue(big.NewFloat(2.4)),
				},
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(3.6)),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value Match",
					"Attribute test[0] value must be one of: [\"3.6\"], got: 1.2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value Match",
					"Attribute test[1] value must be one of: [\"3.6\"], got: 2.4",
				),
			},
		},
		"List elements invalid for multiple validator": {
			val: types.ListValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(big.NewFloat(1.2)),
					types.NumberValue(big.NewFloat(2.4)),
				},
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(3.6)),
				numbervalidator.OneOf(big.NewFloat(4.8)),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value Match",
					"Attribute test[0] value must be one of: [\"3.6\"], got: 1.2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value Match",
					"Attribute test[0] value must be one of: [\"4.8\"], got: 1.2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value Match",
					"Attribute test[1] value must be one of: [\"3.6\"], got: 2.4",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value Match",
					"Attribute test[1] value must be one of: [\"4.8\"], got: 2.4",
				),
			},
		},
		"List elements valid": {
			val: types.ListValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(big.NewFloat(1.2)),
					types.NumberValue(big.NewFloat(2.4)),
				},
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(1.2), big.NewFloat(2.4)),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.ListRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.val,
			}
			response := validator.ListResponse{}
			listvalidator.ValueNumbersAre(testCase.elementValidators...).ValidateList(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
