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

package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/RyazanovAlexander/pipeline-manager/command-executor/v1/cmd"
	"github.com/RyazanovAlexander/pipeline-manager/command-executor/v1/config"
)

func main() {
	config.Load()
	logger := log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	rootCmd := cmd.NewRootCmd(logger, os.Args[1:])
	go func() {
		errCh <- rootCmd.ExecuteContext(context.Background())
	}()

	select {
	case <-sigCh:
		logger.Println()
		logger.Println("Interrupt signal received. Finishing the application...")
		return
	case err := <-errCh:
		if err != nil {
			logger.Fatal(err)
		}
	}
}
