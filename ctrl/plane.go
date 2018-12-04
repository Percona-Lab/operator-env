package ctrl

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/pkg/errors"
	"io"
)

// CreateControlPlane check if Control Plane image exists locally and if not it will pull image from docker hub.
// Then it tries to create Control Plane container and return it ID.
func (c *Controller) CreateControlPlane(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		c.log.Warn().Err(ctx.Err()).Msg("Receive shutdown command")
		return "", ctx.Err()
	default:
	}
	if err := c.checkCPImage(ctx); err != nil {
		return "", errors.Wrap(err, "can't create control plane")
	}
	cpid, err := c.createCPContainer(ctx)
	if err != nil {
		return "", errors.Wrap(err, "can't create control plane container")
	}
	return cpid, nil
}

// createCPContainer tries to create Control Plane container and return it ID.
func (c *Controller) createCPContainer(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		c.log.Warn().Err(ctx.Err()).Msg("Receive shutdown command")
		return "", ctx.Err()
	default:
	}

	cfg := &container.Config{
		Image: "ubuntu", //TODO add control plane image
	}

	hostCfg := &container.HostConfig{}

	netCfg := &network.NetworkingConfig{}

	ctnr, err := c.cli.ContainerCreate(ctx, cfg, hostCfg, netCfg, "op-env-control-plane")
	if err != nil {
		return "", errors.Wrap(err, "can't create control plane")
	}
	return ctnr.ID, nil
}

// checkCPImage check if Control Plane image exists locally and if not it will pull image from docker hub.
func (c *Controller) checkCPImage(ctx context.Context) error {
	images, err := c.cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return errors.Wrap(err, "control plane image checking error")
	}
	for _, image := range images {
		//// TODO add needed image
		//if _, ok := image.Labels["docker.io/library/alpine"]; ok {
		//	return nil
		//}
		c.log.Info().Str("Image ID", image.ID).Msg("Image")
		for k, v := range image.Labels {
			c.log.Info().Str(k, v).Msg("Image")
		}
	}

	// TODO add auth for private repos
	reader, err := c.cli.ImagePull(ctx, "docker.io/library/ubuntu", types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "control plane image checking error")
	}
	io.Copy(c.log, reader)
	return nil
}

//// createCPImage
//func (c *Controller) createCPImage(ctx context.Context) error {
//
//}
