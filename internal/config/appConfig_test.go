package config

import (
	"os"
	"testing"
)

func Test_CreateConfig(t *testing.T) {
	config := AppConfig{
		CertUpdatePolicy: CertUpdatePolicy{
			UpdateBeforeDays: 10,
		},
		Qiniu: QiniuConfig{
			AccessKey: "11111",
			SecretKey: "333333",
			CDNConfigs: []QiniuCDNConfig{
				{
					Bucket:   "aaa",
					Domain:   "aaaa.qiniu.com",
					RegionID: "1z",
				},
				{
					Bucket:   "bbb",
					Domain:   "bbb.qiniu.com",
					RegionID: "1z",
				},
				{
					Bucket:   "ccc",
					Domain:   "ccc.qiniu.com",
					RegionID: "1z",
				},
			},
		},
	}
	f, err := os.Create("testdata/config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err = config.SaveTo(f); err != nil {
		t.Fatal(err)
	}
}
