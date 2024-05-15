package certkit

import (
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type QiniuACME struct {
	user           *User
	qiniuChallenge *ChallengeProvider
}

func NewQiniuACME(user *User, qiniuChallenge *ChallengeProvider) *QiniuACME {
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
