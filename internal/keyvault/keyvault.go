package keyvault

import (
	"keyvault-demo/config"
	"keyvault-demo/internal/pkg/keyvault"
	"keyvault-demo/internal/pkg/keyvault/azurekeyvault"
	"keyvault-demo/internal/pkg/keyvault/hashicorptransit"
)

const (
	AZURE     = "azure"
	HASHICORP = "hashicorp"
)

type keyVault struct {
	kv keyvault.KeyVault
}

var instance *keyVault

func KeyVault() keyvault.KeyVault {
	if instance == nil {
		switch config.Configuration().VaultProvider {
		case AZURE:
			kv := azurekeyvault.NewAzureKeyVault(config.Configuration().VaultURL)
			instance = &keyVault{kv}
		default:
			kv := hashicorptransit.NewHashicorpTransitEngine(config.Configuration().VaultURL, config.Configuration().VaultUsername, config.Configuration().VaultPassword, config.Configuration().HashicorpTransitPath)
			instance = &keyVault{kv}
		}
	}

	return instance.kv
}
