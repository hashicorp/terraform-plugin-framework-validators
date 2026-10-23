// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package resourcevalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type testConfigValidator struct {
	called *bool
}

func (v testConfigValidator) Description(context.Context) string {
	return "test"
}

func (v testConfigValidator) MarkdownDescription(context.Context) string {
	return "test"
}

func (v testConfigValidator) ValidateResource(context.Context, resource.ValidateConfigRequest, *resource.ValidateConfigResponse) {
	*v.called = true
}

func TestPreferWriteOnlyAttributeIfValidator(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		warn           bool
		ifFuncError    bool
		expCalled      bool
		expErrorsCount int
	}{
		"warning-enabled": {
			warn:      true,
			expCalled: true,
		},
		"warning-disabled": {},
		"conditional-diagnostic": {
			ifFuncError:    true,
			expErrorsCount: 1,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			called := false
			v := preferWriteOnlyAttributeIfValidator{
				ifFunc: func(_ context.Context, _ resource.ValidateConfigRequest, resp *PreferWriteOnlyAttributeIfFuncResponse) {
					resp.Warn = testCase.warn
					if testCase.ifFuncError {
						resp.Diagnostics.Append(diag.NewErrorDiagnostic("Conditional Error", "error"))
					}
				},
				validator: testConfigValidator{called: &called},
			}
			resp := &resource.ValidateConfigResponse{}

			v.ValidateResource(context.Background(), resource.ValidateConfigRequest{}, resp)

			if called != testCase.expCalled {
				t.Fatalf("expected wrapped validator called to be %t, got %t", testCase.expCalled, called)
			}

			if got := resp.Diagnostics.ErrorsCount(); got != testCase.expErrorsCount {
				t.Fatalf("expected %d error(s), got %d: %v", testCase.expErrorsCount, got, resp.Diagnostics)
			}
		})
	}
}
