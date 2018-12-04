package ctrl

import (
	"context"
	"github.com/Percona-Lab/operator-env/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Controller type represent testing environment lifecycle controller.
type Controller struct {
	cli *client.Client
	cfg *config.Config
	log *zerolog.Logger
	img map[string]string
}

// NewController return new instance of Controller type.
func NewController(log *zerolog.Logger, cfg *config.Config) (*Controller, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Error().Err(err).Msg("docker client creation failed")
		return nil, errors.Wrap(err, "docker client creation failed")
	}

	// TODO fill up img map with images names and tags

	return &Controller{
		cli: cli,
		log: log,
		cfg: cfg,
		img: make(map[string]string),
	}, nil
}

func (c *Controller) Up(ctx context.Context, platform Platform) error {
	switch platform {

	case Kubernetes:
		c.log.Info().Msg("Starting the Kubernetes cluster")

		cpid, err := c.CreateControlPlane(ctx)
		if err != nil {
			return errors.Wrap(err, "can't create control plane")
		}
		c.log.Info().Str("container ID", cpid).Msg("control plane created")

	case OpenShift:
		c.log.Info().Msg("Starting the OpenShift cluster")

		// TODO need to implement logic for OpenShift

	}
	return nil
}

func (c *Controller) Down(ctx context.Context) error {
	if err := c.cleanup(ctx); err != nil {
		return errors.Wrap(err, "environment shutdown error")
	}
	return nil
}

func (c *Controller) cleanup(ctx context.Context) error {
	containers, err := c.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return errors.Wrap(err, "can't fetch containers list")
	}

	for _, container := range containers {
		if _, ok := c.img[container.Image]; ok {
			if err := c.cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
				Force:         true,
				RemoveVolumes: true,
				RemoveLinks:   true,
			}); err != nil {
				return errors.Wrapf(err, "can't complete cleanup. can't remove container %s", container.ID[:10])
			}
		}
	}
	return nil
}
