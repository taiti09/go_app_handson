package auth

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/taiti09/go_app_handson/clock"
	"github.com/taiti09/go_app_handson/entity"
)

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	Store Store
	Clocker clock.Clocker
}

const (
	RoleKey = "role"
	UserNameKey = "user_name"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

func NewJWTer(s Store, c clock.Clocker) (*JWTer, error) {
	j := &JWTer{Store: s}
	privKey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: private key: %w", err)
	}
	pubkey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: public key: %w", err)
	}
	j.PrivateKey = privKey
	j.PublicKey = pubkey
	j.Clocker = c
	return j, nil
}

func parse(rawKey []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(rawKey,jwk.WithPEM(true))
	if err != nil {
		return nil,err
	}
	return key, nil
}

func (j *JWTer) GenerateToken(ctx context.Context, u entity.User) ([]byte, error) {
	tok, err := jwt.NewBuilder().JwtID(uuid.New().String()).Issuer(`github.com/taiti09/go_app_handson`).Subject("access_token").IssuedAt(j.Clocker.Now()).Expiration(j.Clocker.Now().Add(30*time.Minute)).Claim(RoleKey,u.Role).Claim(UserNameKey,u.Name).Build()
	if err != nil {
		return nil, fmt.Errorf("GetToken failed to build token: %w", err)
	}
	if err := j.Store.Save(ctx,tok.JwtID(),u.ID); err != nil {
		return nil, err
	}

	signed, err := jwt.Sign(tok,jwt.WithKey(jwa.RS256,j.PrivateKey))
	if err != nil {
		return nil, err
	}
	return signed, nil
}