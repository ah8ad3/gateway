package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ah8ad3/gateway/pkg/logger"
)

// SecretKey for security issue
var SecretKey string

var proxyDir, pluginDir, secretDir, dbDir string

func init() {
	test := os.Getenv("TEST")
	if test == "1" {
		proxyDir = "./../../db/proxy.bin"
		pluginDir = "./../../db/plugin.bin"
		secretDir = "./../../db/secret.bin"
		dbDir = "./../../db/"
	} else {
		proxyDir = "db/proxy.bin"
		pluginDir = "db/plugin.bin"
		secretDir = "db/secret.bin"
		dbDir = "db/"
	}
}

// InsertProxy func for insert all proxy that been Marshal and encrypt and save it to db/proxy.bin file
func InsertProxy(proxies []byte) {
	if _, err := ioutil.ReadFile(proxyDir); err != nil {
		proxies = encryptData(proxies)
		err := ioutil.WriteFile(proxyDir, proxies, 0644)
		if err != nil {
			logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "critical",
				Description: err.Error()}})
		}
	} else {
		if err = os.Remove(proxyDir); err != nil {
			logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "critical",
				Description: err.Error()}})
			log.Fatal("Cant remove proxy file")
		}
		logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "log",
			Description: "File proxies.bin removed"}})
		InsertProxy(proxies)
	}
}

// InsertPlugins func for insert all plugins that been Marshal and encrypt and save it to db/plugin.bin file
func InsertPlugins(plugins []byte) {
	if _, err := ioutil.ReadFile(pluginDir); err != nil {
		plugins = encryptData(plugins)
		err := ioutil.WriteFile(pluginDir, plugins, 0644)
		if err != nil {
			logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "critical",
				Description: err.Error()}})
		}
	} else {
		if err = os.Remove(pluginDir); err != nil {
			logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "critical",
				Description: err.Error()}})
			log.Fatal("Cant remove plugin file")
		}
		logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "log",
			Description: "File plugin.bin removed"}})
		InsertPlugins(plugins)
	}
}

// GetProxies to decrypt and get all saved proxy as bson
func GetProxies() []byte {
	data, err := ioutil.ReadFile(proxyDir)
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "critical",
			Description: err.Error()}})
		log.Fatal("proxy.bin not found")
	}
	data = decryptData(data)

	return data
}

// GetPlugins to decrypt and get all saved plugins as bson
func GetPlugins() []byte {
	data, err := ioutil.ReadFile(pluginDir)
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "critical",
			Description: err.Error()}})
		log.Fatal("plugin.bin not found")
	}
	data = decryptData(data)

	return data
}

func createHash(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))
	return hex.EncodeToString(hash.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Description: err.Error(), Event: "critical"},
			Time: time.Now(), Pkg: "db"})
		log.Fatal("Error in Data Encryption")
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Description: err.Error(), Event: "critical"},
			Time: time.Now(), Pkg: "db"})
		log.Fatal("Error in Data Encryption")
	}
	cipherText := gcm.Seal(nonce, nonce, data, nil)
	return cipherText
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Description: err.Error(), Event: "critical"},
			Time: time.Now(), Pkg: "db"})
		log.Fatal("Error in Data Decryption")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Description: err.Error(), Event: "critical"},
			Time: time.Now(), Pkg: "db"})
		log.Fatal("Error in Data Decryption")
	}
	nonceSize := gcm.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Description: err.Error(), Event: "critical"},
			Time: time.Now(), Pkg: "db"})
		log.Fatal("Error in Data Decryption")
	}
	return plaintext
}

func encryptData(data []byte) []byte {
	pass := SecretKey
	if pass == "" {
		log.Fatal("Secret Key can not be empty, for security issue")
	}

	return encrypt(data, pass)
}

func decryptData(data []byte) []byte {
	pass := SecretKey
	if pass == "" {
		log.Fatal("Secret Key can not be empty, for security issue")
	}

	return decrypt(data, pass)
}

func GenerateSecretKey() {
	key := make([]byte, 16)

	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err.Error())
	}

	sec := fmt.Sprintf("%s", key)

	isSaved := saveSecretKey(sec)
	if isSaved {
		SecretKey = sec
	} else {
		str, found := LoadSecretKey()
		if found {
			SecretKey = str
		} else {
			GenerateSecretKey()
		}
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// saveSecretKey save it to file
func saveSecretKey(secret string) bool {
	exist, _ := exists(dbDir)
	if exist {
		if _, err := ioutil.ReadFile(secretDir); err != nil {
			err := ioutil.WriteFile(secretDir, []byte(secret), 0644)
			if err != nil {
				log.Fatal(err.Error())
			}
			return true
		}
		return false
	} else {
		_ = os.Mkdir(dbDir, 0755)
		return saveSecretKey(secret)
	}
}

// LoadSecretKey load it to SecretKey
func LoadSecretKey() (string, bool) {
	data, err := ioutil.ReadFile(secretDir)
	if err != nil {
		return "", false
	}
	return fmt.Sprintf("%s", data), true
}
