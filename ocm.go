package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("ocm")   // name of config file (without extension)
	viper.AddConfigPath("$HOME") // call multiple times to add many search paths

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	type Configuration struct {
		Default  string
		Clusters map[string]map[string]string
	}

	var config Configuration
	var selectedCluster, selectedDefault, clusterPassword string
	var cluster map[string]string
	var clusters []string

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	//default prompt
	if config.Default != "" {
		cluster = config.Clusters[config.Default]

		defaultPrompt := &survey.Select{
			Message: "Use default cluster (" + config.Default + ")",
			Options: []string{"Yes", "No"},
			Default: "Yes",
		}
		surveyError := survey.AskOne(defaultPrompt, &selectedDefault)

		if surveyError != nil {
			os.Exit(0)

		}
		// the option -q uses the default automagically
	}
	// cluster prompt if select default is no or no default option
	if config.Default == "" || selectedDefault == "No" {

		for clusterName := range config.Clusters {
			clusters = append(clusters, clusterName)
		}

		clusterPrompt := &survey.Select{
			Message: "Choose a cluster",
			Options: clusters,
		}
		surveyError := survey.AskOne(clusterPrompt, &selectedCluster)
		if surveyError != nil {
			os.Exit(0)

		}
		fmt.Printf("Cluster selected: %s\n", selectedCluster)

		cluster = config.Clusters[selectedCluster]

	}
	// password promt
	passPrompt := &survey.Password{
		Message: "password",
	}
	passPromptError := survey.AskOne(passPrompt, &clusterPassword)
	if passPromptError != nil {
		os.Exit(0)

	}
	// fmt.Println(reflect.TypeOf(cluster))
	if clusterLogin(cluster, clusterPassword) {
		fmt.Printf("%s Cluster Login succeeded!!\n", selectedCluster)
	}

	if project, ok := cluster["project"]; ok {
		switchProject(project)
	}

	osToken := openshiftToken(cluster["endpoint"])

	if osToken != "" {
		fmt.Printf("Found Token for: %s\n", cluster["endpoint"])

		loginStatus := registryLogin(cluster, osToken)
		fmt.Printf("Docker Login Status: %v\n", loginStatus)

	}
}

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
