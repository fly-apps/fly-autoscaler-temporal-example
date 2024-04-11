Fly Autoscaler Example for Temporal Cloud
=========================================

An example of the autoscaler which scales based on the running workflow
execution count in [Temporal Cloud](https://temporal.io/cloud).

## Loading cert & key secrets

Temporal Cloud uses mutual TLS for authentication which requires PEM-encoded
certificate and key files. In order to use the public `fly-autoscaler` Docker
images, we need to load the file data via [fly secrets][] which are available
as environment variables to the application.

If you're running Linux, the easiest way to load the certificate (`ca.pem`) and
key (`ca.key`) file data as secrets is to use command substitution:

```sh
fly secrets set --stage FAS_TEMPORAL_CERT_DATA="$(<ca.pem)"
```

```sh
fly secrets set --stage FAS_TEMPORAL_KEY_DATA="$(<ca.key)"
```

## Usage

To run this set of example apps, you'll need to deploy the `fly.scaler.toml` and
`fly.worker.toml` as separate applications first. Then you'll need to build and
run `./cmd/fly-autoscaler-execute-workflow` to begin generating Temporal
workflow executions:

```sh
$ go install ./cmd/fly-autoscaler-execute-workflow
```

```sh
# Export environment variables for the tool to use.
export TEMPORAL_ADDRESS="mynamespace.lyeth.tmprl.cloud:7233"
export TEMPORAL_NAMESPACE="mynamespace.lyeth"
export TEMPORAL_TLS_CERT_DATA=$(<ca.pem)
export TEMPORAL_TLS_KEY_DATA=$(<ca.key)

# Executes one new workflow every 1 second, each workflow lasts for 2 seconds.
$ fly-autoscaler-execute-workflow -count 1 -every 1s -duration 2s
```

