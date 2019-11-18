package openshift

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

// Projects  fetch projects of endpoint
func Projects(endpoint string) []string {

	req, err := http.NewRequest("GET", "https://"+endpoint+"/apis/project.openshift.io/v1/projects", nil)
	req.Header.Add("Authorization", "Bearer "+Token(endpoint))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROU] -", err)
	}

	type projects struct {
		Items []struct {
			Metadata struct {
				Name string
			}
		}
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	fetchedProjects := projects{}

	jsonErr := json.Unmarshal(body, &fetchedProjects)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	options := make([]string, len(fetchedProjects.Items))
	for index, project := range fetchedProjects.Items {

		options[index] = project.Metadata.Name
	}

	return options

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
		if strings.Contains(item.Name, convertURL(endpoint)) {
			if token, ok := item.User["token"]; ok {
				return token
			}
		}
	}
	return ""

}

//SetContext foo
func SetContext(project string, endpoint string, user string) {

	context := project + "/" + convertURL(endpoint) + ":443/" + user
	args := [5]string{"oc", "config", "use-context", context}

	_, err := exec.Command(args[0], args[1:4]...).Output()
	if err != nil {
		log.Fatalf("failed running oc  %s\n", err)
	}
	fmt.Printf("Switched to %v\n", context)

}

func checkLogin(endpoint string, token string) bool {

	req, err := http.NewRequest("GET", "https://"+endpoint+":443/api/v1/namespaces", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROU] -", err)
	}

	if resp.Status == "200" {
		return true
	}

	return false
}

func convertURL(url string) string {
	return strings.ReplaceAll(url, ".", "-")
}
