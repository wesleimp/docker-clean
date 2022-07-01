package containers

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func List(c *client.Client, all bool) ([]types.Container, error) {
	containers, err := c.ContainerList(context.Background(), types.ContainerListOptions{
		All: all,
	})

	return containers, err
}

func Remove(c *client.Client, containers []types.Container) error {
	ctx := context.Background()
	for _, container := range containers {
		err := c.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			return err
		}
	}

	return nil
}
