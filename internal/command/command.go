package command

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"

	"github.com/RyazanovAlexander/prodctl/v1/config"
)

type Command struct {
	Name        string
	Description string
	SubCommands []*Command
	Input       map[string]*string
}

func CreateCommandTree() []*Command {
	fset := token.NewFileSet()
	d, _ := parser.ParseDir(fset, config.Config.BundleDirPath, nil, parser.ParseComments)

	for _, f := range d {
		p := doc.New(f, "./", 0)

		for _, f := range p.Funcs {
			fmt.Println("name: ", f.Name)
			fmt.Println("docs: ", f.Doc)
		}
	}

	var s string
	c1 := &Command{
		Name:        "test",
		Description: "aaaa",
		Input: map[string]*string{
			"namespace": &s,
		},
	}

	return []*Command{
		c1,
	}
}
