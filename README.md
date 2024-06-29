# Terraform Provider for All Quiet
This is the terraform provider for the modern incident management and escalation platform [All Quiet](https://allquiet.app).

![Test](https://github.com/AllQuietApp/allquiet-terraform-provider-internal/actions/workflows/test.yml/badge.svg)

## Documentation
Documentation: https://registry.terraform.io/providers/allquiet/allquiet/latest/docs

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

This repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

## Testing

To run the acceptance tests you'll need a real API Key from All Quiet. You can set the key in your shell:

```shell
export ALLQUIET_ENDPOINT=https://localhost:7061/api/public/v1
export ALLQUIET_API_KEY=test
```

Run integration tests

```shell
make testacc
```

