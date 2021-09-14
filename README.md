# protoc-gen-hive

A protoc plugin that generates json-formatted Hive schemas from protobuf files. See [examples](./examples).

Json is a convenient format to intereact with the AWS Glue API and create/update tables in the AWS Glue Data Catalog.

## Installation

```bash
go get github.com/simo7/protoc-gen-hive
```

Alternatively clone the repo and build the plugin:

```bash
go build -o bin/protoc-gen-hive .

export PATH=$PWD/bin:$PATH
```

## Usage

```bash
protoc \
    --hive_out=./ \
    --hive_opt=paths=source_relative \
    examples/person.proto
```

## Well-known Protobuf types

Reference: https://developers.google.com/protocol-buffers/docs/reference/google.protobuf.

The following types are supported:

- [x] `google.protobuf.Timestamp`

## Compatibility

It's tested against the new protobuf API `google.golang.org/protobuf` or version `1.4.0` of the legacy API `github.com/golang/protobuf`.
