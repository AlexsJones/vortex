package secrets

import (
	"errors"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/config"
)

func VaultFetchSecret(storepath, key string) (string, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return "", err
	}

	if client.Token() == "" {
		helper, err := config.DefaultTokenHelper()
		if err != nil {
			return "", err
		}

		token, err := helper.Get()
		if err != nil {
			return "", err
		}

		if token != "" {
			client.SetToken(token)
		}
	}

	data, err := client.Logical().Read(storepath)
	if err != nil {
		return "", err
	}
	if data == nil || data.Data == nil {
		return "", errors.New("No data returned for the given key")
	}
	if _, exist := data.Data[key]; !exist {
		return "", errors.New("Key does not exist inside the returned data map")
	}
	// Sneaky like the ninja
	return data.Data[key].(string), nil
}
