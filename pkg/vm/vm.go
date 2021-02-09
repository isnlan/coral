package vm

import (
	"context"
	"io/ioutil"

	"github.com/isnlan/coral/pkg/utils"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/isnlan/coral/pkg/errors"

	"github.com/isnlan/coral/pkg/protos"
	"github.com/isnlan/coral/pkg/xgrpc"
)

type VM interface {
	Apply(ctx context.Context, data [][]byte) error
	Delete(ctx context.Context, data [][]byte) error
	GetNodeIps(ctx context.Context) ([]string, error)
	GetNsList(ctx context.Context) ([]string, error)
	GetServiceList(ctx context.Context, ns string) ([]string, error)
	GetServicePort(ctx context.Context, ns string, service string) ([]string, error)
	GetDeploymentList(ctx context.Context, ns string) ([]string, error)
	GetDeploymentStatus(ctx context.Context, ns string, deployment string) error
	GetNamespacesPods(ctx context.Context, ns string, label string, filter map[string]string) ([]*protos.Pod, error)
	BuildImage(ctx context.Context, name string, src string) error
	PushImage(ctx context.Context, name string, version string) error
	GetRepositoryUrl(ctx context.Context) string
}

type vmImpl struct {
	client     protos.VMClient
	repository string
}

func (v *vmImpl) Apply(ctx context.Context, data [][]byte) error {
	for _, d := range data {
		_, err := v.client.Apply(ctx, &protos.Data{Data: d})
		if err != nil {
			return errors.Wrapf(err, "apply date: \n %s", string(d))
		}
	}

	return nil
}

func (v *vmImpl) Delete(ctx context.Context, data [][]byte) error {
	for i := range data {
		// 资源删除顺序与资源创建顺序相反
		_, err := v.client.Delete(ctx, &protos.Data{Data: data[len(data)-1-i]})
		if err != nil {
			return err
		}
	}

	return nil

}

func (v *vmImpl) GetNodeIps(ctx context.Context) ([]string, error) {
	ips, err := v.client.GetNodeIps(ctx, &empty.Empty{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ips.Ips, nil
}

func (v *vmImpl) GetNsList(ctx context.Context) ([]string, error) {
	list, err := v.client.GetNamespacesList(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return list.Namespaces, nil
}

func (v *vmImpl) GetServiceList(ctx context.Context, ns string) ([]string, error) {
	list, err := v.client.GetServiceList(ctx, &protos.Namespace{Ns: ns})
	if err != nil {
		return nil, err
	}
	return list.Services, nil
}

func (v *vmImpl) GetServicePort(ctx context.Context, ns string, service string) ([]string, error) {
	ret, err := v.client.GetServicePort(ctx, &protos.RequestServicePort{Ns: ns, Svc: service})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ret.Ports, nil
}

func (v *vmImpl) GetDeploymentList(ctx context.Context, ns string) ([]string, error) {
	list, err := v.client.GetDeploymentList(ctx, &protos.Namespace{Ns: ns})
	if err != nil {
		return nil, err
	}
	return list.Deployments, err
}

func (v *vmImpl) GetDeploymentStatus(ctx context.Context, ns string, deployment string) error {
	status, err := v.client.GetDeploymentStatus(ctx, &protos.RequestDeploymentStatus{Ns: ns, Name: deployment})
	if err != nil {
		return errors.WithStack(err)
	}

	if len(status.Status) > 0 {
		return errors.New(status.Status[0])
	}

	return nil
}

func (v *vmImpl) GetNamespacesPods(ctx context.Context, ns string, label string, filter map[string]string) ([]*protos.Pod, error) {
	pods, err := v.client.GetNamespacesPods(ctx, &protos.RequestNsPods{
		Ns:     ns,
		Label:  label,
		Filter: filter,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return pods.Pods, nil
}

func (v *vmImpl) BuildImage(ctx context.Context, name string, src string) error {
	reader, err := utils.CreateTarStream(src, "Dockerfile")
	if err != nil {
		return errors.WithStack(err)
	}
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = v.client.BuildImage(ctx, &protos.RequestBuildImage{Tag: name, Data: bytes})
	return err
}

func (v *vmImpl) PushImage(ctx context.Context, name string, version string) error {
	_, err := v.client.PushImage(ctx, &protos.RequestPushImage{Name: name, Version: version})
	return err
}

func (v *vmImpl) GetRepositoryUrl(_ context.Context) string {
	return v.repository
}

func New(url, repository string) (*vmImpl, error) {
	cli, err := xgrpc.NewClient(url)
	if err != nil {
		return nil, err
	}

	return &vmImpl{client: protos.NewVMClient(cli), repository: repository}, nil
}
