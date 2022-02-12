/*
MIT License

Copyright The prodctl Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package command

import (
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"

	"github.com/RyazanovAlexander/prodctl/v1/config"
)

type Command struct {
	Name          string
	Description   string
	DirectoryPath string
	SubCommands   []*Command
	Input         map[string]*string
}

func CreateCommandTree() *Command {
	rootCommand := &Command{
		Name:        "root",
		SubCommands: []*Command{},
	}

	createCommandTree(config.Config.BundleDirPath, rootCommand)
	return rootCommand
}

func createCommandTree(curDir string, command *Command) {
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, curDir, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		pkgDesc := doc.New(pkg, "./", 0)

		for _, funcDesc := range pkgDesc.Funcs {
			if strings.Contains(funcDesc.Doc, "prodctl: ignore") {
				continue
			}

			cmd := strings.ToLower(funcDesc.Name)

			text := strings.Split(funcDesc.Doc, "Params:\n")
			description := strings.Replace(text[0], "\n", "", -1)

			inputs := map[string]*string{}
			if len(text) > 1 {
				params := text[1]
				paramsSlice := strings.Split(params, "\n")
				for i := 0; i < len(paramsSlice)-1; i++ {
					param := strings.ReplaceAll(strings.Split(paramsSlice[i], ":")[0], " ", "")

					var ephemeral string
					inputs[param] = &ephemeral
				}
			}

			command.SubCommands = append(command.SubCommands, &Command{
				Name:          cmd,
				Description:   description,
				DirectoryPath: curDir,
				Input:         inputs,
			})
		}
	}

	dirs, err := ioutil.ReadDir(curDir)
	if err != nil {
		panic(err)
	}

	for _, fileInfo := range dirs {
		if fileInfo.IsDir() {
			fParts := strings.Split(fileInfo.Name(), ".")
			name := fParts[0]
			if len(fParts) > 1 {
				name = fParts[1]
			}

			dirCommand := &Command{
				Name: name,
			}

			command.SubCommands = append(command.SubCommands, dirCommand)
			createCommandTree(curDir+"/"+fileInfo.Name(), dirCommand)
		}
	}
}
