package cmd

import (
	"fmt"
	"log"

	// "github.com/davecgh/go-spew/spew"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// byeCmd represents the bye command
var byeCmd = &cobra.Command{
	Use:   "config",
	Short: "ocm config Changes",
	Long:  `Manage configuration file`,
	Run: func(cmd *cobra.Command, args []string) {

		type Configuration struct {
			Default  string
			Clusters map[string]map[string]string
		}

		var config Configuration

		err := viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("unable to decode config into struct, %v", err)
		}

		prettyConfig, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Print(string(prettyConfig))
	},
}

func init() {
	RootCmd.AddCommand(byeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// byeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// byeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
