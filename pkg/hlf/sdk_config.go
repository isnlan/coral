package hlf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type SDKKey struct {
	Path string `yaml:"path"`
}

type SDKCert struct {
	Path string `yaml:"path"`
}

type SDKTLSClient struct {
	Key  SDKKey  `yaml:"key"`
	Cert SDKCert `yaml:"cert"`
}

type SDKTLSCACerts struct {
	Path   string       `yaml:"path"`
	Client SDKTLSClient `yaml:"client"`
}

type SDKRegistrar struct {
	EnrollID     string `yaml:"enrollId"`
	EnrollSecret string `yaml:"enrollSecret"`
}

type Node struct {
	URL         string         `yaml:"url"`
	GrpcOptions SDKGrpcOptions `yaml:"grpcOptions"`
	TLSCACerts  SDKTLSCACerts  `yaml:"tlsCACerts"`
}

type SDKLogging struct {
	Level string `yaml:"level"`
}

type SDKCryptoconfig struct {
	Path string `yaml:"path"`
}

type SDKCryptoStore struct {
	Path string `yaml:"path"`
}

type SDKCredentialStore struct {
	Path        string         `yaml:"path"`
	CryptoStore SDKCryptoStore `yaml:"cryptoStore"`
}

type SDKTLSCerts struct {
	SystemCertPool bool `yaml:"systemCertPool"`
}

type SDKGrpcOptions struct {
	SslTargetNameOverride string `yaml:"ssl-target-name-override"`
	KeepAliveTime         string `yaml:"keep-alive-time"`
	KeepAliveTimeout      string `yaml:"keep-alive-timeout"`
	KeepAlivePermit       bool   `yaml:"keep-alive-permit"`
	FailFast              bool   `yaml:"fail-fast"`
	AllowInsecure         bool   `yaml:"allow-insecure"`
}

type SDKEntity struct {
	Pattern                             string `yaml:"pattern"`
	URLSubstitutionExp                  string `yaml:"urlSubstitutionExp"`
	SslTargetOverrideURLSubstitutionExp string `yaml:"sslTargetOverrideUrlSubstitutionExp"`
	MappedHost                          string `yaml:"mappedHost"`
}

type SDKOrg struct {
	MspID                  string   `yaml:"mspid"`
	CryptoPath             string   `yaml:"cryptoPath"`
	Peers                  []string `yaml:"peers"`
	CertificateAuthorities []string `yaml:"certificateAuthorities"`
}

type SDKClient struct {
	Organization    string             `yaml:"organization"`
	Logging         SDKLogging         `yaml:"logging"`
	Cryptoconfig    SDKCryptoconfig    `yaml:"cryptoconfig"`
	CredentialStore SDKCredentialStore `yaml:"credentialStore"`
	TLSCerts        SDKTLSCerts        `yaml:"tlsCerts"`
}

type SDKChannel struct {
	Peers map[string]interface{} `yaml:"peers"`
}

type SDKCertificateAuthorities struct {
	URL        string        `yaml:"url"`
	TLSCACerts SDKTLSCACerts `yaml:"tlsCACerts"`
	Registrar  SDKRegistrar  `yaml:"registrar"`
}

type SDKEntityMatchers struct {
	Peer                 []SDKEntity `yaml:"peer"`
	Orderer              []SDKEntity `yaml:"orderer"`
	CertificateAuthority []SDKEntity `yaml:"certificateAuthority"`
}

type SDKConfig struct {
	Client                 SDKClient                            `yaml:"client"`
	Channels               map[string]SDKChannel                `yaml:"channels"`
	Organizations          map[string]SDKOrg                    `yaml:"organizations"`
	Orderers               map[string]Node                      `yaml:"orderers"`
	Peers                  map[string]Node                      `yaml:"peers"`
	CertificateAuthorities map[string]SDKCertificateAuthorities `yaml:"certificateAuthorities"`
	EntityMatchers         SDKEntityMatchers                    `yaml:"entityMatchers"`
}

func NewSDKConfig(d []byte) (*SDKConfig, error) {
	config := new(SDKConfig)
	err := yaml.Unmarshal(d, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (sdk *SDKConfig) Into(useIP bool) (*ClientConfig, error) {
	var config ClientConfig

	config.CryptoConfig = CryptoConfig{
		Family:    "ecdsa",
		Algorithm: "P256-SHA256",
		Hash:      "SHA2-256",
	}

	config.Orderers = map[string]OrdererConfig{}
	for k, v := range sdk.Orderers {
		data, err := ioutil.ReadFile(v.TLSCACerts.Path)
		if err != nil {
			return nil, err
		}

		ord := OrdererConfig{
			Host:               v.URL + ":7050",
			UseTLS:             !v.GrpcOptions.AllowInsecure,
			Cert:               string(data),
			ServerNameOverride: v.GrpcOptions.SslTargetNameOverride,
		}
		if useIP {
			for _, e := range sdk.EntityMatchers.Orderer {
				if e.Pattern == k {
					ord.Host = e.URLSubstitutionExp
				}
			}
		}

		config.Orderers[k] = ord
	}

	config.Peers = map[string]PeerConfig{}
	for k, v := range sdk.Peers {
		data, err := ioutil.ReadFile(v.TLSCACerts.Path)
		if err != nil {
			return nil, err
		}

		peer := PeerConfig{
			Host:               v.URL + ":7051",
			UseTLS:             !v.GrpcOptions.AllowInsecure,
			Cert:               string(data),
			ServerNameOverride: v.GrpcOptions.SslTargetNameOverride,
		}
		if useIP {
			for _, e := range sdk.EntityMatchers.Peer {
				if e.Pattern == k {
					peer.Host = e.URLSubstitutionExp
				}
			}
		}

		config.Peers[k] = peer
	}

	config.EventPeers = map[string]PeerConfig{}
	for k, v := range config.Peers {
		config.EventPeers[k] = v
		break
	}

	return &config, nil
}
