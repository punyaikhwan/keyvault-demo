package keyvault

import (
	"keyvault-demo/config"
	"keyvault-demo/internal/pkg/keyvault"
	"keyvault-demo/internal/pkg/keyvault/azurekeyvault"
)

type keyVault struct {
	kv keyvault.KeyVault
}

var instance *keyVault

func KeyVault() keyvault.KeyVault {
	if instance == nil {
		kv := azurekeyvault.NewAzureKeyVault(config.Configuration().VaultURL)
		instance = &keyVault{kv}
	}

	return instance.kv
}
