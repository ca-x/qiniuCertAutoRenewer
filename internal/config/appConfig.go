package config

import (
	"gopkg.in/yaml.v3"
	"io"
)

type AppConfig struct {
	CertUpdatePolicy CertUpdatePolicy `yaml:"cert_update_policy"`
	Qiniu            QiniuConfig      `yaml:"qiniu"`
}

type CertUpdatePolicy struct {
	UpdateBeforeDays int `yaml:"update_before_days"`
}

type QiniuConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`

	CDNConfigs []QiniuCDNConfig `yaml:"cdn_configs"`
}

type QiniuCDNConfig struct {
	Bucket string `yaml:"bucket"`
	Domain string `yaml:"domain"`
	// 华东 z0
	// 华东-浙江 cn-east-2
	// 华北 z1
	// 华南 z2
	// 北美 na0
	// 新加坡 as0

	RegionID string `yaml:"region_id"`
}

func (a *AppConfig) SaveTo(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(a)
}

func LoadFrom(r io.Reader) (*AppConfig, error) {
	conf := &AppConfig{}
	if err := yaml.NewDecoder(r).Decode(conf); err != nil {
		return nil, err
	}
	return conf, nil
}
