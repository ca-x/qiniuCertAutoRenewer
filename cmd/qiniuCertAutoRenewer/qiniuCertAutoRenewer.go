package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"qiniuCertAutoRenewer/internal/certkit"
	"qiniuCertAutoRenewer/internal/config"
	"time"
)

const appVersion = "七牛证书自动续期工具 by czyt v1.0.1"

var (
	logger = slog.Default()
)

func main() {
	var cfg = flag.String("c", "config.yaml", "the config file")
	var version = flag.Bool("v", true, "show version")
	flag.Parse()
	if !*version {
		fmt.Println(appVersion)
		return
	}
	f, err := os.Open(*cfg)
	if err != nil {
		logger.Error("open config file error:", err)
		return
	}
	defer f.Close()
	conf, err := config.LoadFrom(f)
	if err != nil {
		logger.Error("load config file error:", err)
		return
	}
	handleJob(conf)
}

func handleJob(appConfig *config.AppConfig) error {
	// get the cert config
	user, err := certkit.NewUser(appConfig.ACMEConfig.Email)
	if err != nil {
		logger.Error("new user error:", err)
		return err
	}

	certMgr := certkit.NewCertMgr(appConfig.Qiniu.AccessKey, appConfig.Qiniu.SecretKey)

process:
	for _, cdnConfig := range appConfig.CDNConfigs {
		needCreateCert := false
		// check domain cert expire
		tlsEndpoint := certkit.NewTLSEndpoint(cdnConfig.Domain, cdnConfig.SSLPort)
		tlsCerts, err := tlsEndpoint.GetCertificates()
		if err != nil {
			logger.Error("get certificates error", err)
			if appConfig.CertUpdatePolicy.CreateCertificateForFailureOnes {
				needCreateCert = true
			} else {
				return err
			}
		}
		if !needCreateCert {
			for _, cert := range tlsCerts {
				if !cert.IsCA {
					logger.Info("check cert expire for", "domain", cert.DNSNames)
					dateDiff := cert.NotAfter.Sub(time.Now())
					// if not expire ,skip renew certs
					if int64(dateDiff.Hours()/24) > appConfig.CertUpdatePolicy.UpdateBeforeDays {
						logger.Info(
							"the cert info",
							"not after",
							cert.NotAfter,
							"config renew before days",
							appConfig.CertUpdatePolicy.UpdateBeforeDays)
						continue process
					}
				}
			}
		}

		logger.Info("process domain", cdnConfig.Domain)
		provider := certkit.NewChallengeProvider(cdnConfig.Bucket, appConfig.Qiniu.AccessKey, appConfig.Qiniu.SecretKey, cdnConfig.RegionID)
		acme := certkit.NewQiniuACME(user, provider)
		certificate, err := acme.ObtainCertificate(cdnConfig.Domain)
		if err != nil {
			logger.Error("get certificate error", err)
			return err
		}
		cert := certkit.NewCertFrom(certificate)
		uploadCert, err := certMgr.UploadCert(cert)
		if err != nil {
			logger.Error("upload cert error", err)
			return err
		}
		logger.Info("upload cert success", "cert id", uploadCert.CertID)
		logger.Info("start to deploy cert", "target domain", cdnConfig.Domain)
		updateResult, err := certMgr.UpdateHttpsConf(cdnConfig.Domain, uploadCert.CertID, cdnConfig.ForceHttps)
		if err != nil {
			logger.Error("update cert ", "err", err)
			return err
		}
		logger.Info("deploy cert to domain success", "result code", updateResult.Code)
	}
	return nil
}
