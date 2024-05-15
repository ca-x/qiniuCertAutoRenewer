package certkit

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuChallengeProvider struct {
	bucket    string
	accessKey string
	secretKey string
	region    *storage.Region
}

func NewQiniuChallengeProvider(bucket string, accessKey string, secretKey string, region *storage.Region) *QiniuChallengeProvider {
	return &QiniuChallengeProvider{
		bucket:    bucket,
		accessKey: accessKey,
		secretKey: secretKey,
		region:    region,
	}
}

func (p *QiniuChallengeProvider) Present(domain, token, keyAuth string) error {
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

func (p *QiniuChallengeProvider) CleanUp(domain, token, keyAuth string) error {
	return nil
}

type User struct {
	Email        string
	Registration *registration.Resource
	PrivateKey   crypto.PrivateKey
}

func NewUser(email string) (*User, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &User{Email: email, PrivateKey: privateKey}, nil
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.PrivateKey
}

type QiniuACME struct {
	user           *User
	qiniuChallenge *QiniuChallengeProvider
}

func NewQiniuACME(user *User, qiniuChallenge *QiniuChallengeProvider) *QiniuACME {
	return &QiniuACME{
		user:           user,
		qiniuChallenge: qiniuChallenge,
	}
}
func (q *QiniuACME) ObtainCertificate(domain string) (*certificate.Resource, error) {
	config := lego.NewConfig(q.user)
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	q.user.Registration = reg
	err = client.Challenge.SetHTTP01Provider(q.qiniuChallenge)
	if err != nil {
		return nil, err
	}
	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}
	return certificates, nil
}
