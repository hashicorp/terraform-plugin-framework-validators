// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package ephemeralvalidator provides validators to express relationships
// between multiple attributes of an ephemeral resource. For example, checking that
// multiple attributes are not configured at the same time.
//
// These validators are implemented outside the schema, which may be easier to
// implement in provider code generation situations or suit provider code
// preferences differently than those in the schemavalidator package. Those
// validators start on a starting attribute, where relationships can be
// expressed as absolute paths to others or relative to the starting attribute.
package ephemeralvalidator
