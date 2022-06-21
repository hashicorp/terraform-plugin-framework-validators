package mapvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.AttributeValidator = keysAreValidator{}

// keysAreValidator validates that each map key validates against each of the value validators.
type keysAreValidator struct {
	keyValidators []tfsdk.AttributeValidator
}

// Description describes the validation in plain text formatting.
func (v keysAreValidator) Description(ctx context.Context) string {
	var descriptions []string
	for _, validator := range v.keyValidators {
		descriptions = append(descriptions, validator.Description(ctx))
	}

	return fmt.Sprintf("key must satisfy all validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v keysAreValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v keysAreValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	elems, ok := validateMap(ctx, req, resp)
	if !ok {
		return
	}

	for k := range elems {
		request := tfsdk.ValidateAttributeRequest{
			AttributePath:   req.AttributePath.WithElementKeyString(k),
			AttributeConfig: types.String{Value: k},
			Config:          req.Config,
		}

		for _, validator := range v.keyValidators {
			validator.Validate(ctx, request, resp)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}
}

func KeysAre(keyValidators ...tfsdk.AttributeValidator) tfsdk.AttributeValidator {
	return keysAreValidator{
		keyValidators: keyValidators,
	}
}
