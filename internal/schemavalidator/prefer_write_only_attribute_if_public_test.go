// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package schemavalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/dynamicvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
)

func TestPreferWriteOnlyAttributeIfPublicFunctions(t *testing.T) {
	t.Parallel()

	attributeWO := path.MatchRoot("attribute_wo")
	testCases := map[string]func(*bool){
		"bool": func(called *bool) {
			boolvalidator.PreferWriteOnlyAttributeIf(attributeWO, func(context.Context, validator.BoolRequest, *boolvalidator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateBool(context.Background(), validator.BoolRequest{}, &validator.BoolResponse{})
		},
		"dynamic": func(called *bool) {
			dynamicvalidator.PreferWriteOnlyAttributeIf(attributeWO, func(context.Context, validator.DynamicRequest, *dynamicvalidator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateDynamic(context.Background(), validator.DynamicRequest{}, &validator.DynamicResponse{})
		},
		"float32": func(called *bool) {
			float32validator.PreferWriteOnlyAttributeIf(attributeWO, func(context.Context, validator.Float32Request, *float32validator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateFloat32(context.Background(), validator.Float32Request{}, &validator.Float32Response{})
		},
		"float64": func(called *bool) {
			float64validator.PreferWriteOnlyAttributeIf(attributeWO, func(context.Context, validator.Float64Request, *float64validator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateFloat64(context.Background(), validator.Float64Request{}, &validator.Float64Response{})
		},
		"int32": func(called *bool) {
			int32validator.PreferWriteOnlyAttributeIf(attributeWO, func(context.Context, validator.Int32Request, *int32validator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateInt32(context.Background(), validator.Int32Request{}, &validator.Int32Response{})
		},
		"int64": func(called *bool) {
			int64validator.PreferWriteOnlyAttributeIf(attributeWO, func(context.Context, validator.Int64Request, *int64validator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateInt64(context.Background(), validator.Int64Request{}, &validator.Int64Response{})
		},
		"list": func(called *bool) {
			listvalidator.PreferWriteOnlyAttributeIf(attributeWO, func(context.Context, validator.ListRequest, *listvalidator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateList(context.Background(), validator.ListRequest{}, &validator.ListResponse{})
		},
		"map": func(called *bool) {
			mapvalidator.PreferWriteOnlyAttributeIf(attributeWO, func(context.Context, validator.MapRequest, *mapvalidator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateMap(context.Background(), validator.MapRequest{}, &validator.MapResponse{})
		},
		"number": func(called *bool) {
			numbervalidator.PreferWriteOnlyAttributeIf(attributeWO, func(context.Context, validator.NumberRequest, *numbervalidator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateNumber(context.Background(), validator.NumberRequest{}, &validator.NumberResponse{})
		},
		"object": func(called *bool) {
			objectvalidator.PreferWriteOnlyAttributeIf(attributeWO, func(context.Context, validator.ObjectRequest, *objectvalidator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateObject(context.Background(), validator.ObjectRequest{}, &validator.ObjectResponse{})
		},
		"resource": func(called *bool) {
			resourcevalidator.PreferWriteOnlyAttributeIf(path.MatchRoot("attribute"), attributeWO, func(context.Context, resource.ValidateConfigRequest, *resourcevalidator.PreferWriteOnlyAttributeIfFuncResponse) {
				*called = true
			}).ValidateResource(context.Background(), resource.ValidateConfigRequest{}, &resource.ValidateConfigResponse{})
		},
	}

	for name, validate := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			called := false
			validate(&called)

			if !called {
				t.Fatal("expected condition function to be called")
			}
		})
	}
}
