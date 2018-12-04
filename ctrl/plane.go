package ctrl

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/pkg/errors"
	"io"
)

type Plane struct {
	Version   string
	EtcdImg   string
	MasterImg string
	ProxyImg  string
}

// CreateControlPlane check if Control Plane image exists locally and if not it will pull image from docker hub.
// Then it tries to create Control Plane container and return it ID.
func (c *Controller) CreateControlPlane(ctx context.Context) error {
	select {
	case <-ctx.Done():
		c.log.Warn().Err(ctx.Err()).Msg("receive shutdown command")
		return ctx.Err()
	default:
	}
	c.log.Info().Msg("Baking control plane...")

	if err := c.planeCleanup(ctx); err != nil {
		return err
	}

	c.log.Info().Msg("Starting Etcd...")
	etcd, err := c.etcd(ctx)
	if err != nil {
		return errors.Wrap(err, "plane: failed to create etcd")
	}
	c.log.Info().Str("ID", etcd).Msg("Etcd created")

	master, err := c.master(ctx)
	if err != nil {
		return errors.Wrap(err, "plane: failed to create master")
	}
	c.log.Info().Str("ID", master).Msg("Master created")

	proxy, err := c.master(ctx)
	if err != nil {
		return errors.Wrap(err, "plane: failed to create proxy")
	}
	c.log.Info().Str("ID", proxy).Msg("Proxy created")

	return nil
}

func (c *Controller) etcd(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		c.log.Warn().Err(ctx.Err()).Msg("receive shutdown command")
		return "", ctx.Err()
	default:
	}

	if err := c.checkImage(ctx, c.plane.EtcdImg); err != nil {
		return "", errors.Wrap(err, "etcd: failed to check image")
	}

	cfg := &container.Config{
		Image: c.plane.EtcdImg,
		Cmd: []string{
			"/usr/local/bin/etcd",
			"start",
			"--addr=127.0.0.1:4001",
			"--bind-addr=0.0.0.0:4001",
			"--data-dir=/var/etcd/data",
		},
	}

	hostCfg := &container.HostConfig{
		NetworkMode: "host",
		RestartPolicy: container.RestartPolicy{
			Name:              "on-failure",
			MaximumRetryCount: 5,
		},
	}
	netCfg := &network.NetworkingConfig{}

	ctnr, err := c.cli.ContainerCreate(ctx, cfg, hostCfg, netCfg, "op-env-etcd")
	if err != nil {
		return "", errors.Wrap(err, "failed to create etcd")
	}
	if err := c.cli.ContainerStart(ctx, ctnr.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.Wrap(err, "failed to start etcd")
	}

	return ctnr.ID, nil
}

func (c *Controller) master(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		c.log.Warn().Err(ctx.Err()).Msg("receive shutdown command")
		return "", ctx.Err()
	default:
	}

	if err := c.checkImage(ctx, c.plane.MasterImg); err != nil {
		return "", errors.Wrap(err, "master: failed to check image")
	}

	cfg := &container.Config{
		Image: c.plane.EtcdImg,
		Cmd: []string{
			"/hyperkube",
			"kubelet",
			"--api_servers=http://localhost:8080",
			"--v=2",
			"--address=0.0.0.0",
			"--enable-server",
			"--hostname_override=127.0.0.1",
			"--config=/etc/kubernetes/manifests",
		},
	}

	hostCfg := &container.HostConfig{
		NetworkMode: "host",
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name:              "on-failure",
			MaximumRetryCount: 5,
		},
	}
	netCfg := &network.NetworkingConfig{}

	ctnr, err := c.cli.ContainerCreate(ctx, cfg, hostCfg, netCfg, "op-env-master")
	if err != nil {
		return "", errors.Wrap(err, "failed to create master")
	}
	if err := c.cli.ContainerStart(ctx, ctnr.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.Wrap(err, "failed to start master")
	}

	return ctnr.ID, nil
}

func (c *Controller) proxy(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		c.log.Warn().Err(ctx.Err()).Msg("receive shutdown command")
		return "", ctx.Err()
	default:
	}

	if err := c.checkImage(ctx, c.plane.MasterImg); err != nil {
		return "", errors.Wrap(err, "proxy: failed to check image")
	}

	cfg := &container.Config{
		Image: c.plane.EtcdImg,
		Cmd: []string{
			"/hyperkube",
			"proxy",
			"--master=http://127.0.0.1:8080",
			"--v=2",
		},
	}

	hostCfg := &container.HostConfig{
		NetworkMode: "host",
		Privileged:  true,
		RestartPolicy: container.RestartPolicy{
			Name:              "on-failure",
			MaximumRetryCount: 5,
		},
	}
	netCfg := &network.NetworkingConfig{}

	ctnr, err := c.cli.ContainerCreate(ctx, cfg, hostCfg, netCfg, "op-env-proxy")
	if err != nil {
		return "", errors.Wrap(err, "failed to create proxy")
	}
	if err := c.cli.ContainerStart(ctx, ctnr.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.Wrap(err, "failed to start proxy")
	}

	return ctnr.ID, nil
}

// checkImage check if Control Plane image exists locally and if not it will pull image from docker hub.
func (c *Controller) checkImage(ctx context.Context, image string) error {
	select {
	case <-ctx.Done():
		c.log.Warn().Err(ctx.Err()).Msg("Receive shutdown command")
		return ctx.Err()
	default:
	}

	c.log.Info().Str("Image", image).Msg("Pulling image")

	// TODO add auth for private repos
	reader, err := c.cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "can't pull image")
	}
	io.Copy(c.log, reader)
	return nil
}

func (c *Controller) planeCleanup(ctx context.Context) error {
	select {
	case <-ctx.Done():
		c.log.Warn().Err(ctx.Err()).Msg("Receive shutdown command")
		return ctx.Err()
	default:
	}
	containers := []string{"op-env-etcd", "op-env-master", "op-env-proxy"}
	for _, cnt := range containers {
		if err := c.cli.ContainerRemove(ctx, cnt, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			c.log.Warn().Err(err).Msg("Cleanup control plane before start")
		}
	}
	return nil
}
