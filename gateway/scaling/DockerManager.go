package scaling

import (
	"os/exec"
	"fmt"
	"os"
	"strings"
	"github.com/tkrex/IDS/common/models"
	"net/url"
)
type DockerManager struct {

}

func NewDockerManager() *DockerManager {
	return new(DockerManager)
}


func (dockerManager *DockerManager) StartDomainControllerInstance(parentDomain,ownDomain *models.RealWorldDomain) (*models.DomainController, error) {

	envVariables := dockerManager.buildEnvVariables(parentDomain,ownDomain)
	dockerManager.setEnvVariables(envVariables)

	goPath := os.Getenv("GOPATH")
	dockerFilePath := goPath+"/src/github.com/tkrex/IDS/DockerFiles/domainController"

	if error := os.Chdir(dockerFilePath); error != nil {
		fmt.Fprintln(os.Stderr, "Error starting docker compose instance: ", error)
		return nil ,error
	}

	cmdName := "docker-compose"
	cmdArgs := []string{"-p "+ownDomain.Name,"up","-d",}
	if err := exec.Command(cmdName, cmdArgs...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error starting docker compose instance: ", err)
		return nil ,err
	}
	brokerPort := dockerManager.getContainerPort(envVariables["broker"],"9001/tcp")
	restPort := dockerManager.getContainerPort(envVariables["domainController"],"8080/tcp")
	clusterIP := "10.40.53.21"
	brokerURL,_ := url.Parse("ws://"+ clusterIP + ":" + brokerPort)
	restURL,_ := url.Parse("http://"+ clusterIP + ":" + restPort)
	domainController := models.NewDomainController(restURL,brokerURL, ownDomain)
	return domainController, nil
}

func (dockerManager *DockerManager) buildEnvVariables(parentDomain, ownDomain *models.RealWorldDomain) map[string]string{
	domainControllerName := "domainController
	brokerName := "broker"
	dbName := "db"
	envVariables := make(map[string]string)
	envVariables["domainController"] = domainControllerName
	envVariables["db"] = dbName
	envVariables["broker"] = brokerName
	envVariables["own_domain"] = ownDomain.Name
	envVariables["parent_domain"] = parentDomain.Name
	return envVariables
}

func (dockerManager *DockerManager) setEnvVariables(variables map[string]string) {
	for key, value := range variables {
		if error := os.Setenv(key,value); error != nil {
			fmt.Println(error)
		}
	}
}

func (dockerManager *DockerManager) getContainerPort(containerName string, internalPort string) string {
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "docker"
	cmdArgs := []string{"inspect", "-f","'{{index .NetworkSettings.Ports \""+internalPort+"\" 0 \"HostPort\"}}'", containerName}
	fmt.Println(cmdArgs)
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Println("Error while parsing external port for container: ",containerName)
		os.Exit(1)
	}
	externalPort := string(cmdOut)
	externalPort=strings.Replace(externalPort,"'","",-1)
	externalPort=strings.Replace(externalPort,"\n","",-1)
	return externalPort
}

//TODO: Implement
func (dockerManager *DockerManager) StopDomainControllerInstance(domain *models.RealWorldDomain) error {
	var error error
	return error
}
