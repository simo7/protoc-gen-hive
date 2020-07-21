# protoc-gen-gluecatalog

A protoc plugin that generates json-formatted Hive schemas from protobuf files. See [examples](./examples).

Json is a convenient format to intereact with the AWS Glue API and create/update tables in the AWS Glue Data Catalog.

## Installation

```bash
go get github.com/simo7/protoc-gen-gluecatalog
```

Alternatively clone the repo and build the plugin:

```bash
go build -o bin/protoc-gen-hive -ldflags -s .

export PATH=$PWD/bin:$PATH
```

## Usage

```bash
protoc \
    --gluecatalog_out=./ \
    --gluecatalog_opt=paths=source_relative \
    --proto_path=./ \
    examples/person.proto
```

## Compatibility

It's tested against the new protobuf API `google.golang.org/protobuf` or version `1.4.0` of the legacy API `github.com/golang/protobuf`.
