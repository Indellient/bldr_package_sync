# Habitat Builder Syncer

This package is intended to sync habitat packages, and their corresponding keys from an upstream
builder, to a target builder.

These values are driven by configuration.

## Configuration

* `interval`: integer value used to determine the amount of time to sleep after running the process
* `upstream`: the bldr to pull packages/keys from
  * `url`: the url to the corresponding upstream
* `target`: the bldr to push packages/keys to
  * `url`: the url to the corresponding target
  * `authToken`: the authToken to push to the origins in the target (must have access to _all_
    outlined origins)
* `origin`: list of origins and channels to pull packages from
  * `name`: the name of the origin
  * `channels`: a list of channels to pull/push packages to

### Full Example

```
interval = 300

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

## Troubleshooting

If running on a mac and getting SSL related errors, try appending to the env array in the config
file:

```
env = ["SSL_CERT_FILE=/usr/local/etc/openssl/cert.pem"]
```
