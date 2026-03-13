package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"llm-ops/api/grpc"
	apiV1 "llm-ops/api/nexus/api/nexus-api-server/v1"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// @brief:填充明文
func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

// @brief:去除填充数据
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// @brief:AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// @brief:AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// 配置结构体
type Configuration struct {
	AESKey string `yaml:"aes_key"`

	MySQL struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`

	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Kubernetes struct {
		ConfigPath   string `yaml:"config_path"`
		ConfigPathV2 string `yaml:"config_path_v2"`
	} `yaml:"kubernetes"`

	alpha struct {
		ApiURL string `yaml:"api_url"`
		ApiKey string `yaml:"api_key"`
	} `yaml:"alpha"`

	beta struct {
		ApiURL string `yaml:"api_url"`
		ApiKey string `yaml:"api_key"`
	} `yaml:"beta"`

	Tracker TrackerConfig `yaml:"tracker"`

	// Fusion 代理配置
	betaFusion struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"beta_fusion"`

	AlphaFusion struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"alpha_fusion"`
}

type TrackerConfig struct {
	BaseURL string `yaml:"base_url"`
	Token   string `yaml:"token"`
}

var Config Configuration

type NexusCluster struct {
	ID         string
	GrpcURL    string
	Token      string
	Index      int
	IndexAlias int
	Client     apiV1.NexusApiClient `json:"-"`
}

var NexusClusters = map[string]*NexusCluster{
	"us-ca-02": &NexusCluster{
		Index:      3,
		IndexAlias: 7,
		ID:         "us-ca-02",
		GrpcURL:    "grpc-nexus.us-ca-2.gpu-instance.alpha.ai:443",
		Token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJuYW1lIjoiamluZ3h1YW4ifQ.Y8wRDeNDI7hPMaH4BCMtUNkhAz7j9eY0tOURLTVQR9w",
	},
	"us-ca-03": &NexusCluster{
		Index:      5,
		IndexAlias: 8,
		ID:         "us-ca-03",
		GrpcURL:    "grpc-nexus.us-ca-3.gpu-instance.alpha.ai:443",
		Token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoic2VydmVybGVzc19wcm9kIiwibmFtZSI6Imppbmd4dWFuIn0.MQlR-SZwkZuz_iO9Gg34xgq5sp5bEQ5rl_RSgirCeVI",
	},
	"us-01": &NexusCluster{
		Index:   1,
		ID:      "us-01",
		GrpcURL: "grpc-nexus.us-south-1.gpu-instance.alpha.ai:443",
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJuYW1lIjoiamluZ3h1YW4ifQ.pgt1QJNZN_SKniIb9SWE35TEOToXdrnrzCckslLkOPA",
	},
	"us-ca-01": &NexusCluster{
		Index:   2,
		ID:      "us-ca-01",
		GrpcURL: "grpc-nexus.us-ca-1.gpu-instance.alpha.ai:443",
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJuYW1lIjoiamluZ3h1YW4ifQ.8gPySf2Cyon94AvQZJevJaMn9RDNN59FWqJtZUza-xc",
	},
	"us-nyc-01": &NexusCluster{
		Index:      4,
		IndexAlias: 9,
		ID:         "us-nyc-01",
		GrpcURL:    "grpc-nexus.us-nyc-01.gpu-instance.alpha.ai:443",
		Token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJuYW1lIjoiamluZ3h1YW4ifQ.Fy5xAUI-rYXEB2rRdtBvVdLcra8m13uoHTS_9yEX0Is",
	},
	"cn-northwest-2": &NexusCluster{
		Index:   26,
		ID:      "cn-northwest-2",
		GrpcURL: "grpc-nexus.cn-northwest-2.gpu-instance.ppinfra.com:443",
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoic2VydmVybGVzcyIsIm5hbWUiOiJzZXJ2ZXJsZXNzIn0.Sa1tFJBvsAfV1Bx0VHfe1uOVCuajVEZ66V2KkDX23FY",
	},
	"cn-northwest-1": &NexusCluster{
		Index:   27,
		ID:      "cn-northwest-1",
		GrpcURL: "grpc-nexus.cn-northwest-1.gpu-instance.ppinfra.com:443",
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoic2VydmVybGVzc19wcm9kIiwibmFtZSI6InNlcnZlcmxlc3MifQ.l8598ZBm6MSsFfdWacyIs4RI8gAjPz_Nyuuw8lxTctE",
	},
}

func init() {
	log.Println("initing NexusClusters")
	for _, cluster := range NexusClusters {
		conn := grpc.NewApiServerClient(&grpc.RemoteServerOptions{
			Address:  cluster.GrpcURL,
			Token:    cluster.Token,
			Insecure: false,
		})
		client, err := conn.Client()
		if err != nil {
			log.Println(fmt.Sprintf("failed to create client: %v, url: %s, token: %s", err, cluster.GrpcURL, cluster.Token))
			return
		}
		if err != nil {
			log.Println(fmt.Sprintf("failed to create client: %v, url: %s, token: %s", err, cluster.GrpcURL, cluster.Token))
			continue
		}
		cluster.Client = client
	}
}

func LoadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading config file %s: %v", configPath, err)
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	if Config.AESKey == "" {
		return nil
	}

	var aeskey = []byte(Config.AESKey)

	mysqlPassword := Config.MySQL.Password
	bytesPass, err := base64.StdEncoding.DecodeString(mysqlPassword)
	if err != nil {
		return err
	}
	tpass, err := AesDecrypt(bytesPass, aeskey)
	if err != nil {
		return err
	}
	Config.MySQL.Password = string(tpass)
	log.Println(fmt.Sprintf("mysqlPassword: %s", Config.MySQL.Password))

	trackerToken := Config.Tracker.Token
	bytesPass, err = base64.StdEncoding.DecodeString(trackerToken)
	if err != nil {
		return err
	}
	tpass, err = AesDecrypt(bytesPass, aeskey)
	if err != nil {
		return err
	}
	Config.Tracker.Token = string(tpass)
	log.Println(fmt.Sprintf("trackerToken: %s", Config.Tracker.Token))

	alphaApiKey := Config.Alpha.ApiKey
	bytesPass, err = base64.StdEncoding.DecodeString(alphaApiKey)
	if err != nil {
		return err
	}
	tpass, err = AesDecrypt(bytesPass, aeskey)
	if err != nil {
		return err
	}
	Config.Alpha.ApiKey = string(tpass)
	log.Println(fmt.Sprintf("alphaApiKey: %s", Config.Alpha.ApiKey))

	betaApiKey := Config.Beta.ApiKey
	bytesPass, err = base64.StdEncoding.DecodeString(betaApiKey)
	if err != nil {
		return err
	}
	tpass, err = AesDecrypt(bytesPass, aeskey)
	if err != nil {
		return err
	}
	Config.Beta.ApiKey = string(tpass)
	log.Println(fmt.Sprintf("betaApiKey: %s", Config.Beta.ApiKey))

	// 解密 beta Fusion Token
	if Config.betaFusion.Token != "" {
		betaFusionToken := Config.betaFusion.Token
		bytesPass, err = base64.StdEncoding.DecodeString(betaFusionToken)
		if err != nil {
			return err
		}
		tpass, err = AesDecrypt(bytesPass, aeskey)
		if err != nil {
			return err
		}
		Config.betaFusion.Token = string(tpass)
		log.Println(fmt.Sprintf("betaFusionToken: %s", Config.betaFusion.Token))
	}

	// 解密 alpha Fusion Token
	if Config.AlphaFusion.Token != "" {
		AlphaFusionToken := Config.AlphaFusion.Token
		bytesPass, err = base64.StdEncoding.DecodeString(AlphaFusionToken)
		if err != nil {
			return err
		}
		tpass, err = AesDecrypt(bytesPass, aeskey)
		if err != nil {
			return err
		}
		Config.AlphaFusion.Token = string(tpass)
		log.Println(fmt.Sprintf("AlphaFusionToken: %s", Config.AlphaFusion.Token))
	}

	return nil
}

func GetClusterByID(id string) *NexusCluster {
	for _, cluster := range NexusClusters {
		if cluster.ID == id {
			return cluster
		}
	}
	return nil
}

func GetClusterByIndex(index int) *NexusCluster {
	for _, cluster := range NexusClusters {
		if cluster.Index == index || cluster.IndexAlias == index {
			return cluster
		}
	}
	return nil
}

func GetClusterNames() []string {
	names := make([]string, 0, len(NexusClusters))
	for _, nCluster := range NexusClusters {
		names = append(names, nCluster.ID)
	}
	return names
}
