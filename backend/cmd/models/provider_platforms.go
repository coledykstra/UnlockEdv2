package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
)

type ProviderPlatformType string

const (
	CanvasOSS   ProviderPlatformType = "canvas_oss"
	CanvasCloud ProviderPlatformType = "canvas_cloud"
	Kolibri     ProviderPlatformType = "kolibri"
)

type ProviderPlatformState string

const (
	Enabled  ProviderPlatformState = "enabled"
	Disabled ProviderPlatformState = "disabled"
	Archived ProviderPlatformState = "archived"
)

type ProviderPlatform struct {
	DatabaseFields
	Type                   ProviderPlatformType  `gorm:"size:100" json:"type"`
	Name                   string                `gorm:"size:255" json:"name"`
	Description            string                `gorm:"size:1024" json:"description"`
	IconUrl                string                `gorm:"size:255" json:"icon_url"`
	AccountID              string                `gorm:"size:64" json:"account_id"`
	AccessKey              string                `gorm:"size:255" json:"access_key"`
	BaseUrl                string                `gorm:"size:255" json:"base_url"`
	State                  ProviderPlatformState `gorm:"size:100" json:"state"`
	ExternalAuthProviderId string                `gorm:"size:128" json:"external_auth_provider_id"`

	Programs             []Program             `json:"-"`
	ProviderUserMappings []ProviderUserMapping `json:"-"`
}

func (ProviderPlatform) TableName() string {
	return "provider_platforms"
}

func (provider *ProviderPlatform) DecryptAccessKey() (string, error) {
	key := os.Getenv("APP_KEY")
	hashedKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashedKey[:])
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(provider.AccessKey)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

func (provider *ProviderPlatform) EncryptAccessKey() (string, error) {
	key := os.Getenv("APP_KEY")
	hashedKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashedKey[:])
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(provider.AccessKey))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(provider.AccessKey))
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}
