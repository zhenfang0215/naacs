package notesasacnofigserver

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNotionClient(t *testing.T) {
	config := &NotionProviderConfig{
		AppSecret:  "",
		Env:        Prod_Environment,
		DatabaseId: "",
	}
	client := NewNaaCSClient(config)
	cfg := &TestSetting{}
	err := client.GetConfig("project", cfg)

	cfgBytes, _ := json.Marshal(cfg)
	fmt.Println(string(cfgBytes))
	fmt.Println(err)
	fmt.Println(cfg)
}

type TestSetting struct {
	K1 string   `json:"k1"`
	K2 []string `json:"k2"`
}

func TestWolaiClient(t *testing.T) {
	config := &WolaiProviderConfig{
		AppId:     "",
		AppSecret: "",
		Env:       Dev_Environment,
		BlockId:   "",
	}
	client := NewNaaCSClient(config)
	cfg := &TestSetting{}
	err := client.GetConfig("project", cfg)
	fmt.Println(err)
	fmt.Println(cfg)
}
