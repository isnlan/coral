package dockervm

import (
	"bytes"
	"io"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/snlansky/coral/pkg/errors"
)

type DockerVM struct {
	client *docker.Client
	config *docker.AuthConfiguration
}

func New(config *docker.AuthConfiguration) *DockerVM {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	return &DockerVM{client: client, config: config}
}

//  buildContext is tar file
func (s *DockerVM) BuildImage(tag string, reader io.Reader) (string, error) {
	var buf bytes.Buffer
	opts := docker.BuildImageOptions{
		Name:                tag,
		Dockerfile:          "Dockerfile",
		NoCache:             true,
		SuppressOutput:      true,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		InputStream:         reader,
		OutputStream:        &buf,
	}

	err := s.client.BuildImage(opts)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(buf.Bytes()), nil
}

// tag必须小写字母+数字
// src目录必须包含Dockerfile
func (s *DockerVM) BuildImageBySource(tag, src string) (string, error) {
	var buf bytes.Buffer
	opts := docker.BuildImageOptions{
		Name:                tag,
		Dockerfile:          "Dockerfile",
		NoCache:             true,
		SuppressOutput:      true,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		OutputStream:        &buf,
		ContextDir:          src,
	}

	err := s.client.BuildImage(opts)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(buf.Bytes()), nil
}

func (s *DockerVM) PushImage(name string, tag string) (string, error) {
	var buf bytes.Buffer
	pushOpts := docker.PushImageOptions{
		Name:         name,
		Tag:          tag,
		Registry:     s.config.ServerAddress,
		OutputStream: &buf,
	}
	err := s.client.PushImage(pushOpts, *s.config)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}
