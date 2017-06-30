package secret

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"keyvault/config"
	"net/http"
)

type SecretResp struct {
	Value string `json:"value"`
	Id    string `json:"id"`
}

func GetSecret(cfg config.Config, token string, key string) string {
	client := http.Client{}
	secretUrl := fmt.Sprintf("%s://%s.%s/secrets/%s?api-version=2016-10-01", config.KeyVaultScheme, cfg.Name, config.KeyVaultBaseUrl, key)
	req, _ := http.NewRequest(http.MethodGet, secretUrl, nil)
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("get secret error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Printf("get secret error: %v\n", string(data))
		return ""
	}
	var respData SecretResp
	json.Unmarshal(data, &respData)

	return respData.Value
}
