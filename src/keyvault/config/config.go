package config

const (
	KeyVaultScheme  = "https"
	KeyVaultBaseUrl = "vault.azure.cn"
	GrantType       = "client_credentials"
)

type Config struct {
	Name         string `yaml:"name"`
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}
