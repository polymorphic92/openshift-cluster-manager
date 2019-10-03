package openshift

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// #####################################################################################
// #####################################################################################
// #####################################################################################

// Login into an openshift cluster
func Login(cluster map[string]string, password string) bool {

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

// Project to switch to
func Project(project string) {
	args := [5]string{"oc", "project", project}
	_, err := exec.Command(args[0], args[1:3]...).Output()
	if err != nil {
		log.Fatalf("failed running oc project %s\n", err)
		// if login error call Login
	}
	fmt.Printf("Switched to %v\n", project)
}

// Token of the sigined in user
func Token(endpoint string) string {

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
