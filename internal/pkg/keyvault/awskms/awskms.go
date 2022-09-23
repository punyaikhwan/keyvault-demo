package awskms

import (
	"context"
	"encoding/base64"
	"fmt"
	"keyvault-demo/config"
	"keyvault-demo/internal/pkg/keyvault"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type awsKeyVault struct {
	vaultURL string
	client   *kms.KMS
}

func (az *awsKeyVault) getClient() *kms.KMS {
	if az.client == nil {
		fmt.Println("Creating AWS client...")
		cred := credentials.NewStaticCredentials(config.Configuration().AWSAccessKeyID, config.Configuration().AWSSecretKey, "")
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String("ap-southeast-1"),
			Credentials: cred,
		})
		if err != nil {
			log.Fatalf("failed to obtain a credential: %v", err)
		}

		az.client = kms.New(sess)
	}

	return az.client
}

func NewAWSKeyVault(vaultURL string) keyvault.KeyVault {
	return &awsKeyVault{vaultURL, nil}
}

func (az *awsKeyVault) Encrypt(ctx context.Context, plaintext string, keyName string, keyVersion string) (result keyvault.EncryptionResult, err error) {
	encryptRes, err := az.getClient().Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(keyName),
		Plaintext: []byte(plaintext),
	})

	if err != nil {
		return result, err
	}

	result.Result = base64.StdEncoding.EncodeToString(encryptRes.CiphertextBlob)
	result.KeyVersion = keyVersion
	return result, nil
}

func (az *awsKeyVault) Decrypt(ctx context.Context, encryptedBase64 string, keyName string, keyVersion string) (decryptedBase64 string, err error) {
	encryptedByte, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return decryptedBase64, err
	}
	decryptedRes, err := az.getClient().Decrypt(&kms.DecryptInput{
		CiphertextBlob: encryptedByte,
	})

	if err != nil {
		return "", err
	}

	return string(decryptedRes.Plaintext), nil
}

// func (az *awsKeyVault) getLatestKeyVersion(ctx context.Context, keyName string) (keyVersion string, err error) {
// 	key, err := az.getClient().ListAliasesWithContext(ctx, &kms.ListAliasesInput{
// 		KeyID:
// 	})
// 	if err != nil {
// 		return keyVersion, err
// 	}
// 	return key.Key.KID.Version(), nil
// }
