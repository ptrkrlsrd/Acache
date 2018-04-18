// Copyright © 2018 Petter Karlsrud petterkarlsrud@me.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/bolt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/ptrkrlsrd/acache/pkg/acache"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var store acache.Store

const DBName = "acache.db"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "acache",
	Short: "Simple API cacher and server",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initDB)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/acache/acache.json)")
}

func initDB() {
	home, err := homedir.Dir()
	db, err := bolt.Open(fmt.Sprintf("%s/.config/acache/%s", home, DBName), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	store = acache.NewCache(db)

}

func configPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return home + ".config/acache/", nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".config/acache" (without extension).
		configPath, err := configPath()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
