package modelgen

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type FieldDefine struct {
	Name string
	Type string
}

type StructDefine struct {
	Comment []string
	Name    string
	Fields  []FieldDefine
}

func GetStructDefinesFromFile(file string) (string, []StructDefine, error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return "", nil, errors.Join(err, errors.New("read src file failed"))
	}
	f := token.NewFileSet()
	p, err := parser.ParseFile(f, file, src, parser.ParseComments)
	if err != nil {
		return "", nil, errors.Join(err, errors.New("parse file failed"))
	}
	pkg := p.Name.String()
	ret := make([]StructDefine, 0)
	for _, v := range p.Decls {
		var d StructDefine
		d.Fields = make([]FieldDefine, 0)
		stc, ok := v.(*ast.GenDecl)
		if !ok || stc.Tok != token.TYPE {
			continue
		}
		if stc.Doc != nil {
			d.Comment = strings.Split(strings.TrimRight(stc.Doc.Text(), "\n"), "\n")
		}
		for _, spec := range stc.Specs {
			tp, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			stp, ok := tp.Type.(*ast.StructType)
			if !ok || !stp.Struct.IsValid() {
				continue
			}
			d.Name = tp.Name.Name
			for _, field := range stp.Fields.List {
				for _, name := range field.Names {
					d.Fields = append(d.Fields, FieldDefine{
						Name: name.String(),
						Type: string(src[field.Type.Pos()-1 : field.Type.End()-1]),
					})
				}
			}
		}
		if len(d.Name) == 0 {
			continue
		}

		ret = append(ret, d)
	}
	return pkg, ret, nil
}
