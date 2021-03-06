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

package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

// Config is global object that holds all application level variables.
var Config appConfig

type appConfig struct {
	BundleDirPath string `mapstructure:"BUNDLE_DIR_PATH"`
}

// Load loads config from environment variables
func Load() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("PRODCTL")

	viper.BindEnv("BUNDLE_DIR_PATH")

	if err := viper.Unmarshal(&Config); err != nil {
		return err
	}

	if Config.BundleDirPath == "" {
		Config.BundleDirPath = "/bundle"
		if _, err := os.Stat(Config.BundleDirPath); errors.Is(err, os.ErrNotExist) {
			Config.BundleDirPath = "fakes/.bundle"
		}
	}

	return nil
}
