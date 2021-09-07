package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang/glog"

	hiveOpts "github.com/simo7/protoc-gen-gluecatalog/hive_options"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var protoToHive = map[string]string{
	"bool":    "boolean",
	"int32":   "int",
	"int64":   "bigint",
	"uint32":  "int",
	"uint64":  "bigint",
	"float":   "float",
	"float64": "double",
	"double":  "double",
	"string":  "string",
	"[]byte":  "binary",
	"enum":    "string",
}

func main() {
	var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if f.Generate {
				generateFile(gen, f)
			}
		}
		return nil
	})
}

func getIndent(indentLevel int) string {
	return strings.Repeat(" ", indentLevel*2)
}

func generateFile(gen *protogen.Plugin, file *protogen.File) {
	path, _ := filepath.Split(file.GeneratedFilenamePrefix)

	for _, message := range file.Messages {
		opts := message.Desc.Options()

		if !opts.ProtoReflect().IsValid() {
			continue
		}

		optValue := proto.GetExtension(opts, hiveOpts.E_MessageOpts)
		tableName := optValue.(*hiveOpts.MessageOptions).GetTableName()

		if tableName == "" {
			continue
		}

		filename := path + tableName + ".json"
		g := gen.NewGeneratedFile(filename, file.GoImportPath)

		g.P("[")

		for i, field := range message.Fields {
			g.P("  {")

			generateField(g, field.Desc, 2)

			if i == len(message.Fields)-1 {
				g.P("  }")
				continue
			}
			g.P("  },")
		}

		g.P("]")
	}
}

func generateField(g *protogen.GeneratedFile, field protoreflect.FieldDescriptor, indentLevel int) {
	fieldName := string(field.Name())
	protoKind := field.Kind().String()
	fieldType := generateFieldType(field, fieldName, protoKind)

	g.P(fmt.Sprintf(`%s"name": "%s",`, getIndent(indentLevel), fieldName))
	g.P(fmt.Sprintf(`%s"type": "%s"`, getIndent(indentLevel), fieldType))
}

func generateFieldType(field protoreflect.FieldDescriptor, fieldName string, protoKind string) string {
	fieldType := protoToHive[protoKind]

	if protoKind != "message" && fieldType == "" {
		glog.Fatalf("type %v for field %v was not recognized", protoKind, fieldName)
	}

	opts := field.Options()
	if opts.ProtoReflect().IsValid() {
		optValue := proto.GetExtension(opts, hiveOpts.E_FieldOpts)
		fieldType = optValue.(*hiveOpts.FieldOptions).GetTypeOverride()
	}

	if protoKind == "message" {
		fds := field.Message().Fields()

		if fds.Len() > 0 {
			if isProtoTimestamp(fds.Get(0)) {
				return "timestamp"
			}
		}

		var messageFields []string

		for i := 0; i < fds.Len(); i++ {
			field := fds.Get(i)
			fieldName := string(field.Name())
			protoKind := field.Kind().String()
			structField := fmt.Sprintf("%s:%s",
				fieldName, generateFieldType(field, fieldName, protoKind),
			)
			messageFields = append(messageFields, structField)
		}

		fieldType = fmt.Sprintf("struct<%s>", strings.Join(messageFields, ","))
	}

	if field.IsList() {
		return fmt.Sprintf("array<%s>", fieldType)
	}

	return fieldType
}

func isProtoTimestamp(fd protoreflect.FieldDescriptor) bool {
	return strings.HasPrefix(string(fd.FullName()), "google.protobuf.Timestamp.")
}
