package secrets

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/hashicorp/vault/api"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate = validator.New()
)

type vault struct {
	token   string `validate:"required"`
	address string `validate:"gt=1"`
}

// NewVaultFromEnv creates a vault SecretObtainer
// using the usual environment settings that the vault cli would use.
func NewVaultFromEnv() (*vault, error) {
	v := &vault{
		address: "http://localhost:8200",
	}
	if f, err := os.Open(path.Join(os.Getenv("HOME"), ".vault-token")); !os.IsNotExist(err) {
		buff, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		v.token = string(buff)
	}
	if token, set := os.LookupEnv("VAULT_TOKEN"); set && len(v.token) == 0 {
		v.token = token
	}
	if addr, set := os.LookupEnv("VAULT_ADDR"); set {
		v.address = addr
	}
	return v, validate.Struct(v)
}

func VaultFetchSecret(storepath, key string) (string, error) {
	v, err := NewVaultFromEnv()
	if err != nil {
		return "", err
	}
	client, err := api.NewClient(&api.Config{
		Address: v.address,
	})
	if err != nil {
		return "", err
	}
	client.SetToken(v.token)
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
