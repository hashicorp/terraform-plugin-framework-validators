package float64validator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAtLeastValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         tftypes.Value
		f           func(context.Context, tftypes.Value) (attr.Value, error)
		min         float64
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
			min: 0.90,
		},
		"null number": {
			val: tftypes.NewValue(tftypes.Number, nil),
			f:   types.NumberType.ValueFromTerraform,
			min: 0.90,
		},
		"valid integer as Number": {
			val: tftypes.NewValue(tftypes.Number, 2),
			f:   types.NumberType.ValueFromTerraform,
			min: 0.90,
		},
		"valid integer as Float64": {
			val: tftypes.NewValue(tftypes.Number, 2),
			f:   types.Float64Type.ValueFromTerraform,
			min: 0.90,
		},
		"valid float as Number": {
			val: tftypes.NewValue(tftypes.Number, 2.2),
			f:   types.NumberType.ValueFromTerraform,
			min: 0.90,
		},
		"valid float as Float64": {
			val: tftypes.NewValue(tftypes.Number, 2.2),
			f:   types.Float64Type.ValueFromTerraform,
			min: 0.90,
		},
		"valid float as Number min": {
			val: tftypes.NewValue(tftypes.Number, 0.9),
			f:   types.NumberType.ValueFromTerraform,
			min: 0.90,
		},
		"valid float as Float64 min": {
			val: tftypes.NewValue(tftypes.Number, 0.9),
			f:   types.Float64Type.ValueFromTerraform,
			min: 0.90,
		},
		"too small float as Number": {
			val:         tftypes.NewValue(tftypes.Number, -1.1111),
			f:           types.NumberType.ValueFromTerraform,
			min:         0.90,
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
			AtLeast(test.min).Validate(ctx, request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
