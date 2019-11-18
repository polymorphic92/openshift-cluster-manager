package cmd

import (
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/polymorphic92/openshift-cluster-manager/openshift"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Cluster/project selection",
	Long:  `Select an openshift cluster and project then set it as the current context`,
	Run: func(cmd *cobra.Command, args []string) {

		clusterFlag, _ := cmd.Flags().GetString("cluster")
		projectFlag, _ := cmd.Flags().GetString("project")
		defaultFlag, _ := cmd.Flags().GetBool("default")

		var config struct {
			Default  string
			Clusters map[string]map[string]string
		}
		var cluster map[string]string
		var project string
		err := viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("unable to decode config into struct, %v", err)
		}

		// use default if -d flag is set and the clusterFlag is not set
		if defaultFlag && clusterFlag == "" {
			cluster = config.Clusters[config.Default]
		}

		if clusterFlag != "" && config.Clusters[clusterFlag] != nil {
			cluster = config.Clusters[clusterFlag]
		}

		if cluster == nil {
			var clusters []string
			for clusterName := range config.Clusters {
				clusters = append(clusters, clusterName)
			}
			cluster = config.Clusters[selectFrom(clusters, "Select Cluster")]
		}

		// use default if  d flag is set and the projectFlag is not set
		if defaultFlag && projectFlag == "" {
			project = cluster["project"]
		}

		if projectFlag != "" {
			project = projectFlag
		}

		if project == "" {
			projects := openshift.Projects(cluster["endpoint"])
			project = selectFrom(projects, "Select Project:")
		}
		// fmt.Printf("Selected Cluster: %v\n", cluster["endpoint"])
		// fmt.Printf("Selected Project: %v\n", project)
		openshift.SetContext(project, cluster["endpoint"], cluster["user"])

	},
}

func selectFrom(collection []string, surveyMessage string) string {
	var selected string
	clusterPrompt := &survey.Select{
		Message: surveyMessage,
		Options: collection,
	}
	surveyError := survey.AskOne(clusterPrompt, &selected, survey.WithPageSize(50))
	if surveyError != nil {
		os.Exit(0)

	}
	return selected
}

func init() {
	RootCmd.AddCommand(selectCmd)

	selectCmd.Flags().BoolP("default", "d", false, "Use Config Default")
	selectCmd.Flags().StringP("cluster", "c", "", "Select Default project")
	selectCmd.Flags().StringP("project", "p", "", "Select Default project")

}
