package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	// "github.com/AlecAivazis/survey/v2"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ocm",
	Short: "Manage Openshift Clusters",
	Long:  `Manage Openshift Clusters`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/ocm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra-example" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("ocm")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("error:", err)
	}
}

// #####################################################################################
// #####################################################################################
// #####################################################################################

func clusterLogin(cluster map[string]string, password string) bool {

	cmd := [7]string{"oc", "login", "https://" + cluster["endpoint"], "-u", cluster["user"], "-p", password}
	res, err := exec.Command(cmd[0], cmd[1:7]...).Output()
	if err != nil {
		log.Fatalf("oc login for %s falied:\n\n%s", cluster["endpoint"], res)
	}

	if res != nil {
		return true
	}

	return false
}

func switchProject(project string) {
	args := [5]string{"oc", "project", project}
	_, err := exec.Command(args[0], args[1:3]...).Output()
	if err != nil {
		log.Fatalf("failed running oc project %s\n", err)
	}
	fmt.Printf("Switched to %v\n", project)
}

func openshiftToken(endpoint string) string {

	type openshiftConfig struct {
		Users []struct {
			Name string
			User map[string]string
		}
	}

	args := [5]string{"oc", "config", "view", "-o", "json"}
	out, err := exec.Command(args[0], args[1:5]...).Output()
	if err != nil {
		log.Fatalf("failed  getting openshift configwith %s\n", err)
	}

	osConfig := openshiftConfig{}

	jsonErr := json.Unmarshal(out, &osConfig)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for _, item := range osConfig.Users {
		if strings.Contains(item.Name, strings.ReplaceAll(endpoint, ".", "-")) {
			if token, ok := item.User["token"]; ok {
				return token
			}
		}
	}
	return ""

}

func registryLogin(cluster map[string]string, password string) bool {

	cmd := [7]string{"docker", "login", cluster["docker-registry"], "-u", cluster["user"], "-p", password}
	res, err := exec.Command(cmd[0], cmd[1:7]...).Output()
	if err != nil {
		log.Fatalf("docker login for %s falied:\n\n%s", cluster["endpoint"], res)
	}

	if res != nil {
		return true
	}

	return false
}
