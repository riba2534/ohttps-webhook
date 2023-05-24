package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Payload struct {
	CertificateName           string `json:"certificateName"`
	CertificateDomains        string `json:"certificateDomains"`
	CertificateCertKey        string `json:"certificateCertKey"`
	CertificateFullchainCerts string `json:"certificateFullchainCerts"`
	CertificateExpireAt       string `json:"certificateExpireAt"`
}

type RequestBody struct {
	Timestamp int64   `json:"timestamp"`
	Payload   Payload `json:"payload"`
	Sign      string  `json:"sign"`
}

type ResponseBody struct {
	Success bool `json:"success"`
}

func main() {
	certKeyPath := os.Getenv("CERT_KEY_PATH")
	fullChainPath := os.Getenv("FULL_CHAIN_PATH")
	callbackToken := os.Getenv("CALLBACK_TOKEN")

	if certKeyPath == "" || fullChainPath == "" || callbackToken == "" {
		panic("CERT_KEY_PATH, FULL_CHAIN_PATH or CALLBACK_TOKEN not set")
	}

	r := gin.Default()

	r.POST("/webhook", func(c *gin.Context) {
		var body RequestBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		expectedSign := md5.Sum([]byte(strconv.FormatInt(body.Timestamp, 10) + ":" + callbackToken))
		if body.Sign != strings.ToLower(hex.EncodeToString(expectedSign[:])) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sign"})
			return
		}

		err := ioutil.WriteFile(fullChainPath, []byte(body.Payload.CertificateFullchainCerts), 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to write fullchain certificate"})
			return
		}

		err = ioutil.WriteFile(certKeyPath, []byte(body.Payload.CertificateCertKey), 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to write private key"})
			return
		}

		c.JSON(http.StatusOK, ResponseBody{Success: true})
	})

	r.Run()
}
