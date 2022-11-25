package numbervalidator_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
)

func TestAllValidatorValidateNumber(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val        types.Number
		validators []validator.Number
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: types.NumberValue(big.NewFloat(1.2)),
			validators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(3)),
				numbervalidator.OneOf(big.NewFloat(5)),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value Match",
					"Attribute test value must be one of: [\"3\"], got: 1.2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value Match",
					"Attribute test value must be one of: [\"5\"], got: 1.2",
				),
			},
		},
		"valid": {
			val: types.NumberValue(big.NewFloat(1.2)),
			validators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(1.2)),
				numbervalidator.NoneOf(big.NewFloat(2.4)),
			},
			expected: nil,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := validator.NumberRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.NumberResponse{}
			numbervalidator.All(test.validators...).ValidateNumber(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
