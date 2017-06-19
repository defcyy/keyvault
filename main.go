package main

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"fmt"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"io/ioutil"
	"os"
	"text/template"
)

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

type TokenResp struct {
	TokenType   string `json:"token_type""`
	Resource    string `json:"resource""`
	AccessToken string `json:"access_token""`
}

type SecretResp struct {
	Value string `json:"value"`
	Id    string `json:"id"`
}

var config Config

func main() {

	var (
		configPath   string
		renderTpl    string
		renderOutput string
		keyStr       string
	)

	renderCmd := &cobra.Command{
		Use:   "render [key]",
		Short: "",
		Long:  "",
		PreRun: func(cmd *cobra.Command, args []string) {
			configSetup(configPath)
		},
		Run: func(cmd *cobra.Command, args []string) {
			tpl, err := template.ParseFiles(renderTpl)
			if err != nil {
				fmt.Printf("config templete error: %v\n", err)
				os.Exit(1)
			}

			out, _ := os.OpenFile(renderOutput, os.O_RDWR|os.O_CREATE, 0644)
			token, err := getToken()
			if err != nil {
				fmt.Printf("get token error: %v\n", err)
				return
			}

			keys := strings.Split(keyStr, ",")
			fmt.Printf("keys %v\n", keys)
			var data = make(map[string]string)
			for _, key := range keys {
				secret := getSecret(token, key)
				fmt.Printf(" --- kEY: %s --- \n", key)
				data[key] = secret
			}

			tpl.Execute(out, data)
			fmt.Printf("END\n")
		},
	}

	renderCmd.PersistentFlags().StringVar(&configPath, "config", "", "config path")
	renderCmd.PersistentFlags().StringVar(&renderTpl, "tpl", "config.tpl", "template path")
	renderCmd.PersistentFlags().StringVar(&renderOutput, "out", "config.out", "output path")
	renderCmd.PersistentFlags().StringVar(&keyStr, "keys", "", "keyvault keys")

	rootCmd := &cobra.Command{Use: "keyvault"}
	rootCmd.AddCommand(renderCmd)

	rootCmd.Execute()
}

func configSetup(configPath string) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("config error: %v\n", err)
		os.Exit(1)
	}
	yaml.Unmarshal(data, &config)
}

func getToken() (string, error) {

	opsUrl := fmt.Sprintf("%s://%s.%s/secrets?api-version=2016-10-01", KeyVaultScheme, config.Name, KeyVaultBaseUrl)
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
	v.Set("grant_type", GrantType)
	v.Set("client_id", config.ClientId)
	v.Set("client_secret", config.ClientSecret)
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

func getSecret(token string, key string) string {
	client := http.Client{}
	secretUrl := fmt.Sprintf("%s://%s.%s/secrets/%s?api-version=2016-10-01", KeyVaultScheme, config.Name, KeyVaultBaseUrl, key)
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
