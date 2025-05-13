## 0.18.0 (May 13, 2025)

NOTES:

* all: This Go module has been updated to Go 1.23 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.23 release notes](https://go.dev/doc/go1.23) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#276](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/276))

BUG FIXES:

* all: Fixed bug with PreferWriteOnlyAttribute validator not running when client is capable of using write-only attributes. ([#287](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/287))

## 0.17.0 (February 19, 2025)

FEATURES:

* boolvalidator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* dynamicvalidator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* float32validator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* float64validator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* int32validator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* int64validator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* listvalidator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* mapvalidator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* numbervalidator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* objectvalidator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* resourcevalidator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))
* stringvalidator: Added `PreferWriteOnlyAttribute` validator ([#263](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/263))

## 0.16.0 (December 12, 2024)

FEATURES:

* listvalidator: Added `NoNullValues` validator ([#245](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/245))
* mapvalidator:  Added `NoNullValues` validator ([#245](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/245))
* setvalidator:  Added `NoNullValues` validator ([#245](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/245))
* dynamicvalidator: New package which contains `types.Dynamic` specific validators ([#249](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/249))

BUG FIXES:

* Fixed bug with `ConflictsWith` and `AlsoRequires` validators where unknown values would raise invalid diagnostics during `terraform validate`. ([#251](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/251))

## 0.15.0 (October 31, 2024)

FEATURES:

* ephemeralvalidator: Introduce new package with declarative validators for ephemeral resource configurations ([#242](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/242))

## 0.14.0 (October 17, 2024)

NOTES:

* all: This Go module has been updated to Go 1.22 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.22 release notes](https://go.dev/doc/go1.22) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#229](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/229))
* all: Previously, creating validators with invalid data would result in a `nil` value being returned and a panic from `terraform-plugin-framework`. This has been updated to return an implementation diagnostic referencing the invalid data/validator during config validation. ([#235](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/235))

FEATURES:

* boolvalidator: Added `Equals` validator ([#232](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/232))

ENHANCEMENTS:

* all: Implemented parameter interfaces for all value-based validators. This allows these validators to be used with provider-defined functions. ([#235](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/235))

## 0.13.0 (July 09, 2024)

NOTES:

* all: This Go module has been updated to Go 1.21 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.21 release notes](https://go.dev/doc/go1.21) before upgrading. Any consumers building on earlier Go versions may experience errors ([#198](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/198))

FEATURES:

* int32validator: New package which contains int32 type specific validators ([#222](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/222))
* float32validator: New package which contains float32 type specific validators ([#222](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/222))

ENHANCEMENTS:

* listvalidator: Added `ValueInt32sAre` and `ValueFloat32sAre` validators ([#222](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/222))
* mapvalidator: Added `ValueInt32sAre` and `ValueFloat32sAre` validators ([#222](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/222))
* setvalidator: Added `ValueInt32sAre` and `ValueFloat32sAre` validators ([#222](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/222))

BUG FIXES:

* Add the attribute from the request to the at-least-one validator's error message ([#199](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/199))

## 0.12.0 (August 30, 2023)

ENHANCEMENTS:

* boolvalidator: Added `All`, `Any`, and `AnyWithAllWarnings` validators ([#158](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/158))
* datasourcevalidator: Added `All`, `Any`, and `AnyWithAllWarnings` validators ([#158](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/158))
* providervalidator: Added `All`, `Any`, and `AnyWithAllWarnings` validators ([#158](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/158))
* resourcevalidator: Added `All`, `Any`, and `AnyWithAllWarnings` validators ([#158](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/158))

## 0.11.0 (August 03, 2023)

NOTES:

* This Go module has been updated to Go 1.19 per the [Go support policy](https://golang.org/doc/devel/release.html#policy). Any consumers building on earlier Go versions may experience errors. ([#117](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/117))

ENHANCEMENTS:

* int64validator: Added `equalToProductOf` validator ([#129](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/129))

BUG FIXES:

* stringvalidator: Removed double quoting in `Description` returned from `NoneOf`, `NoneOfCaseInsensitive`, `OneOf` and `OneOfCaseInsensitive` validators ([#152](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/152))

## 0.10.0 (February 08, 2023)

ENHANCEMENTS:

* listvalidator: Added `IsRequired` validator ([#107](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/107))
* setvalidator: Added `IsRequired` validator ([#107](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/107))
* objectvalidator: Added `IsRequired` validator ([#107](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/107))

# 0.9.0 (December 20, 2022)

ENHANCEMENTS:

* listvalidator: Added `UniqueValues` validator ([#88](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/88))
* stringvalidator: Added `UTF8LengthAtLeast`, `UTF8LengthAtMost`, and `UTF8LengthBetween` validators ([#87](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/87))

# 0.8.0 (December 13, 2022)

NOTES:

* all: Support terraform-plugin-framework version 1.0.0 types handling ([#83](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/83))

# 0.7.0 (November 30, 2022)

BREAKING CHANGES:

* all: Migrated implementations to support terraform-plugin-framework version 0.17.0 `datasource/schema`, `provider/schema`, and `resource/schema` packages with type-specific validation ([#80](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/80))
* listvalidator: The `ValuesAre` validator has been removed and split into element type-specific validators in the same package, such as `StringValuesAre` ([#80](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/80))
* mapvalidator: The `ValuesAre` validator has been removed and split into element type-specific validators in the same package, such as `StringValuesAre` ([#80](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/80))
* metavalidator: The `All` and `Any` validators have been removed and split into type-specific packages, such as `stringvalidator.Any` ([#80](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/80))
* schemavalidator: The `AlsoRequires`, `AtLeastOneOf`, `ConflictsWith`, and `ExactlyOneOf` validators have been removed and split into type-specific packages, such as `stringvalidator.ConflictsWith` ([#80](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/80))
* setvalidator: The `ValuesAre` validator has been removed and split into element type-specific validators in the same package, such as `StringValuesAre` ([#80](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/80))

FEATURES:

* boolvalidator: New package which contains boolean type specific validators ([#80](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/80))
* objectvalidator: New package which contains object type specific validators ([#80](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/80))

# 0.6.0 (November 17, 2022)

NOTES:
* all: This Go module has been updated for deprecations in terraform-plugin-framework version 0.15.0 ([#72](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/72))
* all: This Go module has been updated to make it compatible with the breaking changes in terraform-plugin-framework version 0.16.0 ([#77](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/77))

BUG FIXES:
* mapvalidator: Updated `KeysAre()` to return all errors instead of just the first ([#74](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/74))

# 0.5.0 (August 30, 2022)

NOTES:

* This Go module has been updated to Go 1.18 per the [Go support policy](https://golang.org/doc/devel/release.html#policy). Any consumers building on earlier Go versions may experience errors. ([#55](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/55))

FEATURES:

* Introduced `datasourcevalidator` package with `AtLeastOneOf()`, `Conflicting()`, `ExactlyOneOf()`, and `RequiredTogether()` validation functions ([#60](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/60))
* Introduced `providervalidator` package with `AtLeastOneOf()`, `Conflicting()`, `ExactlyOneOf()`, and `RequiredTogether()` validation functions ([#60](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/60))
* Introduced `resourcevalidator` package with `AtLeastOneOf()`, `Conflicting()`, `ExactlyOneOf()`, and `RequiredTogether()` validation functions ([#60](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/60))

BUG FIXES:

* all: Included missing attribute path details in error diagnostics since they are currently not output by Terraform ([#61](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/61))

# 0.4.0 (July 20, 2022)

FEATURES:

* Introduced `metavalidator` package with `Any()`, `AnyWithAllWarnings()`, and `All()` validation functions ([#43](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/43))
* Introduced `schemavalidator` package with 4 new validation functions: `RequiredWith()`, `ConflictsWith()`, `AtLeastOneOf()`, `ExactlyOneOf()` ([#32](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/32))

ENHANCEMENTS:

* int64validator: Added `AtLeastSumOf()`, `AtMostSumOf()` and `EqualToSumOf()` validation functions ([#29](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/29))

# 0.3.0 (June 29, 2022)

FEATURES:

* Introduced `listvalidator` package with `ValuesAre()` validation functions ([#37](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/37))
* Introduced `mapvalidator` package with `KeysAre()` and `ValuesAre()` validation functions ([#38](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/38))
* Introduced `numbervalidator` package with `OneOf()` and `NoneOf()` validation functions ([#42](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/42))
* Introduced `setvalidator` package with `ValuesAre()` validation function ([#36](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/36))

ENHANCEMENTS:

* float64validator: Added `OneOf()` and `NoneOf()` validation functions ([#42](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/42))
* int64validator: Added `OneOf()` and `NoneOf()` validation functions ([#42](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/42))
* listvalidator: Added `SizeAtLeast()`, `SizeAtMost()` and `SizeBetween` validation functions ([#41](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/41))
* mapvalidator: Added `SizeAtLeast()`, `SizeAtMost()` and `SizeBetween` validation functions ([#39](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/39))
* setvalidator: Added `SizeAtLeast()`, `SizeAtMost()` and `SizeBetween` validation functions ([#40](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/40))
* stringvalidator: Added `OneOf()` and `NoneOf()` (case sensitive), and `OneOfCaseInsensitive()` and `NoneOfCaseInsensitive()` (case insensitive) validation functions ([#45](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/45))

# 0.2.0 (June 7, 2022)

BREAKING CHANGES:

* Fixed package naming for `int64validator`: was misnamed `validate` ([#25](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/25))

# 0.1.0 (May 25, 2022)

FEATURES:

* Introduced `float64validator` package with `AtLeast()`, `AtMost()`, and `Between()` validation functions ([#18](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/18))
* Introduced `int64validator` package with `AtLeast()`, `AtMost()`, and `Between()` validation functions ([#21](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/21))
* Introduced `stringvalidator.RegexMatches()` validation function ([#23](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/23))
* Introduced `stringvalidator` package with `LengthAtLeast()`, `LengthAtMost()`, and `LengthBetween()` validation functions ([#22](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/22))
