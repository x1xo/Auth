package utils

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/x1xo/Auth/src/databases"
	"github.com/x1xo/Auth/src/databases/models"
	"go.mongodb.org/mongo-driver/bson"
)

type Certs struct {
	PrivateKey string `json:"private_key" bson:"private_key"`
	PublicKey  string `json:"public_key" bson:"public_key"`
}

var privateKey *rsa.PrivateKey
var PublicKey *string
var PublicJWTKey *jwk.Key

func RandomId(lenght int) (string, error) {
	randomBytes := make([]byte, lenght)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := hex.EncodeToString(randomBytes)
	return randomString, nil
}

func GenerateJWKS() (*jwk.Set, error) {
	if PublicJWTKey == nil {
		if err := LoadCertificates(); err != nil {
			return nil, err
		}
	}

	jwkSet := jwk.NewSet()
	jwkSet.Add(*PublicJWTKey)
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

	jwkKey, err := jwk.ParseKey([]byte(*PublicKey), jwk.WithPEM(true))
	if err != nil {
		return err
	}
	jwk.AssignKeyID(jwkKey)
	jwkKey.Set(jwk.AlgorithmKey, jwa.RS256)
	PublicJWTKey = &jwkKey

	return nil
}

func CreateSesssion(userId, tokenId, ipAddress, userAgent string, expires int) error {
	redis := databases.GetRedis()

	ipInfo, err := GetIPInfo(ipAddress)
	if err != nil {
		return err
	}

	session := models.UserSession{
		UserId:    userId,
		TokenId:   tokenId,
		IPAddress: *ipInfo,
		UserAgent: userAgent,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(expires)),
	}

	jsonSession, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return redis.Set(context.Background(), userId+"_"+tokenId, string(jsonSession), time.Duration(expires)).Err()
}

func GetIPInfo(ipAddress string) (*models.IPAddressInfo, error) {
	infoReq, err := http.Get("https://ipinfo.io/" + ipAddress + "/json")
	if err != nil {
		return nil, err
	}

	defer infoReq.Body.Close()

	body, err := io.ReadAll(infoReq.Body)
	if err != nil {
		return nil, err
	}

	var ipInfo models.IPAddressInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return nil, err
	}

	return &ipInfo, nil
}

func GenerateToken(userId string) (string, string, error) {
	token := jwt.New()

	token.Set("sub", userId)
	token.Set("jti", uuid.New().String())
	token.Set("iat", time.Now().Unix())
	token.Set("exp", time.Now().Add(time.Hour*3).Unix())

	headers := jws.NewHeaders()
	headers.Set("kid", (*PublicJWTKey).KeyID())

	signedToken, err := jwt.Sign(token, jwa.RS256, privateKey, jwt.WithHeaders(headers))
	if err != nil {
		log.Println("[Error] Couldn't sign a token: \n", err)
		return "", "", err
	}

	tokenId, _ := token.Get("jti")
	return string(signedToken), tokenId.(string), nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseString(tokenString, jwt.WithVerify(jwa.RS256, *PublicJWTKey))
	if err != nil {
		return nil, err
	}

	userId := token.Subject()
	if userId == "" {
		return nil, errors.New("invalid token")
	}

	tokenId, _ := token.Get("jti")
	err = databases.GetRedis().Get(context.Background(), userId+"_"+tokenId.(string)).Err()
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return &token, nil
}
