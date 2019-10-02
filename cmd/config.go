package cmd

import (
	"fmt"
	"log"

	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "ocm config Changes",
	Long:  `Manage configuration file`,
	Run: func(cmd *cobra.Command, args []string) {

		var config struct {
			Default  string
			Clusters map[string]map[string]string
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("Unable to decode config, %v", err)
		}

		prettyConfig, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Print(string(prettyConfig))
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
}
