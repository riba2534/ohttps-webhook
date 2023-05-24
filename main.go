package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Payload struct {
	CertificateName           string   `json:"certificateName"`
	CertificateDomains        []string `json:"certificateDomains"`
	CertificateCertKey        string   `json:"certificateCertKey"`
	CertificateFullchainCerts string   `json:"certificateFullchainCerts"`
	CertificateExpireAt       int64    `json:"certificateExpireAt"`
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
	tlsPath := os.Getenv("TLS_PATH")
	callbackToken := os.Getenv("CALLBACK_TOKEN")

	if tlsPath == "" || callbackToken == "" {
		panic("tlsPath or callbackToken not set")
	}
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world!")
	})

	r.POST("/webhook", func(c *gin.Context) {
		var body RequestBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("body is ", MarshalAny(body))
		expectedSign := md5.Sum([]byte(strconv.FormatInt(body.Timestamp, 10) + ":" + callbackToken))
		if body.Sign != strings.ToLower(hex.EncodeToString(expectedSign[:])) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sign"})
			return
		}
		err := ioutil.WriteFile(tlsPath+"fullchain.cer", []byte(body.Payload.CertificateFullchainCerts), 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to write fullchain certificate"})
			return
		}
		err = ioutil.WriteFile(tlsPath+"cert.key", []byte(body.Payload.CertificateCertKey), 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to write private key"})
			return
		}
		c.JSON(http.StatusOK, ResponseBody{Success: true})
	})
	r.Run(":4321")
}

func MarshalAny(i interface{}) string {
	s, _ := json.Marshal(i)
	return string(s)
}
