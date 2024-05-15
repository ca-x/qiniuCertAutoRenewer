package certkit

import "testing"

func Test_GetCertList(t *testing.T) {
	mgr := NewCertMgr("11111", "2222")
	info, err := mgr.GetDomainInfo("test.czyt.tech")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v \r\n", info)
	certInfo, err := mgr.GetCertInfo("1111")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v \r\n", certInfo)
}
