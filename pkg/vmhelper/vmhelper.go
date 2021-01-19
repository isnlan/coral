package vmhelper

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/snlansky/coral/pkg/utils"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/snlansky/coral/pkg/errors"

	"github.com/snlansky/coral/pkg/net"
	"github.com/snlansky/coral/pkg/protos"
)

type IVM interface {
	Apply(data [][]byte) error
	Delete(data [][]byte) error
	GetNodeIps() ([]string, error)
	GetNsList() ([]string, error)
	GetServiceList(ns string) ([]string, error)
	GetServicePort(ns string, service string) ([]string, error)
	GetDeploymentList(ns string) ([]string, error)
	GetDeploymentStatus(ns string, deployment string) error
	GetNamespacesPods(ns string, label string, filter map[string]string) ([]*protos.Pod, error)
	BuildImage(name string, src string) error
	PushImage(name string, version string) error
	GetRepositoryUrl() string
}

type vm struct {
	cli        *net.Client
	repository string
}

func (v *vm) Apply(data [][]byte) error {
	conn, err := v.cli.Get()
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	cli := protos.NewVMClient(conn.ClientConn)
	for _, d := range data {
		_, err = cli.Apply(v.getContext(), &protos.Data{Data: d})
		if err != nil {
			return errors.Wrapf(err, "apply date: \n %s", string(d))
		}
	}

	return nil
}

func (v *vm) Delete(data [][]byte) error {
	conn, err := v.cli.Get()
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	cli := protos.NewVMClient(conn.ClientConn)

	for i := range data {
		// 资源删除顺序与资源创建顺序相反
		_, err = cli.Delete(v.getContext(), &protos.Data{Data: data[len(data)-1-i]})
		if err != nil {
			return err
		}
	}

	return nil

}

func (v *vm) GetNodeIps() ([]string, error) {
	conn, err := v.cli.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	cli := protos.NewVMClient(conn.ClientConn)
	ips, err := cli.GetNodeIps(v.getContext(), &empty.Empty{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ips.Ips, nil
}

func (v *vm) GetNsList() ([]string, error) {
	conn, err := v.cli.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	cli := protos.NewVMClient(conn.ClientConn)
	list, err := cli.GetNamespacesList(v.getContext(), &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return list.Namespaces, nil
}

func (v *vm) GetServiceList(ns string) ([]string, error) {
	conn, err := v.cli.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	cli := protos.NewVMClient(conn.ClientConn)
	list, err := cli.GetServiceList(v.getContext(), &protos.Namespace{Ns: ns})
	if err != nil {
		return nil, err
	}
	return list.Services, nil
}

func (v *vm) GetServicePort(ns string, service string) ([]string, error) {
	conn, err := v.cli.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	cli := protos.NewVMClient(conn.ClientConn)
	ret, err := cli.GetServicePort(v.getContext(), &protos.RequestServicePort{Ns: ns, Svc: service})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ret.Ports, nil
}

func (v *vm) GetDeploymentList(ns string) ([]string, error) {
	conn, err := v.cli.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	cli := protos.NewVMClient(conn.ClientConn)
	list, err := cli.GetDeploymentList(v.getContext(), &protos.Namespace{Ns: ns})
	if err != nil {
		return nil, err
	}
	return list.Deployments, err
}

func (v *vm) GetDeploymentStatus(ns string, deployment string) error {
	conn, err := v.cli.Get()
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	cli := protos.NewVMClient(conn.ClientConn)
	status, err := cli.GetDeploymentStatus(v.getContext(), &protos.RequestDeploymentStatus{Ns: ns, Name: deployment})
	if err != nil {
		return errors.WithStack(err)
	}

	if len(status.Status) > 0 {
		return errors.New(status.Status[0])
	}

	return nil
}

func (v *vm) GetNamespacesPods(ns string, label string, filter map[string]string) ([]*protos.Pod, error) {
	conn, err := v.cli.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	cli := protos.NewVMClient(conn.ClientConn)
	pods, err := cli.GetNamespacesPods(v.getContext(), &protos.RequestNsPods{
		Ns:     ns,
		Label:  label,
		Filter: filter,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return pods.Pods, nil
}

func (v *vm) BuildImage(name string, src string) error {
	conn, err := v.cli.Get()
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	reader, err := utils.CreateTarStream(src, "Dockerfile")
	if err != nil {
		return errors.WithStack(err)
	}
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.WithStack(err)
	}

	cli := protos.NewVMClient(conn.ClientConn)
	_, err = cli.BuildImage(context.Background(), &protos.RequestBuildImage{Tag: name, Data: bytes})
	return err
}

func (v *vm) PushImage(name string, version string) error {
	conn, err := v.cli.Get()
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	cli := protos.NewVMClient(conn.ClientConn)
	_, err = cli.PushImage(context.Background(), &protos.RequestPushImage{Name: name, Version: version})
	return err
}

func (v *vm) GetRepositoryUrl() string {
	return v.repository
}

func (v *vm) getContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	return ctx
}

func New(url, repository string) (*vm, error) {
	cli, err := net.New(url)
	if err != nil {
		return nil, err
	}
	return &vm{cli: cli, repository: repository}, nil
}
