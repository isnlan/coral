package project

import (
	"path/filepath"

	"github.com/isnlan/coral/pkg/utils"

	"github.com/isnlan/coral/pkg/errors"
)

const (
	// 所有链的目录
	BaasChains = "chains"
	// 所有DApp目录
	BaasDApps = "contracts"
	// 合约上传目录
	BaasDAppsUpload = "uploads"
	// 所有项目目录
	BaasProjects = "projects"
	// 所有合约解压目录
	BaasDappsExtract = "dappsextract"
)

type Config struct {
	NfsServer       string `mapstructure:"nfsServer"`       // nfs server ip
	NfsPath         string `mapstructure:"nfsPath"`         // 在baas根目录下nfs共享目录
	LocalDataShared string `mapstructure:"localDataShared"` // 本地数据共享目录
	Template        string `mapstructure:"template"`        // 在baas根目录下fabric k8s模板目录
}

var Project *project

//工程目录
type project struct {
	NfsServer           string
	BaasNfsShared       string
	BaasLocalDataShared string
	BaasTemplatePath    string
}

func Init(cfg *Config) {
	Project = &project{
		NfsServer:           cfg.NfsServer,
		BaasNfsShared:       cfg.NfsPath,
		BaasLocalDataShared: cfg.LocalDataShared,
		BaasTemplatePath:    cfg.Template,
	}
}

func (p *project) LocalChainArtifactPath(chain string) string {
	path := filepath.Join(p.BaasLocalDataShared, BaasChains, chain)
	_ = utils.CreatedDir(path)
	return path
}

func (p *project) NfsArtifactPath(chain string) string {
	path := filepath.Join(p.BaasNfsShared, BaasChains, chain)
	_ = utils.CreatedDir(path)
	return path
}

func (p *project) LocalDappExtractPath() string {
	path := filepath.Join(p.BaasLocalDataShared, BaasDappsExtract)
	_ = utils.CreatedDir(path)
	return path
}

func (p *project) LocalDappProjectExtractPath(dappId string) string {
	return filepath.Join(p.BaasLocalDataShared, BaasDappsExtract, dappId)
}

func (p *project) NfsDappExtractPath(dappId string) string {
	path := filepath.Join(p.BaasNfsShared, BaasDappsExtract, dappId)
	return path
}

func (p *project) LocalUserProjectArtifactPath(user string) string {
	path := filepath.Join(p.BaasLocalDataShared, BaasProjects, user)
	_ = utils.CreatedDir(path)
	return path
}

func (p *project) LocalDappArtifactPath(dapp string) string {
	path := filepath.Join(p.BaasLocalDataShared, BaasDApps)
	_ = utils.CreatedDir(path)

	return filepath.Join(path, dapp)
}

func (p *project) LocalDappUploadPath(dapp string) string {
	path := filepath.Join(p.BaasLocalDataShared, BaasDAppsUpload)
	_ = utils.CreatedDir(path)

	return filepath.Join(path, dapp)
}

func (p *project) LocalRemoveChain(chain string) error {
	err := utils.RemoveDir(filepath.Join(p.BaasLocalDataShared, BaasChains, chain))
	if err != nil {
		return errors.Wrap(err, "remove project dir error")
	}
	return nil
}

func (p *project) LocalTemplatePath(sub ...string) string {
	if len(sub) > 0 {
		return filepath.Join(p.BaasTemplatePath, sub[0])
	} else {
		return p.BaasTemplatePath
	}
}
