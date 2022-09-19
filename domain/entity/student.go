package entity

import (
	"context"
	"keyvault-demo/internal/keyvault"
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID         uuid.UUID `gorm:"primaryKey,autoIncrement" json:"id"`
	NIM        string    `gorm:"nim" json:"nim"`
	Name       string    `gorm:"name" json:"name"`
	NIK        string    `gorm:"nik" json:"nik"`
	Phone      string    `gorm:"phone" json:"phone"`
	KeyVersion string    `gorm:"key_version" json:"-"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

var keyName = "demo"

func (s *Student) GetAndSetID() *Student {
	s.ID = uuid.New()
	return s
}

func (s *Student) EncryptNIK() (err error) {
	encrypted, err := keyvault.KeyVault().Encrypt(context.TODO(), s.NIK, keyName, s.KeyVersion)
	if err != nil {
		return err
	}
	s.NIK = encrypted.Result
	s.KeyVersion = encrypted.KeyVersion
	return nil
}

func (s *Student) EncryptPhone() (err error) {
	encrypted, err := keyvault.KeyVault().Encrypt(context.TODO(), s.Phone, keyName, s.KeyVersion)
	if err != nil {
		return err
	}
	s.Phone = encrypted.Result
	s.KeyVersion = encrypted.KeyVersion
	return nil
}

func (s *Student) DecryptNIK() (err error) {
	decrypted, err := keyvault.KeyVault().Decrypt(context.TODO(), s.NIK, keyName, s.KeyVersion)
	if err != nil {
		return err
	}
	s.NIK = decrypted
	return nil
}

func (s *Student) DecryptPhone() (err error) {
	decrypted, err := keyvault.KeyVault().Decrypt(context.TODO(), s.Phone, keyName, s.KeyVersion)
	if err != nil {
		return err
	}
	s.Phone = decrypted
	return nil
}
