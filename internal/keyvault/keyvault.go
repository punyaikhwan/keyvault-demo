package keyvault

import (
	"keyvault-demo/config"
	"keyvault-demo/internal/pkg/keyvault"
	"keyvault-demo/internal/pkg/keyvault/hashicorptransit"
)

type keyVault struct {
	kv keyvault.KeyVault
}

var instance *keyVault

func KeyVault() keyvault.KeyVault {
	if instance == nil {
		kv := hashicorptransit.NewHashicorpTransitEngine(config.Configuration().VaultURL, config.Configuration().VaultUsername, config.Configuration().VaultPassword)
		instance = &keyVault{kv}
	}

	return instance.kv
}
