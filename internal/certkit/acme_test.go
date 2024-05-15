package certkit

import (
	"github.com/qiniu/go-sdk/v7/storage"
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
	provider := NewQiniuChallengeProvider(testBucket, testAccessKey, testSecretKey, &storage.ZoneXinjiapo)
	acme := NewQiniuACME(user, provider)
	certificate, err := acme.ObtainCertificate("test.czyt.tech")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get certificate:%+v \r\n", *certificate)
}
