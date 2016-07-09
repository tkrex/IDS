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


func (scalingManager *ScalingManager) StartDomainControllerInstance(domain *models.RealWorldDomain) (*models.DomainController, error) {

	envVariables := scalingManager.buildEnvVariables(domain)
	scalingManager.setEnvVariables(envVariables)

	goPath := os.Getenv("GOPATH")
	dockerFilePath := goPath+"/src/github.com/tkrex/IDS/DockerFiles/domainController"

	var err    error
	cmdName := "cd"
	cmdArgs := []string{dockerFilePath}
	if err = exec.Command(cmdName, cmdArgs...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error changing directory to ", err)
		return nil ,err
	}

	cmdName = "docker-compose"
	cmdArgs = []string{"up","-d"}
	if err = exec.Command(cmdName, cmdArgs...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error starting docker compose instance: ", err)
		return nil ,err
	}
	brokerPort := scalingManager.getContainerPort(envVariables["broker"])
	restPort := scalingManager.getContainerPort(envVariables["domainController"])
	clusterIP := "10.40.53.21"
	brokerURL,_ := url.Parse("ws://"+ clusterIP + ":" + string(brokerPort))
	restURL,_ := url.Parse("http://"+ clusterIP + ":" + string(restPort))
	domainController := models.NewDomainController(restURL,brokerURL,domain)
	return domainController, nil
}





func (scalingManager *ScalingManager) buildEnvVariables(domain *models.RealWorldDomain) map[string]string{
	domainControllerName := "domainController-" + domain.Name
	brokerName := "broker-" + domain.Name
	dbName := "db-"+domain.Name
	envVariables := make(map[string]string)
	envVariables["domainController"] = domainControllerName
	envVariables["db"] = dbName
	envVariables["broker"] = brokerName
	return envVariables
}

func (scalingManager *ScalingManager) setEnvVariables(variables map[string]string) {
	for key, value := range variables {
		cmdName := "export"
		cmdArgs := []string{key+"="+value}
		fmt.Println(cmdArgs)
		if err := exec.Command(cmdName, cmdArgs...).Run; err != nil {
			fmt.Fprintln(os.Stderr, "ERROR Setting Env Variables", err)
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

