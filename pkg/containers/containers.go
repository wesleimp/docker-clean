package containers

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func List(client *client.Client, all bool) ([]types.Container, error) {
	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{
		All: all,
	})

	return containers, err
}
