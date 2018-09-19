package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

var clientNumber int

func main() {
	clientNumber = 0
	for range time.Tick(3 * time.Second) {
		addClient()
		clientNumber++
	}
}

func makeSpec(image string, externPort int) swarm.ServiceSpec {
	max := uint64(1)

	spec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: "sneks_testing_client_" + strconv.Itoa(clientNumber),
			Labels: map[string]string{
				"com.docker.stack.image":     image,
				"com.docker.stack.namespace": "sneks_testing",
			},
		},
		TaskTemplate: swarm.TaskSpec{
			RestartPolicy: &swarm.RestartPolicy{
				MaxAttempts: &max,
				Condition:   swarm.RestartPolicyConditionNone,
			},
			ContainerSpec: swarm.ContainerSpec{
				Image: image,
				Env:   []string{"GSPORT=" + strconv.Itoa(externPort)},
				Labels: map[string]string{
					"com.docker.stack.image":     image,
					"com.docker.stack.namespace": "sneks_testing",
				},
			},
		},
	}
	return spec
}

func makeOpts() types.ServiceCreateOptions {
	authConfig := types.AuthConfig{
		Username: "parthmehrotra",
		Password: PASSWORD,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	return types.ServiceCreateOptions{
		EncodedRegistryAuth: authStr,
	}
}

func addClient() {
	dockerClient, dockerErr := client.NewEnvClient()
	if dockerErr != nil {
		fmt.Println("DOCKER ERROR")
		fmt.Println(dockerErr)
		return
	}

	createResponse, serviceErr :=
		dockerClient.ServiceCreate(
			context.Background(),
			makeSpec("moneygames/test-client:master", currentPort),
			makeOpts())

	fmt.Println(createResponse)
	if serviceErr != nil {
		fmt.Println("Service ERROR")
		fmt.Println(serviceErr)
		return
	}
}
