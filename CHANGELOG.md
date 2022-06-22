# 0.3.0 (unreleased)

FEATURES:
* Introduced `mapvalidator` package with `ValuesAre()` validation function ([#13](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/13))

# 0.2.0 (June 7, 2022)

BREAKING CHANGES:

* Fixed package naming for `int64validator`: was misnamed `validate` ([#25](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/25))

# 0.1.0 (May 25, 2022)

FEATURES:

* Introduced `float64validator` package with `AtLeast()`, `AtMost()`, and `Between()` validation functions ([#18](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/18))
* Introduced `int64validator` package with `AtLeast()`, `AtMost()`, and `Between()` validation functions ([#21](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/21))
* Introduced `stringvalidator.RegexMatches()` validation function ([#23](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/23))
* Introduced `stringvalidator` package with `LengthAtLeast()`, `LengthAtMost()`, and `LengthBetween()` validation functions ([#22](https://github.com/hashicorp/terraform-plugin-framework-validators/issues/22))
