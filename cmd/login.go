package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/polymorphic92/openshift-cluster-manager/openshift"
	"github.com/polymorphic92/openshift-cluster-manager/registry"

	"github.com/AlecAivazis/survey/v2"
)

// helloCmd represents the hello command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Cluster login",
	Long: `Logs into one or more openshift cluster(s).
flags: --user --password will be used against all clusters`,
	Run: func(cmd *cobra.Command, args []string) {

		var config struct {
			Default  string
			Clusters map[string]map[string]string
		}

		inputPassword, _ := cmd.Flags().GetString("password")
		var password string
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("Unable to decode config, %v", err)
		}

		if len(args) > 0 {
			for _, cluster := range args {

				selectedCluster := config.Clusters[cluster]
				if selectedCluster != nil {
					fmt.Println("logging into cluster" + selectedCluster["endpoint"])
					password = promtPass(inputPassword)
					if openshift.Login(selectedCluster, password) {
						fmt.Printf("%s Cluster Login succeeded!!\n", cluster)

						if registry.Login(selectedCluster, openshift.Token(selectedCluster["endpoint"])) {
							fmt.Printf("%s Docker Login succeeded!!\n", cluster)
						}
					}
				}

			}
		} else {

			for name := range config.Clusters {

				fmt.Println("logging into cluster: " + config.Clusters[name]["endpoint"])
				password = promtPass(inputPassword)
				if openshift.Login(config.Clusters[name], password) {
					fmt.Printf("%s Cluster Login succeeded!!\n", name)
					if registry.Login(config.Clusters[name], openshift.Token(config.Clusters[name]["endpoint"])) {
						fmt.Printf("%s Docker Login succeeded!!\n", name)
					}
				}
			}

		}

	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("user", "u", "", "The Cluster User")
	loginCmd.Flags().StringP("password", "p", "", "Password to the Cluster User")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func promtPass(input string) string {

	var password string
	if input == "" {
		fmt.Printf("\n")
		passwordPrompt := &survey.Password{
			Message: "password",
		}
		passPromptError := survey.AskOne(passwordPrompt, &password)
		if passPromptError != nil {
			os.Exit(0)
		}
		fmt.Printf("\n")
	} else {
		password = input
	}

	return password
}
