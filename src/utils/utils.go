package utils

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/x1xo/Auth/src/databases"
	"go.mongodb.org/mongo-driver/bson"
)

type Certs struct {
	PrivateKey string `json:"private_key" bson:"private_key"`
	PublicKey  string `json:"public_key" bson:"public_key"`
}

var privateKey *rsa.PrivateKey
var PublicKey *string

func GenerateJWKS() (*jwk.Set, error) {
	if PublicKey == nil {
		if err := LoadCertificates(); err != nil {
			return nil, err
		}
	}
	jwkKey, err := jwk.ParseKey([]byte(*PublicKey), jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}

	jwkSet := jwk.NewSet()
	jwkSet.Add(jwkKey)
	return &jwkSet, nil
}

func LoadCertificates() error {
	var certs *Certs
	err := databases.GetMongoDatabase().Collection("certs").FindOne(context.Background(), bson.M{}).Decode(&certs)
	if err != nil {
		return err
	}

	block, _ := pem.Decode([]byte(certs.PrivateKey))
	if block == nil {
		return err
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	privateKey = privKey.(*rsa.PrivateKey)

	PublicKey = &certs.PublicKey
	return nil
}

func RandomId(lenght int) (string, error) {
	randomBytes := make([]byte, lenght)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := hex.EncodeToString(randomBytes)
	return randomString, nil
}

func GenerateToken(userId string) (string, string, error) {
	token := jwt.New()

	token.Set("sub", userId)
	token.Set("tokenId", uuid.New().String())
	token.Set("iat", time.Now().Unix())
	token.Set("exp", time.Now().Add(time.Hour*3).Unix())

	signedToken, err := jwt.Sign(token, jwa.RS256, privateKey)
	if err != nil {
		log.Println("[Error] Couldn't sign a token: \n", err)
		return "", "", err
	}
	tokenId, _ := token.Get("tokenId")
	return string(signedToken), tokenId.(string), nil
}
