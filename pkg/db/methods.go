package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"github.com/ah8ad3/gateway/pkg/logger"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// InsertProxy func for insert all proxy that been Marshal and encrypt and save it to db/proxy.bin file
func InsertProxy(proxies []byte) {
	if _, err := ioutil.ReadFile("db/proxy.bin"); err != nil {
		proxies = encryptData(proxies)
		err := ioutil.WriteFile("db/proxy.bin", proxies, 0644)
		if err != nil {
			logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "critical",
				Description: err.Error()}})
		}
	}else {
		if err = os.Remove("db/proxy.bin"); err != nil {
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
	if _, err := ioutil.ReadFile("db/plugin.bin"); err != nil {
		plugins = encryptData(plugins)
		err := ioutil.WriteFile("db/plugin.bin", plugins, 0644)
		if err != nil {
			logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "critical",
				Description: err.Error()}})
		}
	}else {
		if err = os.Remove("db/plugin.bin"); err != nil {
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
	data, err := ioutil.ReadFile("db/proxy.bin")
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
	data, err := ioutil.ReadFile("db/plugin.bin")
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
	pass := os.Getenv("SECRET_KEY")
	if pass == "" {
		log.Fatal("Secret Key can not be empty, for security issue")
	}

	return encrypt(data, pass)
}

func decryptData(data []byte) []byte {
	pass := os.Getenv("SECRET_KEY")
	if pass == "" {
		log.Fatal("Secret Key can not be empty, for security issue")
	}

	return decrypt(data, pass)
}
