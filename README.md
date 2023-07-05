# ARCHIVED

Because we decided against using Airbyte at my workplace, I haven't had any need to maintain this,
so I'm glad others have taken the mantel. I hope this has served as a good starting point for
him and others to continue working on this. For now, please use this version instead of mine as
it's more actively managed.

### [josephjohncox/terraform-provider-aribyte](https://github.com/josephjohncox/terraform-provider-airbyte)

# Terraform Provider Airbyte (Terraform Plugin Framework)

Terraform Provider for [Airbyte](https://airbyte.io/).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18
- [Local Airbyte Instance](https://docs.airbyte.com/deploying-airbyte/local-deployment/)

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

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

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
