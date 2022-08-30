package int64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleAtMostSumOf() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					// Validate this integer value must be at most the
					// summed integer values of other_attr1 and other_attr2.
					int64validator.AtMostSumOf(path.Expressions{
						path.MatchRoot("other_attr1"),
						path.MatchRoot("other_attr2"),
					}...),
				},
			},
			"other_attr1": {
				Required: true,
				Type:     types.Int64Type,
			},
			"other_attr2": {
				Required: true,
				Type:     types.Int64Type,
			},
		},
	}
}
