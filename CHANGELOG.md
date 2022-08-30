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
