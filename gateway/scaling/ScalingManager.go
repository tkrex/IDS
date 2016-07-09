package scaling

import (
	"os/exec"
	"fmt"
	"os"
	"strconv"
	"github.com/tkrex/IDS/common/models"
	"net/url"
)
type ScalingManager struct {

}

func NewScalingManager() *ScalingManager {
	return new(ScalingManager)
}


func (scalingManager *ScalingManager) StartDomainControllerInstance(parentDomain,ownDomain *models.RealWorldDomain) (*models.DomainController, error) {

	envVariables := scalingManager.buildEnvVariables(parentDomain,ownDomain)
	scalingManager.setEnvVariables(envVariables)

	goPath := os.Getenv("GOPATH")
	dockerFilePath := goPath+"/src/github.com/tkrex/IDS/DockerFiles/domainController"

	if error := os.Chdir(dockerFilePath); error != nil {
		fmt.Fprintln(os.Stderr, "Error starting docker compose instance: ", error)
		return nil ,error
	}

	cmdName := "docker-compose"
	cmdArgs := []string{"up","-d"}
	if err := exec.Command(cmdName, cmdArgs...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error starting docker compose instance: ", err)
		return nil ,err
	}
	brokerPort := scalingManager.getContainerPort(envVariables["broker"])
	restPort := scalingManager.getContainerPort(envVariables["domainController"])
	clusterIP := "10.40.53.21"
	brokerURL,_ := url.Parse("ws://"+ clusterIP + ":" + string(brokerPort))
	restURL,_ := url.Parse("http://"+ clusterIP + ":" + string(restPort))
	domainController := models.NewDomainController(restURL,brokerURL, ownDomain)
	return domainController, nil
}





func (scalingManager *ScalingManager) buildEnvVariables(parentDomain, ownDomain *models.RealWorldDomain) map[string]string{
	domainControllerName := "domainController-" + ownDomain.Name
	brokerName := "broker-" + ownDomain.Name
	dbName := "db-"+ ownDomain.Name
	envVariables := make(map[string]string)
	envVariables["domainController"] = domainControllerName
	envVariables["db"] = dbName
	envVariables["broker"] = brokerName
	envVariables["own_domain"] = ownDomain.Name
	envVariables["parent_domain"] = parentDomain.Name
	return envVariables
}

func (scalingManager *ScalingManager) setEnvVariables(variables map[string]string) {
	for key, value := range variables {
		if error := os.Setenv(key,value); error != nil {
			fmt.Println(error)
		}
	}
}

func (scalingManager *ScalingManager) getContainerPort(containerName string) int {
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "docker"
	cmdArgs := []string{"inspect", "-f","'{{index .NetworkSettings.Ports \"9001/tcp\" 0 \"HostPort\"}}'", containerName}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
		os.Exit(1)
	}
	port,_ := strconv.Atoi(string(cmdOut))
	return port
}

