# Habitat Builder Syncer

This package is intended to sync habitat packages, and their corresponding keys from an upstream
builder, to a target builder.

These values are driven by configuration.

## Usage

For the preferred usage on the habitat package please review the [Habitat package README](./habitat/README.md)

Currently there is only one mode to execute in, continuity as an agent where the sync process happens followed
by a sleep and re-execution.

### Example

This will run the sync process with the given configuration file.

```
bldr_package_sync --config user.toml sync
```

## Configuration

* `interval`: integer value used to determine the amount of time to sleep after running the process
* `log_level`: the level to print logs at (debug, info, warn, error)
* `temp_dir`: the temp directory to stage files in
* `env`: additional environment variables to use when shelling out (could be proxy or ssl, see
  [troubleshooting](#troubleshooting))
* `features`: list of beta features to include
  * options: `PACKAGE_CONSTRAINTS`
* `upstream`: the bldr to pull packages/keys from
  * `url`: the url to the corresponding upstream
* `target`: the bldr to push packages/keys to
  * `url`: the url to the corresponding target
  * `authToken`: the authToken to push to the origins in the target (must have access to _all_
    outlined origins)
* `origin`: list of origins and channels to pull packages from
  * `name`: the name of the origin
  * `channels`: a list of channels to pull/push packages to
* `package`: list of package constraints (this feature is beta)
  * `name`: the name of the origin
  * `contraint`: the constraint on the package

### Full Example

```
interval = 300
env = []
log_level = "info"
temp_dir = "/tmp"

[upstream]
url = "https://bldr.habitat.sh"

[target]
url = ""
authToken = ""

[[origin]]
name = "habitat"
channels = ["stable", "on-prem-stable"]

[[origin]]
name = "core"
channels = ["stable"]
```

#### Package Contraints (Beta)

```
interval = 300
env = []
log_level = "info"
temp_dir = "/tmp"
features = ["PACKAGE_CONSTRAINTS"]

[upstream]
url = "https://bldr.habitat.sh"

[target]
url = ""
authToken = ""

[[origin]]
name = "habitat"
channels = ["stable", "on-prem-stable"]

[[origin]]
name = "core"
channels = ["stable"]

[[package]]
name = "core/hab"
constraint = "< 0.80.0"

[[package]]
name = "core/hab-pkg-export-docker"
constraint = "< 0.80.0"
```

## Troubleshooting

If running on a mac and getting SSL related errors, try appending to the env array in the config
file:

```
env = ["SSL_CERT_FILE=/usr/local/etc/openssl/cert.pem"]
```
