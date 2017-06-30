package auth

import (
	"fmt"
	"net/url"
	"regexp"
)

import (
	"encoding/json"
	"io/ioutil"
	"keyvault/config"
	"net/http"
	"strings"
)

type TokenResp struct {
	TokenType   string `json:"token_type""`
	Resource    string `json:"resource""`
	AccessToken string `json:"access_token""`
}

func GetToken(cfg config.Config) (string, error) {

	opsUrl := fmt.Sprintf("%s://%s.%s/secrets?api-version=2016-10-01", config.KeyVaultScheme, cfg.Name, config.KeyVaultBaseUrl)
	opsResp, err := http.Get(opsUrl)

	if err != nil {
		fmt.Printf("error %v\n", err)
		return "", err
	}
	defer opsResp.Body.Close()

	authInfo := opsResp.Header.Get("Www-Authenticate")
	authReg, _ := regexp.Compile(`authorization="[^"]+"`)
	resReg, _ := regexp.Compile(`resource="[^"]+"`)

	loginUrl := strings.Trim(authReg.FindString(authInfo)[14:], "\"")
	oauthUrl := loginUrl + "/oauth2/token"
	resUrl := strings.Trim(resReg.FindString(authInfo)[9:], "\"")

	fmt.Printf("oauth url: %v , resource url: %v\n", oauthUrl, resUrl)

	v := url.Values{}
	v.Set("grant_type", config.GrantType)
	v.Set("client_id", cfg.ClientId)
	v.Set("client_secret", cfg.ClientSecret)
	v.Set("resource", resUrl)

	authResp, err := http.PostForm(oauthUrl, v)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return "", err
	}
	defer authResp.Body.Close()

	data, _ := ioutil.ReadAll(authResp.Body)
	var respData TokenResp
	if err := json.Unmarshal(data, &respData); err != nil {
		fmt.Printf("error %v\n", err)
		return "", err
	}

	token := fmt.Sprintf("%s %s", respData.TokenType, respData.AccessToken)
	return token, nil
}
