// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

// Package listresourcevalidator provides validators to express relationships
// between multiple attributes of a list config. For example, checking that
// multiple attributes are not configured at the same time.
//
// These validators are implemented outside the schema, which may be easier to
// implement in provider code generation situations or suit provider code
// preferences differently than those in the schemavalidator package. Those
// validators start on a starting attribute, where relationships can be
// expressed as absolute paths to others or relative to the starting attribute.
package listresourcevalidator
