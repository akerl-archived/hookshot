hookshot
=========

[![Build Status](https://img.shields.io/travis/com/akerl/hookshot.svg)](https://travis-ci.com/akerl/hookshot)
[![GitHub release](https://img.shields.io/github/release/akerl/hookshot.svg)](https://github.com/akerl/hookshot/releases)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-green.svg)](https://tldrlegal.com/license/mit-license)

Webhook invoker. Designed to replace [dock0/scheduled_build](https://github.com/dock0/scheduled_build) for hitting URLs when run.

## Usage

1. Create a Lambda with the payload.zip generated in the "Installation" section below.
2. In the Environment Variables for the Lambda, set `S3_BUCKET` and `S3_KEY` to refer to an S3 bucket and S3 key where you will store the configuration file. The bucket/file must be readable by the Lambda.
3. Create a file at that bucket/key with the configuration:

```
---
targets:
  example-target:
    method: POST
    url: https://registry.hub.docker.com/u/dock0/arch/trigger/blahblah/
  other-example:
    url: https://example.com
```

You can list as many targets as you'd like. The `method` setting is optional, and defaults to `GET`.

## Installation

The methods below describe how to create a payload.zip that can be used for AWS Lambdas.

### Official build process

This requires that you have Docker installed and running. It will launch a Docker b
uild container, build the binary, and create a zip file for loading into AWS Lambda
. The zip file can be found at `./pkg/payload.zip`.

```
make
```

### Local pkgforge build

This doesn't require Docker but does require that you have [the pkgforge gem](https://github.com/akerl/pkgforge) installed. It builds a zip file at `./pkg/payload.zip`

```
pkgforge build
```

### Local manual build

This method has no deps other than golang, make, and zip. You have to manually create the zip file.

```
make local
cp ./bin/hookshot_linux ./main
zip payload.zip ./main
```

## License

hookshot is released under the MIT License. See the bundled LICENSE file for details.
