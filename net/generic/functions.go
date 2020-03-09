package generic

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-tcp-common/log"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-tcp-client/common"
	"os"
	"strings"
)

var Logger log.Logger = nil

// Apply and Open connection for a specific Connection Handler, based on the Connection HAndler Configuration and the input Parameters / network config files data
func ConnectHandlerViaConfig(connConfig module.ConnectionConfig, handler ConnectionHandler, host defaults.HostValue, netConfig *module.NetProtocolType, depConfig *module.DeployConfig) (NetworkClient, error) {
	var globalError error
	var hostRef string = host.HostName
	if hostRef == "" {
		hostRef = host.IpAddress
	}
	addr := fmt.Sprintf("%s:%s", hostRef, host.Port)
	if connConfig.UseUserPassword {
		globalError = handler.ConnectWithPasswd(addr, netConfig.UserName, netConfig.Password)
	} else if connConfig.UseUserKeyPassphrase {
		var keyFilePath string = netConfig.KeyFile
		if strings.Index(keyFilePath, ":") < 0 &&
			strings.Index(keyFilePath, "/") != 0 &&
			strings.Index(keyFilePath, "\\") != 0 {
			keyFilePath = depConfig.WorkDir + string(os.PathSeparator) + keyFilePath
		}
		globalError = handler.ConnectWithKeyAndPassphrase(addr, netConfig.UserName, keyFilePath, netConfig.Passphrase)
	} else if connConfig.UseUserKey {
		var keyFilePath string = netConfig.KeyFile
		if strings.Index(keyFilePath, ":") < 0 &&
			strings.Index(keyFilePath, "/") != 0 &&
			strings.Index(keyFilePath, "\\") != 0 {
			keyFilePath = depConfig.WorkDir + string(os.PathSeparator) + keyFilePath
		}
		globalError = handler.ConnectWithKey(addr, netConfig.UserName, keyFilePath)
	} else if connConfig.UseTLSCertificates {
		var keyFilePath string = netConfig.KeyFile
		var certFilePath string = netConfig.Certificate
		if strings.Index(keyFilePath, ":") < 0 &&
			strings.Index(keyFilePath, "/") != 0 &&
			strings.Index(keyFilePath, "\\") != 0 {
			keyFilePath = depConfig.WorkDir + string(os.PathSeparator) + keyFilePath
		}
		if strings.Index(certFilePath, ":") < 0 &&
			strings.Index(certFilePath, "/") != 0 &&
			strings.Index(certFilePath, "\\") != 0 {
			certFilePath = depConfig.WorkDir + string(os.PathSeparator) + certFilePath
		}
		var caCertPath string = netConfig.CaCert
		if strings.Index(caCertPath, ":") < 0 &&
			strings.Index(caCertPath, "/") != 0 &&
			strings.Index(caCertPath, "\\") != 0 {
			caCertPath = depConfig.WorkDir + string(os.PathSeparator) + caCertPath
		}
		
		certificate := common.CertificateKeyPair{
			Key:  keyFilePath,
			Cert: certFilePath,
		}
		globalError = handler.ConnectWithCertificate(hostRef, host.Port, certificate, caCertPath)
	} else {
		Logger.Error(string(netConfig.NetProtocol) + ": Unable to determine the Connection mode, between: password, key, key passphrase, certificates")
		return nil, errors.New(string(netConfig.NetProtocol) + ": Unable to determine the connection mode")
	}
	if globalError != nil {
		return nil, globalError
	}
	return handler.GetClient(), nil
}
