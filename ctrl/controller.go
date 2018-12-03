package ctrl

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"io"
)

// Controller type represent testing environment lifecycle controller.
type Controller struct {
	cli  *client.Client
	log  *zerolog.Logger
	imgs map[string]string
}

// NewController return new instance of Controller type.
func NewController(log *zerolog.Logger) (*Controller, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Error().Err(err).Msg("docker client creation failed")
		return nil, errors.Wrap(err, "docker client creation failed")
	}

	// TODO fill up imgs map with images names and tags

	return &Controller{
		cli:  cli,
		log:  log,
		imgs: make(map[string]string),
	}, nil
}

func (c *Controller) Up(ctx context.Context, platform Platform) error {
	switch platform {
	case Kubernetes:
		c.log.Info().Msg("Starting the Kubernetes cluster")
		if err := c.createCP(ctx); err != nil {
			return errors.Wrap(err, "can't create control plane")
		}

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

func (c *Controller) createCP(ctx context.Context) (string, error) {
	ctnrCfg := &container.Config{
		Image: "", //TODO add control plane image

	}
	cntrHostCfg := &container.HostConfig{}
	cntrNetCfg := &network.NetworkingConfig{}

	ctnr, err := c.cli.ContainerCreate(ctx, ctnrCfg, cntrHostCfg, cntrNetCfg, "op-env-control-plane")
	if err != nil {
		return "", errors.Wrap(err, "can't create control plane")
	}
	return ctnr.ID, nil
}

func (c *Controller) createNode(ctx context.Context) error {
	// TODO add node creation
	return nil
}

// checkCPImage check if Control Plane image exists locally.
// It will pull image if it does not exist locally.
func (c *Controller) checkCPImage(ctx context.Context) error {
	images, err := c.cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return errors.Wrap(err, "control plane image checking error")
	}
	for _, image := range images {
		// TODO add needed image
		if _, ok := image.Labels["docker.io/library/alpine"]; ok {
			return nil
		}
	}

	// TODO add auth for private repos
	reader, err := c.cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "control plane image checking error")
	}
	io.Copy(c.log, reader)
	return nil
}

func (c *Controller) cleanup(ctx context.Context) error {
	containers, err := c.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return errors.Wrap(err, "can't fetch containers list")
	}

	for _, container := range containers {
		if _, ok := c.imgs[container.Image]; ok {
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
