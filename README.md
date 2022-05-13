[![PkgGoDev](https://pkg.go.dev/badge/github.com/hashicorp/terraform-plugin-framework-validators)](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework-validators)

# Terraform Plugin Framework Validators

terraform-plugin-framework-validators is a Go module containing common use case validators for [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) types. It aims to provide generic type validator functionality that should be applicable to the broader framework-based provider ecosystem.

This Go module is not intended to define all possible validations. In particular, many validators that relate to specific string formats, encodings, and other specifics should instead be implemented separately in custom attribute types (e.g. a type implementing the `attr.TypeWithValidators` interface). Many of these custom types can be discovered by a conventional repository naming prefix of `terraform-plugin-framework-type-`.

## Terraform Plugin Framework Compatibility

This Go module is typically kept up to date with the latest `terraform-plugin-framework` releases to ensure all validator functionality is available.

## Go Compatibility

This Go module follows `terraform-plugin-framework` Go compatibility.

Currently that means Go **1.17** must be used when developing and testing code.

## Contributing

See [`.github/CONTRIBUTING.md`](https://github.com/hashicorp/terraform-plugin-framework-validators/blob/main/.github/CONTRIBUTING.md)

## License

[Mozilla Public License v2.0](https://github.com/hashicorp/terraform-plugin-framework-validators/blob/main/LICENSE)
