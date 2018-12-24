# bldr_package_sync

This package is intended to sync habitat packages, and their corresponding keys from an upstream
builder, to a target builder.

These values are driven by configuration.

## Maintainers

* Skyler Layne (<skyler.layne@indellient.com>)

## Type of Package

This is a service package

## Usage

This package is intended to be run in isolation or along side a depot instance.

## Update Strategies

The recommended update strategy for this packages is `at-once` that way when ever new releases are
published they can be consumed right away.

Checkout [the update strategy documentation](https://www.habitat.sh/docs/using-habitat/#update-strategy)
for information on the
strategies Habitat supports.

### Configuration Updates

The best way to interact with this packages is via user tomls located in
`/hab/user/bldr_package_sync/config/user.toml`, for an example please have a look at the
[default.toml](./default.toml)

Checkout the [configuration update](https://www.habitat.sh/docs/using-habitat/#configuration-updates)
documentation for more
information on what configuration updates are and how they are executed.
