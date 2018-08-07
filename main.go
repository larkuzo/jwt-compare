package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Result struct {
	time     int64
	size     int
	transfer int64
}

func result(time int64, size int, transfer int64) Result {
	return Result{time, size, transfer}
}

const TestCount = 50

func main() {
	rsaPrivateKey, _ := rsa.GenerateKey(rand.Reader, 512)
	ecdsaPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	hmacKey := make([]byte, 24)
	rand.Read(hmacKey)

	// First test JWT
	payload := map[string]interface{}{"data": "this is a signed token"}
	jwtRSA(rsaPrivateKey, payload)
	jwtECDSA(ecdsaPrivateKey, payload)
	jwtHMAC(hmacKey, payload)

	// Preparation
	num := TestCount
	rsaResult := make([]Result, num)
	ecdsaResult := make([]Result, num)
	hmacResult := make([]Result, num)

	http.HandleFunc("/rsa", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		token := jwtRSA(rsaPrivateKey, map[string]interface{}{"data": "this is a signed token"})
		end := time.Since(start)
		data, _ := json.Marshal(token)
		w.Write(data)
		rsaResult[TestCount-num] = result(end.Nanoseconds(), len(token), 0)
	})

	http.HandleFunc("/ecdsa", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		token := jwtECDSA(ecdsaPrivateKey, map[string]interface{}{"data": "this is a signed token"})
		end := time.Since(start)
		data, _ := json.Marshal(token)
		w.Write(data)
		ecdsaResult[TestCount-num] = result(end.Nanoseconds(), len(token), 0)
	})

	http.HandleFunc("/hmac", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		token := jwtHMAC(hmacKey, map[string]interface{}{"data": "this is a signed token"})
		end := time.Since(start)
		data, _ := json.Marshal(token)
		w.Write(data)
		hmacResult[TestCount-num] = result(end.Nanoseconds(), len(token), 0)
	})

	// Start server
	if port, exists := os.LookupEnv("PORT"); exists {
		go http.ListenAndServe(":"+port, nil)
	} else {
		go http.ListenAndServe(":8080", nil)
	}

	var response *http.Response
	var start time.Time

	// First test Request
	response, _ = http.Get("http://localhost:8080/rsa")
	ioutil.ReadAll(response.Body)
	response, _ = http.Get("http://localhost:8080/ecdsa")
	ioutil.ReadAll(response.Body)
	response, _ = http.Get("http://localhost:8080/hmac")
	ioutil.ReadAll(response.Body)

	// Server test and measure
	for num > 0 {
		start = time.Now()
		response, _ = http.Get("http://localhost:8080/rsa")
		ioutil.ReadAll(response.Body)
		rsaResult[TestCount-num].transfer = time.Since(start).Nanoseconds()

		start = time.Now()
		response, _ = http.Get("http://localhost:8080/ecdsa")
		ioutil.ReadAll(response.Body)
		ecdsaResult[TestCount-num].transfer = time.Since(start).Nanoseconds()

		start = time.Now()
		response, _ = http.Get("http://localhost:8080/hmac")
		ioutil.ReadAll(response.Body)
		hmacResult[TestCount-num].transfer = time.Since(start).Nanoseconds()

		num -= 1
	}

	// Write result
	fmt.Println("RSA")
	fmt.Println("generate (ns),length,transfer (ns)")
	for i := 0; i < TestCount; i++ {
		fmt.Print(rsaResult[i].time, ",", rsaResult[i].size, ",", rsaResult[i].transfer)
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("ECDSA")
	fmt.Println("generate (ns),length,transfer (ns)")
	for i := 0; i < TestCount; i++ {
		fmt.Print(ecdsaResult[i].time, ",", ecdsaResult[i].size, ",", ecdsaResult[i].transfer)
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("HMAC")
	fmt.Println("generate (ns),length,transfer (ns)")
	for i := 0; i < TestCount; i++ {
		fmt.Print(hmacResult[i].time, ",", hmacResult[i].size, ",", hmacResult[i].transfer)
		fmt.Println()
	}
}

func jwtRSA(privateKey *rsa.PrivateKey, payload interface{}) string {
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.RS256,
		Key:       privateKey,
	}, nil)

	if err != nil {
		panic(err)
	}

	return jwtToken(signer, payload)
}

func jwtECDSA(privateKey *ecdsa.PrivateKey, payload interface{}) string {
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.ES256,
		Key:       privateKey,
	}, nil)

	if err != nil {
		panic(err)
	}

	return jwtToken(signer, payload)
}

func jwtHMAC(privateKey []byte, payload interface{}) string {
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       privateKey,
	}, nil)

	if err != nil {
		panic(err)
	}

	return jwtToken(signer, payload)
}

func jwtToken(signer jose.Signer, payload interface{}) string {
	builder := jwt.Signed(signer).Claims(payload)
	token, err := builder.CompactSerialize()
	if err != nil {
		panic(err)
	}

	return token
}
