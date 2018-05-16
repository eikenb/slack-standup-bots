package main

// How to store a secret key/token in your repo in an encrypted format.

// Run this command to encrypt your "secret"
//
// aws kms encrypt --key-id alias/key --plaintext "secret" \
//    --query CiphertextBlob --output text
//
// Take the output and save that in your repo as the key.
// At runtime decrypt the key to use it as seen below.

// In this case the base64 encoded, encrypted token is stored in an environment
// variable. This grabs it and decrypts it for use.

// It is done using a chan like this to have it work like a future.
// Run this ASAP and hopefully the decrypted token will be ready when needed.
// Run using `go decryptToken()`

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

var slackToken chan string

func init() {
	slackToken = make(chan string)
}

func decryptToken() {
	sess, err := session.NewSession(
		&aws.Config{Region: aws.String("us-west-2")})
	if err != nil {
		fmt.Println("NewSession Error:", err)
		os.Exit(1)
	}
	val, ok := os.LookupEnv("kmsEncryptedToken")
	if !ok {
		fmt.Println("Missing kmsEncryptedToken environment variable.")
		os.Exit(1)
	}
	blob, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		fmt.Println("Base64 Decode Error:", err)
		os.Exit(1)
	}

	svc := kms.New(sess)
	result, err := svc.Decrypt(&kms.DecryptInput{CiphertextBlob: blob})
	if err != nil {
		fmt.Println("Decrypt Error:", err)
		os.Exit(1)
	}

	token := string(result.Plaintext)
	for {
		slackToken <- token
	}
}
