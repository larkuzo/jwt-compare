package main

import (
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"crypto/rsa"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
	"time"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

func main() {
	rsaPrivateKey, _ := rsa.GenerateKey(rand.Reader, 512)
	ecdsaPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	hmacKey := []byte("secret")

	// First test
	jwtRSA(rsaPrivateKey, map[string]interface{}{"data": "this is a signed token"})
	jwtECDSA(ecdsaPrivateKey, map[string]interface{}{"data": "this is a signed token"})
	jwtHMAC(hmacKey, map[string]interface{}{"data": "this is a signed token"})

	http.HandleFunc("/rsa", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		token := jwtRSA(rsaPrivateKey, map[string]interface{}{"data": "this is a signed token"})
		end := time.Since(start)
		data, _ := json.Marshal(token)
		w.Write(data)
		fmt.Println(end.Nanoseconds())
		fmt.Println(len(token))
	})

	http.HandleFunc("/ecdsa", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		token := jwtECDSA(ecdsaPrivateKey, map[string]interface{}{"data": "this is a signed token"})
		end := time.Since(start)
		data, _ := json.Marshal(token)
		w.Write(data)
		fmt.Println(end.Nanoseconds())
		fmt.Println(len(token))
	})

	http.HandleFunc("/hmac", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		token := jwtHMAC(hmacKey, map[string]interface{}{"data": "this is a signed token"})
		end := time.Since(start)
		data, _ := json.Marshal(token)
		w.Write(data)
		fmt.Println(end.Nanoseconds())
		fmt.Println(len(token))
	})

	go http.ListenAndServe(":8080", nil)

	var response *http.Response
	for i := 0; i < 50; i++ {
		response, _ = http.Get("http://localhost:8080/rsa")
		ioutil.ReadAll(response.Body)
		response, _ = http.Get("http://localhost:8080/ecdsa")
		ioutil.ReadAll(response.Body)
		response, _ = http.Get("http://localhost:8080/hmac")
		ioutil.ReadAll(response.Body)
	}
}

func testJWT() {
	//op := 50
	//testRSA(rsaPrivateKey, op)
	//testECDSA(ecdsaPrivateKey, op)
	//testHMAC(hmacKey, op)
}

func testRSA(key *rsa.PrivateKey, op int) {
	plot := make([]int64, op)
	for i := 0; i < op; i++ {
		start := time.Now()
		jwtRSA(key, map[string]interface{}{"data": "this is a signed token"})
		end := time.Since(start)
		plot[i] = end.Nanoseconds()
	}
	fmt.Println("RSA")
	for i := 0; i < op; i++ {
		fmt.Println(plot[i])
	}
	fmt.Println()
}

func testECDSA(key *ecdsa.PrivateKey, op int) {
	plot := make([]int64, op)
	for i := 0; i < op; i++ {
		start := time.Now()
		jwtECDSA(key, map[string]interface{}{"data": "this is a signed token"})
		end := time.Since(start)
		plot[i] = end.Nanoseconds()
	}

	fmt.Println("ECDSA")
	for i := 0; i < op; i++ {
		fmt.Println(plot[i])
	}
	fmt.Println()
}

func testHMAC(key []byte, op int) {
	plot := make([]int64, op)
	for i := 0; i < op; i++ {
		start := time.Now()
		jwtHMAC(key, map[string]interface{}{"data": "this is a signed token"})
		end := time.Since(start)
		plot[i] = end.Nanoseconds()
	}

	fmt.Println("HMAC")
	for i := 0; i < op; i++ {
		fmt.Println(plot[i])
	}
	fmt.Println()
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
