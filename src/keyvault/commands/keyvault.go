package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"keyvault/auth"
	"keyvault/config"
	"keyvault/secret"
	"os"
	"strings"
	"text/template"
)

var (
	configPath   string
	renderTpl    string
	renderOutput string
	keyStr       string
	cfg          config.Config
)

var KeyvaultCmd = &cobra.Command{Use: "keyvault"}

var renderCmd = &cobra.Command{
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
		token, err := auth.GetToken(cfg)
		if err != nil {
			fmt.Printf("get token error: %v\n", err)
			return
		}

		keys := strings.Split(keyStr, ",")
		fmt.Printf("keys %v\n", keys)
		var data = make(map[string]string)
		for _, key := range keys {
			secret := secret.GetSecret(cfg, token, key)
			fmt.Printf(" --- kEY: %s --- \n", key)
			data[key] = secret
		}

		tpl.Execute(out, data)
		fmt.Printf("END\n")
	},
}

func init() {
	renderCmd.PersistentFlags().StringVar(&configPath, "config", "", "config path")
	renderCmd.PersistentFlags().StringVar(&renderTpl, "tpl", "config.tpl", "template path")
	renderCmd.PersistentFlags().StringVar(&renderOutput, "out", "config.out", "output path")
	renderCmd.PersistentFlags().StringVar(&keyStr, "keys", "", "keyvault keys")
}

func Execute() {
	KeyvaultCmd.AddCommand(renderCmd)

	if err := KeyvaultCmd.Execute(); err != nil {
		fmt.Printf("commnd error: %v\n", err)
		os.Exit(1)
	}
}

func configSetup(configPath string) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("config error: %v\n", err)
		os.Exit(1)
	}
	yaml.Unmarshal(data, &cfg)
}
