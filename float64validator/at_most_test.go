package float64validator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         tftypes.Value
		f           func(context.Context, tftypes.Value) (attr.Value, error)
		max         float64
		expectError bool
	}
	tests := map[string]testCase{
		"not a number": {
			val:         tftypes.NewValue(tftypes.Bool, true),
			f:           types.BoolType.ValueFromTerraform,
			expectError: true,
		},
		"unknown number": {
			val: tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			f:   types.NumberType.ValueFromTerraform,
			max: 2.00,
		},
		"null number": {
			val: tftypes.NewValue(tftypes.Number, nil),
			f:   types.NumberType.ValueFromTerraform,
			max: 2.00,
		},
		"valid integer as Number": {
			val: tftypes.NewValue(tftypes.Number, 1),
			f:   types.NumberType.ValueFromTerraform,
			max: 2.00,
		},
		"valid integer as Float64": {
			val: tftypes.NewValue(tftypes.Number, 1),
			f:   types.Float64Type.ValueFromTerraform,
			max: 2.00,
		},
		"valid float as Number": {
			val: tftypes.NewValue(tftypes.Number, 1.1),
			f:   types.NumberType.ValueFromTerraform,
			max: 2.00,
		},
		"valid float as Float64": {
			val: tftypes.NewValue(tftypes.Number, 1.1),
			f:   types.Float64Type.ValueFromTerraform,
			max: 2.00,
		},
		"valid float as Number max": {
			val: tftypes.NewValue(tftypes.Number, 2.00),
			f:   types.NumberType.ValueFromTerraform,
			max: 2.00,
		},
		"valid float as Float64 max": {
			val: tftypes.NewValue(tftypes.Number, 2.00),
			f:   types.Float64Type.ValueFromTerraform,
			max: 2.00,
		},
		"too large float as Number": {
			val:         tftypes.NewValue(tftypes.Number, 3.00),
			f:           types.NumberType.ValueFromTerraform,
			max:         2.00,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			ctx := context.TODO()
			val, err := test.f(ctx, test.val)

			if err != nil {
				t.Fatalf("got unexpected error: %s", err)
			}

			request := tfsdk.ValidateAttributeRequest{
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
				AttributeConfig: val,
			}
			response := tfsdk.ValidateAttributeResponse{}
			AtMost(test.max).Validate(ctx, request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
