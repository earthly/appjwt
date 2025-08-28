package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// loadPrivateKey loads an RSA private key from a PEM file.
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading private key: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	return privateKey, nil
}

// generateJWT generates a signed JWT for the given GitHub App ID.
func generateJWT(appID string, key *rsa.PrivateKey) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iss": appID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signed, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("error signing JWT: %w", err)
	}

	return signed, nil
}

func main() {
	// Define CLI flags
	var (
		appID   = flag.String("app-id", "", "GitHub App ID (required)")
		keyPath = flag.String("key-file", "", "Path to GitHub App private key PEM file (required)")
	)

	flag.Parse()

	// Validate required flags
	if *appID == "" || *keyPath == "" {
		flag.Usage()
		log.Fatal("Both --app-id and --key-file are required")
	}

	// Load private key
	privateKey, err := loadPrivateKey(*keyPath)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Generate JWT
	jwtToken, err := generateJWT(*appID, privateKey)
	if err != nil {
		log.Fatalf("Failed to generate JWT: %v", err)
	}

	// Print result
	fmt.Println(jwtToken)
}
