Fly Autoscaler Example
======================

An example of the autoscaler which scales based on the running workflow execution
count in Temporal.


## Usage

To run the Temporal autoscaler example, you'll need a Temporal Cloud account, a
namespace set up, and a `ca.pem` and `ca.key` file.

First, export your environment variables:

```sh
export TEMPORAL_ADDRESS="mynamespace.lyeth.tmprl.cloud:7233"
export TEMPORAL_NAMESPACE="mynamespace.lyeth"
export TEMPORAL_TLS_CERT_DATA=$(<ca.pem)
export TEMPORAL_TLS_KEY_DATA=$(<ca.key)
```

Then build & run the worker:

```sh
$ go install ./cmd/fly-autoscaler-temporal-worker
$ fly-autoscaler-temporal-worker
```

Then build and start the command for starting workflows.

```sh
$ go install ./cmd/fly-autoscaler-execute-workflow
$ fly-autoscaler-execute-workflow
```

