package notesasacnofigserver

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/lemonnekogh/guolai"
)

const (
	WOLAI_APP_ID     = "WOLAI_APP_ID"
	WOLAI_APP_SECRET = "WOLAI_APP_SECRET"
	CONFIG_BLOCK_ID  = "CONFIG_BLOCK_ID"
	ENV              = "ENV"

	WOLAI_DB_COLUMN_APP_NAME = "app_name"
	WOLAI_DB_COLUMN_ENV      = "env"
)

type WolaiProvider struct {
	config     *WolaiProviderConfig
	woLaiToken string
	client     *guolai.WolaiAPI
}
type WolaiProviderConfig struct {
	AppId     string
	AppSecret string
	Env       Environment
	BlockId   string
}

func (*WolaiProviderConfig) GetProvider() NotesProvider {
	return Wolai_NoteProvider
}

func NewWolaiProvider(config *WolaiProviderConfig) *WolaiProvider {
	if config.AppId == "" {
		config.AppId = os.Getenv(WOLAI_APP_ID)
	}
	if config.AppSecret == "" {
		config.AppSecret = os.Getenv(WOLAI_APP_SECRET)
	}
	if config.BlockId == "" {
		config.BlockId = os.Getenv(CONFIG_BLOCK_ID)
	}
	if config.Env == 0 {
		config.Env = Dev_Environment
		if os.Getenv(ENV) == "prod" || os.Getenv(ENV) == "production" {
			config.Env = Prod_Environment
		}
	}
	token, err := CreateWolaiToken(config.AppId, config.AppSecret)
	if err != nil {
		panic(err)
	}
	client := guolai.New(token)
	return &WolaiProvider{
		config:     config,
		woLaiToken: token,
		client:     client,
	}
}

func (naacc *WolaiProvider) GetConfig(appName string, target interface{}) error {
	configItems, err := naacc.getConfigFromWolaiDatabase()
	if err != nil {
		return err
	}
	configs := filterWolaiDatabaseItem(configItems, func(drd *guolai.DatabaseRowData) bool {
		return drd.Data[WOLAI_DB_COLUMN_APP_NAME].Value == appName
	})
	if len(configs) == 0 {
		return ErrConfigNotFound
	}
	pageId := configs[0].PageId
	pageApiResp, err := naacc.getConfigContentFromWolaiPage(pageId)
	if err != nil {
		return err
	}
	if len(pageApiResp) == 0 {
		return ErrConfigNotFound
	}

	var targetItem guolai.RichText
	for _, item := range pageApiResp {
		if item.Type == "code" {
			for _, content := range item.Content {
				if content.Type == "text" {
					targetItem = content
					goto end
				}
			}
		}
	}
end:
	err = json.Unmarshal([]byte(targetItem.Title), target)
	if err != nil {
		return err
	}
	return nil
}

func (naacc *WolaiProvider) getConfigFromWolaiDatabase() ([]guolai.DatabaseRowData, error) {
	blockApiResp, err := naacc.client.GetDatabase(naacc.config.BlockId)
	if err != nil {
		return nil, err
	}
	items := filterWolaiDatabaseItem(blockApiResp.Rows, func(drd *guolai.DatabaseRowData) bool {
		return drd.Data[WOLAI_DB_COLUMN_ENV].Value == naacc.config.Env.String()
	})
	return items, nil
}

func (naacc *WolaiProvider) getConfigContentFromWolaiPage(pageId string) ([]guolai.BlockApiResponse, error) {
	blockChildrenApiResp, err := naacc.client.GetBlockChildren(pageId)
	if err != nil {
		return nil, err
	}

	return blockChildrenApiResp, nil
}

func filterWolaiDatabaseItem(rows []guolai.DatabaseRowData, filter func(*guolai.DatabaseRowData) bool) []guolai.DatabaseRowData {
	var result []guolai.DatabaseRowData
	for _, row := range rows {
		if ok := filter(&row); ok {
			result = append(result, row)
		}
	}
	return result
}

func CreateWolaiToken(appId string, appSecret string) (string, error) {
	createTokenResp, err := guolai.CreateToken(appId, appSecret)
	if err != nil {
		return "", err
	}
	if createTokenResp != nil && createTokenResp.AppToken != "" {
		return createTokenResp.AppToken, nil
	}
	return "", errors.New("CreateToken Resp is nil")
}
