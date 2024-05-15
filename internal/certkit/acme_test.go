package certkit

import (
	"testing"
)

const (
	testBucket    = "test"
	testAccessKey = "test"
	testSecretKey = "test"
)

func TestQiniuACME_ObtainCertificate(t *testing.T) {
	user, err := NewUser("cert@czyt.tech")
	if err != nil {
		t.Error(err)
	}
	provider := NewChallengeProvider(testBucket, testAccessKey, testSecretKey, "z1")
	acme := NewQiniuACME(user, provider)
	certificate, err := acme.ObtainCertificate("test.czyt.tech")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get certificate:%+v \r\n", *certificate)
}
