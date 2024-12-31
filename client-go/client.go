package notesasacnofigserver

import (
	"fmt"
)

const (
	DATABASE_COLUMN_APP_NAME = "app_name"
	DATABASE_COLUMN_ENV      = "env"
)

type NotesProvider int8

const (
	Wolai_NoteProvider NotesProvider = iota + 1
	Notion_NoteProvider
)

type Environment int8

const (
	Dev_Environment Environment = iota + 1
	Prod_Environment
)

func (e Environment) String() string {
	if e == Prod_Environment {
		return "prod"
	}
	return "dev"
}

func NewEnvironmentFromString(env string) Environment {
	if env == "prod" {
		return Prod_Environment
	}
	return Dev_Environment
}

type NaaCSClient struct {
	Provider NotesProvider
}

type ProviderService interface {
	// GetConfig
	GetConfig(appName string, target interface{}) error
}

type ProviderConfig interface {
	GetProvider() NotesProvider
}

var (
	ErrConfigNotFound = fmt.Errorf("config not found")
	ErrColumnError    = fmt.Errorf("colunm error")
)

func NewNaaCSClient(cfg ProviderConfig) ProviderService {
	provider := cfg.GetProvider()
	if provider == Wolai_NoteProvider {
		config, ok := cfg.(*WolaiProviderConfig)
		if !ok {
			panic(fmt.Errorf("provider config error %d", provider))
		}
		return NewWolaiProvider(config)
	} else if provider == Notion_NoteProvider {
		config, ok := cfg.(*NotionProviderConfig)
		if !ok {
			panic(fmt.Errorf("provider config error %d", provider))
		}
		return NewNotionProvider(config)
	}
	panic(fmt.Errorf("not support provider %d", provider))
}
