package registry

import (
	"log"
	"os/exec"
)

//Login to docker registry
func Login(cluster map[string]string, password string) bool {

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
