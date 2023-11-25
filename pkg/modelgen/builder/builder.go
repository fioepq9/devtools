package builder

import (
	"context"

	"github.com/dave/jennifer/jen"
	"github.com/fioepq9/devtools/pkg/modelgen"
)

func NewBuilderFunc(structBuilderName string, structName string) jen.Code {
	return jen.Func().
		Id("New" + structBuilderName).Params().
		Op("*").Id(structBuilderName).
		Block(
			// return &<buidler>{data: &<struct>{}}
			jen.Return(jen.Op("&").Id(structBuilderName).Values(jen.Id("data").Op(":").Op("&").Id(structName).Values())),
		)
}

func FieldAssignmentFunc(structBuilderName, fieldName, fieldType string) jen.Code {
	return jen.Func().Params(jen.Id("b").Id("*"+structBuilderName)).
		Id(fieldName).
		Params(jen.Id("v").Id(fieldType)).
		Id("*"+structBuilderName).
		Block(
			jen.Id("b").Dot("data").Dot(fieldName).Op("=").Id("v"),
			jen.Return(jen.Id("b")),
		)
}

func BuildFunc(structBuilderName string, structName string) jen.Code {
	return jen.Func().
		Params(jen.Id("b").Id("*" + structBuilderName)).
		Id("Build").Params().
		Id("*" + structName).
		Block(jen.Return(jen.Id("b").Dot("data")))
}

type StructBuilderDefine struct {
	Fields []modelgen.FieldDefine
}

type GenBuilderOptions struct {
	Package string
	Path    string
	Defines map[string]StructBuilderDefine
}

type GenBuilderOption func(*GenBuilderOptions)

func WithPackage(pkg string) GenBuilderOption {
	return func(option *GenBuilderOptions) {
		option.Package = pkg
	}
}

func WithPath(path string) GenBuilderOption {
	return func(option *GenBuilderOptions) {
		option.Path = path
	}
}

func WithDefines(defines map[string]StructBuilderDefine) GenBuilderOption {
	return func(option *GenBuilderOptions) {
		option.Defines = defines
	}
}

func Generate(ctx context.Context, option ...GenBuilderOption) error {
	var opt GenBuilderOptions
	for _, o := range option {
		o(&opt)
	}

	f := jen.NewFile(opt.Package)

	for structName, define := range opt.Defines {
		builderName := structName + "Builder"

		fieldAssignmentFns := make([]jen.Code, 0)
		for _, field := range define.Fields {
			fieldAssignmentFns = append(fieldAssignmentFns, FieldAssignmentFunc(builderName, field.Name, field.Type))
		}

		builderStruct := jen.Type().Id(builderName).Struct(jen.Id("data").Id("*" + structName))

		f.Add(builderStruct, jen.Line())

		f.Add(NewBuilderFunc(builderName, structName), jen.Line())
		for _, fn := range fieldAssignmentFns {
			f.Add(fn, jen.Line())
		}
		f.Add(BuildFunc(builderName, structName), jen.Line())
	}

	return f.Save(opt.Path)
}
