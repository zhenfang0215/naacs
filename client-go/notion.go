package notesasacnofigserver

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/dstotijn/go-notion"
)

// 先在 wolai 手动操作关联继承 https://developers.notion.com/docs/create-a-notion-integration
// 根据 pageid 获取下面的 databaseid https://developers.notion.com/docs/working-with-page-content
//
//	--- 或者直接粘贴 full page 的 dtabaseid
//
// 然后查询 database 下面的 pageid list https://developers.notion.com/docs/working-with-page-content
// 根据 pageid 获取页面的 blocks https://developers.notion.com/reference/block，可以获取到内容
const (
	NOTION_APP_SECRET  = "NOTION_APP_SECRET"
	NOTION_DATABASE_ID = "NOTION_DATABASE_ID"
)

type NotionProvider struct {
	*notion.Client
	config *NotionProviderConfig
	ctx    context.Context
}
type NotionProviderConfig struct {
	AppSecret  string
	Env        Environment
	DatabaseId string
}

func (*NotionProviderConfig) GetProvider() NotesProvider {
	return Notion_NoteProvider
}

func NewNotionProvider(config *NotionProviderConfig) *NotionProvider {
	client := notion.NewClient(config.AppSecret)
	return &NotionProvider{
		Client: client,
		config: config,
		ctx:    context.Background(),
	}
}

func (np *NotionProvider) GetConfig(appName string, target interface{}) error {
	databases, err := np.QueryDatabase(np.ctx, np.config.DatabaseId, &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: DATABASE_COLUMN_ENV,
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				Select: &notion.SelectDatabaseQueryFilter{
					Equals: np.config.Env.String(),
				},
			},
			Timestamp: "",
			Or:        nil,
			And:       nil,
		},
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	})
	if err != nil {
		return err
	}
	// 找到 pageid
	pageId := ""
	for _, database := range databases.Results {
		databaseProperties, ok := database.Properties.(notion.DatabasePageProperties)
		if !ok {
			return ErrColumnError
		}
		appNameColumn, appNameOk := databaseProperties[DATABASE_COLUMN_APP_NAME]
		if !(appNameOk) {
			return ErrColumnError
		}
		if appNameColumn.Title[0].PlainText == appName {
			pageId, err = extractNotionID(database.URL)
			if err != nil {
				return err
			}
		}
	}
	if pageId == "" {
		return fmt.Errorf("no pageid")
	}
	blocks, err := np.FindBlockChildrenByID(np.ctx, pageId, &notion.PaginationQuery{})
	if err != nil {
		return err
	}
	text := ""
	for _, block := range blocks.Results {
		code, ok := block.(*notion.CodeBlock)
		if ok {
			text = code.RichText[0].Text.Content
			goto end
		}
	}
	if text == "" {
		return fmt.Errorf("no code block")
	}
end:
	err = json.Unmarshal([]byte(text), target)
	if err != nil {
		return err
	}
	return nil
}

func extractNotionID(url string) (string, error) {
	// 定义正则表达式，提取最后一个 `-` 后面的部分
	re := regexp.MustCompile(`-(\w+)$`)

	// 执行正则匹配
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("无法从 URL 中提取 ID")
	}

	return matches[1], nil
}
