package ctrl

import (
	"context"
	"github.com/Percona-Lab/operator-env/config"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Controller type represent testing environment lifecycle controller.
type Controller struct {
	cli   *client.Client
	cfg   *config.Config
	log   *zerolog.Logger
	plane *Plane
}

// NewController return new instance of Controller type.
func NewController(log *zerolog.Logger, cfg *config.Config) (*Controller, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Error().Err(err).Msg("docker client creation failed")
		return nil, errors.Wrap(err, "docker client creation failed")
	}

	plane := &Plane{
		Version:   "1.13",
		EtcdImg:   "gcr.io/google-containers/etcd:3.2.24",
		MasterImg: "gcr.io/google-containers/hyperkube:v1.13.0",
		ProxyImg:  "gcr.io/google-containers/hyperkube:v1.13.0",
	}

	return &Controller{
		cli:   cli,
		log:   log,
		cfg:   cfg,
		plane: plane,
	}, nil
}

// Up will bring up Kubernetes or Openshift cluster
func (c *Controller) Up(ctx context.Context, platform Platform) error {
	switch platform {

	case Kubernetes:
		c.log.Info().Msg("Starting the Kubernetes cluster")

		if err := c.CreateControlPlane(ctx); err != nil {
			return errors.Wrap(err, "can't create control plane")
		}
		c.log.Info().Msg("control plane created")

	case OpenShift:
		c.log.Info().Msg("Starting the OpenShift cluster")

		// TODO need to implement logic for OpenShift
		return errors.New("the OpenShift platform currently not supported")
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
	//containers, err := c.cli.ContainerList(ctx, types.ContainerListOptions{})
	//if err != nil {
	//	return errors.Wrap(err, "can't fetch containers list")
	//}
	//
	//for _, container := range containers {
	//	if _, ok := c.img[container.Image]; ok {
	//		if err := c.cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
	//			Force:         true,
	//			RemoveVolumes: true,
	//			RemoveLinks:   true,
	//		}); err != nil {
	//			return errors.Wrapf(err, "can't complete cleanup. can't remove container %s", container.ID[:10])
	//		}
	//	}
	//}
	return nil
}
