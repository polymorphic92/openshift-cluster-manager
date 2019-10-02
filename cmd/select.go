package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// helloCmd represents the hello command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Cluster selection",
	Long:  `Select an openshift cluster`,
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

		cluster, _ := cmd.Flags().GetString("cluster")
		project, _ := cmd.Flags().GetString("project")

		if cluster == "" {
			fmt.Println("Select Cluster ...")
		} else {
			fmt.Println("Using " + cluster)
		}

		if project == "" {
			fmt.Println("Select Project ...")
		} else {
			fmt.Println("Using " + project)
		}

	},
}

func init() {
	RootCmd.AddCommand(selectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	selectCmd.Flags().StringP("cluster", "c", "", "the cluster to login to")
	selectCmd.Flags().StringP("project", "p", "", "the project within the cluster")

}
