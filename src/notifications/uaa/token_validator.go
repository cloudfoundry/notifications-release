package uaa

import (
	"errors"
	"fmt"
	"sync"

	gouaa "github.com/cloudfoundry-community/go-uaa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pivotal-golang/lager"
)

type keysFetcher interface {
	TokenKeys() ([]gouaa.JWK, error)
}

type TokenValidator struct {
	keysFetcher keysFetcher
	keyMap      map[string]gouaa.JWK
	keyMutex    sync.RWMutex
	logger      lager.Logger
}

func NewTokenValidator(logger lager.Logger, keysFetcher keysFetcher) *TokenValidator {
	logger = logger.Session("uaa.token.validator")
	return &TokenValidator{
		logger:      logger,
		keysFetcher: keysFetcher,
		keyMap:      make(map[string]gouaa.JWK),
	}
}

func (v *TokenValidator) LoadSigningKeys() error {
	v.logger.Info("loading.keys")
	keys, err := v.keysFetcher.TokenKeys()

	if err != nil {
		v.logger.Error("loading.keys.failed", err)
		return err
	}

	keyMap := make(map[string]gouaa.JWK, len(keys))
	keyNames := make([]string, 0, len(keys))

	for _, key := range keys {
		keyMap[key.Kid] = key
		keyNames = append(keyNames, key.Kid)
	}

	v.logger.Info("loaded.keys", lager.Data{
		"keys": keyNames,
	})

	v.keyMutex.Lock()
	defer v.keyMutex.Unlock()
	v.keyMap = keyMap

	return nil
}

func (v *TokenValidator) findKey(id string) (gouaa.JWK, bool) {
	v.keyMutex.RLock()
	defer v.keyMutex.RUnlock()

	key, ok := v.keyMap[id]
	return key, ok
}

func (v *TokenValidator) lookUp(keyId string) (string, error) {
	key, ok := v.findKey(keyId)

	if ok {
		return key.Value, nil
	}

	if err := v.LoadSigningKeys(); err != nil {
		return "", err
	}

	key, ok = v.findKey(keyId)

	if !ok {
		return "", fmt.Errorf("unknown key with id %s", keyId)
	}

	return key.Value, nil
}

func (v *TokenValidator) Parse(rawToken string) (*jwt.Token, error) {
	return jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		switch t.Method {
		case
			jwt.SigningMethodRS256,
			jwt.SigningMethodRS384,
			jwt.SigningMethodRS512:
			break
		default:
			return nil, fmt.Errorf("Unsupported signing method %v", t.Method.Alg())
		}

		keyId, ok := t.Header["kid"].(string)

		if !ok {
			return nil, errors.New("Unable to lookup key id for the token")
		}

		key, err := v.lookUp(keyId)

		if err != nil {
			return nil, err
		}

		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
		if err != nil {
			return nil, err
		}

		return pubKey, nil
	})
}
