package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/polymorphic92/openshift-cluster-manager/openshift"
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
		var cluster string
		err := viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("unable to decode config into struct, %v", err)
		}

		useDefault, _ := cmd.Flags().GetBool("default")

		if useDefault && config.Clusters[config.Default] != nil {
			cluster = config.Default
		} else {
			var clusters []string
			for clusterName := range config.Clusters {
				clusters = append(clusters, clusterName)
			}
			cluster = selectFrom(clusters, "Select Cluster")
		}
		fmt.Println("Select cluster: " + cluster)

		openshift.Project(config.Clusters[cluster]["project"])

		// if cluster["project"] == "" {
		// 	fmt.Println("Select Project ...")
		// }

	},
}

func selectFrom(collection []string, surveyMessage string) string {
	var selected string
	clusterPrompt := &survey.Select{
		Message: surveyMessage,
		Options: collection,
	}
	surveyError := survey.AskOne(clusterPrompt, &selected)
	if surveyError != nil {
		os.Exit(0)

	}
	return selected
}

func init() {
	RootCmd.AddCommand(selectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	// selectCmd.Flags().StringP("cluster", "c", "", "the cluster to login to")
	selectCmd.Flags().BoolP("default", "d", false, "Select Default Cluster")

}
