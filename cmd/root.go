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

package cmd

import (
	"log"

	c "github.com/RyazanovAlexander/prodctl/v1/internal/command"
	"github.com/RyazanovAlexander/prodctl/v1/internal/executor"
	"github.com/spf13/cobra"
)

var globalUsage = `A console utility that manages the deployment, upgrade and removal of the selected bundle resources.

Common actions for prodctl:

- prodctl [directory] ... [directoryN] command [param1] ... [paramN]

Examples:

- prodctl version
- prodctl deploy --namespace test --release first
- prodctl environment deploy --clusterName dev --resourceGroup devrg
- prodctl release engine remove --namespace test --release first
- prodctl release test run --category smoke --name helloWorld
`

// The path to the bundle directory, passed as a command line argument.
var bundleDirPath string

// DefaultBundleDirectoryPath is the path to the default bundle directory.
const DefaultBundleDirectoryPath = "."

// NewRootCmd creates new root cmd.
func NewRootCmd(logger *log.Logger, args []string) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "prodctl",
		Short: "Product bundle control",
		Long:  globalUsage,
	}

	rootCommand := c.CreateCommandTree()
	CreateCommands(rootCommand, cobraCmd, logger)

	cobraCmd.AddCommand(
		newVersionCmd(logger),
	)

	return cobraCmd
}

func CreateCommands(command *c.Command, cobraCmd *cobra.Command, logger *log.Logger) {
	for i := 0; i < len(command.SubCommands); i++ {
		subCommand := command.SubCommands[i]
		cmd := &cobra.Command{
			Use:   subCommand.Name,
			Short: subCommand.Description,
			Run: func(cmd *cobra.Command, args []string) {
				executor.RunCommand(subCommand, logger)
			},
		}

		for key := range subCommand.Input {
			cmd.PersistentFlags().StringVar(subCommand.Input[key], key, "", "")
		}

		cobraCmd.AddCommand(cmd)
		CreateCommands(subCommand, cmd, logger)
	}
}
