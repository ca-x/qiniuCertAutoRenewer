package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"time"
)

type AppConfig struct {
	ACMEConfig       ACMEConfig       `yaml:"acme_config"`
	CertUpdatePolicy CertUpdatePolicy `yaml:"cert_update_policy"`
	Qiniu            QiniuConfig      `yaml:"qiniu"`
	CDNConfigs       []QiniuCDNConfig `yaml:"cdn_configs"`
	DelayPerTask     time.Duration    `yaml:"delay_per_task"`
}

type ACMEConfig struct {
	Email string `yaml:"email"`
}
type CertUpdatePolicy struct {
	CreateCertificateForFailureOnes bool  `yaml:"create_certificate_for_failure_ones"`
	UpdateBeforeDays                int64 `yaml:"update_before_days"`
}

type QiniuConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

type QiniuCDNConfig struct {
	Bucket string `yaml:"bucket"`
	// 华东 z0
	// 华东-浙江 cn-east-2
	// 华北 z1
	// 华南 z2
	// 北美 na0
	// 新加坡 as0

	RegionID   string `yaml:"region_id"`
	Domain     string `yaml:"domain"`
	SSLPort    int    `yaml:"ssl_port"`
	ForceHttps bool   `yaml:"force_https,omitempty"`
}

func (a *AppConfig) SaveTo(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(a)
}

func LoadFrom(r io.Reader) (*AppConfig, error) {
	conf := &AppConfig{}
	if err := yaml.NewDecoder(r).Decode(conf); err != nil {
		return nil, err
	}
	if conf.DelayPerTask == 0 {
		conf.DelayPerTask = 10 * time.Second
	}
	return conf, nil
}
