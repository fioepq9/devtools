package modelgen

import (
	"strings"

	"github.com/fioepq9/devtools/pkg/modelgen"
	"github.com/fioepq9/devtools/pkg/modelgen/builder"
	"github.com/fioepq9/devtools/pkg/util"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

type Flag struct {
	Builder bool
}

func (f *Flag) ParseComment(comment string) *Flag {
	for _, v := range strings.Split(comment, " ") {
		if v == "builder" {
			f.Builder = true
		}
	}
	return f
}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "modelgen",
		Usage: "modelgen [packages]",
		Action: func(c *cli.Context) error {
			gofiles, err := util.AllGoFiles(c.Context, strings.Join(c.Args().Slice(), ","))
			if err != nil {
				return err
			}

			for _, f := range gofiles {
				pkg, v, err := modelgen.GetStructDefinesFromFile(f)
				if err != nil {
					return err
				}
				if len(v) == 0 {
					continue
				}

				needBuilder := make([]modelgen.StructDefine, 0)
				for _, vv := range v {
					var flag Flag
					for _, c := range vv.Comment {
						if strings.HasPrefix(c, "+modelgen") {
							flag.ParseComment(c)
						}
					}
					if flag.Builder {
						needBuilder = append(needBuilder, vv)
					}
				}
				if len(needBuilder) != 0 {
					err = builder.GenBuilder(
						c.Context,
						builder.WithPackage(pkg),
						builder.WithPath(strings.TrimRight(f, ".go")+"_builder.gen.go"),
						builder.WithDefines(lo.SliceToMap(needBuilder, func(item modelgen.StructDefine) (string, builder.StructBuilderDefine) {
							return item.Name, builder.StructBuilderDefine{Fields: item.Fields}
						})),
					)
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
}
