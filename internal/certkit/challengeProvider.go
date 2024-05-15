package certkit

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type ChallengeProvider struct {
	bucket    string
	accessKey string
	secretKey string
	region    *storage.Region
}

// NewChallengeProvider  create a qiniu challenge provider
func NewChallengeProvider(bucket string, accessKey string, secretKey string, regionId string) *ChallengeProvider {
	region, b := storage.GetRegionByID(storage.RegionID(regionId))
	if !b {
		return nil
	}
	return &ChallengeProvider{
		bucket:    bucket,
		accessKey: accessKey,
		secretKey: secretKey,
		region:    &region,
	}
}

func (p *ChallengeProvider) Present(domain, token, keyAuth string) error {
	_ = dns01.GetChallengeInfo(domain, keyAuth)
	putPolicy := storage.PutPolicy{
		Scope: p.bucket,
	}
	mac := qbox.NewMac(p.accessKey, p.secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.Region = p.region
	cfg.UseHTTPS = true
	cfg.UseCdnDomains = false
	challengePath := http01.ChallengePath(token)
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": challengePath[1:],
		},
	}
	data := []byte(keyAuth)
	dataLen := int64(len(data))
	err := formUploader.Put(
		context.Background(),
		&ret, upToken,
		challengePath[1:],
		bytes.NewReader(data),
		dataLen,
		&putExtra,
	)
	if err != nil {
		return err
	}
	fmt.Println(ret.Key, ret.Hash)

	return nil
}

func (p *ChallengeProvider) CleanUp(domain, token, keyAuth string) error {
	return nil
}
