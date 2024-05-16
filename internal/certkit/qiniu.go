package certkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"io"
	"net/http"
	"strconv"
	"time"
)

type CodeErr struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type DomainInfo struct {
	CodeErr
	Name               string    `json:"name"`
	PareDomain         string    `json:"pareDomain"`
	Type               string    `json:"type"`
	Cname              string    `json:"cname"`
	TestURLPath        string    `json:"testURLPath"`
	Protocol           string    `json:"protocol"`
	Platform           string    `json:"platform"`
	GeoCover           string    `json:"geoCover"`
	QiniuPrivate       bool      `json:"qiniuPrivate"`
	OperationType      string    `json:"operationType"`
	OperatingState     string    `json:"operatingState"`
	OperatingStateDesc string    `json:"operatingStateDesc"`
	CreateAt           time.Time `json:"createAt"`
	ModifyAt           time.Time `json:"modifyAt"`
	HTTPS              struct {
		CertID     string `json:"certId"`
		ForceHTTPS bool   `json:"forceHttps"`
	} `json:"https"`
	CouldOperateBySelf bool   `json:"couldOperateBySelf"`
	RegisterNo         string `json:"registerNo"`
}

type Cert struct {
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
	CA         string `json:"ca"`
	Pri        string `json:"pri"`
}

func NewCertFrom(certificate *certificate.Resource) Cert {
	return Cert{
		Name:       certificate.Domain,
		CommonName: certificate.Domain,
		CA:         string(certificate.Certificate),
		Pri:        string(certificate.PrivateKey),
	}
}

type UploadCertResp struct {
	CodeErr
	CertID string `json:"certID"`
}

type CertInfo struct {
	CodeErr
	Cert struct {
		CertID           string    `json:"certid"`
		Name             string    `json:"name"`
		UID              int       `json:"uid"`
		CommonName       string    `json:"common_name"`
		DNSNames         []string  `json:"dnsnames"`
		CreateTime       TimeStamp `json:"create_time"`
		NotBefore        TimeStamp `json:"not_before"`
		NotAfter         TimeStamp `json:"not_after"`
		OrderID          string    `json:"orderid"`
		ProductShortName string    `json:"product_short_name"`
		ProductType      string    `json:"product_type"`
		Encrypt          string    `json:"encrypt"`
		EncryptParameter string    `json:"encryptParameter"`
		Enable           bool      `json:"enable"`
		Ca               string    `json:"ca"`
		Pri              string    `json:"pri"`
	} `json:"cert"`
}

type TimeStamp struct {
	time.Time
}

func (t *TimeStamp) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(i, 0)
	return nil
}

type HTTPSConf struct {
	CertID     string `json:"certid"`
	ForceHttps bool   `json:"forceHttps"`
}

const APIHost = "http://api.qiniu.com"

type CertMgr struct {
	*qbox.Mac
}

func NewCertMgr(accessKey, secretKey string) *CertMgr {
	return &CertMgr{
		Mac: qbox.NewMac(accessKey, secretKey),
	}
}

func (c *CertMgr) makeRequest(method string, path string, body interface{}) (resData []byte, err error) {
	urlStr := fmt.Sprintf("%s%s", APIHost, path)
	reqData, _ := json.Marshal(body)
	req, reqErr := http.NewRequest(method, urlStr, bytes.NewReader(reqData))
	if reqErr != nil {
		err = reqErr
		return
	}
	accessToken, signErr := c.SignRequest(req)
	if signErr != nil {
		err = signErr
		return
	}

	req.Header.Add("Authorization", "QBox "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = respErr
		return
	}
	defer resp.Body.Close()

	resData, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		err = ioErr
		return
	}

	return
}

func (c *CertMgr) GetDomainInfo(domain string) (*DomainInfo, error) {
	b, err := c.makeRequest("GET", "/domain/"+domain, nil)
	if err != nil {
		return nil, err
	}
	info := &DomainInfo{}
	if err = json.NewDecoder(bytes.NewReader(b)).Decode(info); err != nil {
		return nil, err
	}

	if info.Code > 200 {
		return nil, fmt.Errorf("%d: %s", info.Code, info.Error)
	}
	return info, nil
}

func (c *CertMgr) GetCertInfo(certID string) (*CertInfo, error) {
	b, err := c.makeRequest("GET", "/sslcert/"+certID, nil)
	if err != nil {
		return nil, err
	}
	info := &CertInfo{}
	if err = json.NewDecoder(bytes.NewReader(b)).Decode(info); err != nil {
		return nil, err
	}
	if info.Code > 200 {
		return nil, fmt.Errorf("%d: %s", info.Code, info.Error)
	}
	return info, nil
}

func (c *CertMgr) UploadCert(cert Cert) (*UploadCertResp, error) {
	b, err := c.makeRequest("POST", "/sslcert", cert)
	if err != nil {
		return nil, err
	}
	resp := &UploadCertResp{}
	if err = json.NewDecoder(bytes.NewReader(b)).Decode(resp); err != nil {
		return nil, err
	}
	if resp.Code > 200 {
		return nil, fmt.Errorf("%d: %s", resp.Code, resp.Error)
	}
	return resp, nil
}

func (c *CertMgr) UpdateHttpsConf(domain, certID string, forceHttps bool) (*CodeErr, error) {
	b, err := c.makeRequest("PUT", "/domain/"+domain+"/httpsconf", HTTPSConf{
		CertID:     certID,
		ForceHttps: forceHttps,
	})
	if err != nil {
		return nil, err
	}
	resp := &CodeErr{}

	if err = json.NewDecoder(bytes.NewReader(b)).Decode(resp); err != nil {
		return nil, err
	}

	if resp.Code > 200 {
		return nil, fmt.Errorf("%d: %s", resp.Code, resp.Error)
	}
	return resp, nil
}

func (c *CertMgr) DeleteCert(certID string) (*CodeErr, error) {
	b, err := c.makeRequest("DELETE", "/sslcert/"+certID, nil)
	if err != nil {
		return nil, err
	}
	resp := &CodeErr{}
	if err = json.NewDecoder(bytes.NewReader(b)).Decode(resp); err != nil {
		return nil, err
	}

	if resp.Code > 200 {
		return nil, fmt.Errorf("%d: %s", resp.Code, resp.Error)
	}
	return resp, nil
}

func (c *CertMgr) DomainSSLize(domain, certID string) (*CodeErr, error) {
	b, err := c.makeRequest("PUT", "/domain/"+domain+"/sslize", HTTPSConf{
		CertID:     certID,
		ForceHttps: true,
	})
	if err != nil {
		return nil, err
	}
	resp := &CodeErr{}
	if err = json.NewDecoder(bytes.NewReader(b)).Decode(resp); err != nil {
		return nil, err
	}
	if resp.Code > 200 {
		return nil, fmt.Errorf("%d: %s", resp.Code, resp.Error)
	}
	return resp, nil
}
