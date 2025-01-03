# Nacs - Turn Notes into a Configuration Server

- [中文介绍](https://github.com/zhenfang0215/nacs/blob/main/README.md)

![nacs](https://github.com/zhenfang0215/nacs/blob/main/doc/nacs-github-banner-img.png)


Nacs is a tool that allows you to use your notes as configuration files. Have you ever had these concerns😕:
1. Are you still feeling insecure and inelegant because your configuration files are scattered everywhere? 
2. Are you still troubled by not being able to find a suitable configuration center for your project?
3. Are you still feeling that it's not cost-effective to set up a heavy configuration center? 

Now you've found the solution. You can store your service configurations in your note-taking software and then use Nacs to retrieve configuration information anywhere. This is really🤣:
1. Lightweight (no need to deploy a service)
2. Centralized (all in your notes)
3. Elegant (in my opinion)
4. Secure (as stated by the note-taking software)

- [Supported Note Brands](#supported-note-brands)
- [Database Format Required for Notes](#database-format-required-for-notes)
- [Usage Method](#usage-method)
- [Usage Method for Each Note Brand](#usage-method-for-each-note-brand)
  - [Notion](#notion)
  - [Wolai](#wolai)


## Supported Note Brands
| Brand | Official Website | API Documentation |
|:--:|:--:|:--:|
| notion | [Official Website](https://www.notion.com/) | [API Documentation](https://developers.notion.com/docs/getting-started)|
| wolai | [Official Website](https://www.wolai.com/) | [API Documentation](https://www.wolai.com/7FB9PLeqZ1ni9FfD11WuUi)|


## Database Format Required for Notes
Here are three fixed requirements, please forgive me, experts: Currently, only configuration files can be set through a database with a fixed format; and you need to manually create the database; and the internal configuration can only be code, including JSON format configuration information. The fixed format of this database and table is as follows:
![database_table](https://github.com/zhenfang0215/nacs/blob/main/doc/database_sample.png)


## Usage Method
1. Configure the JSON configuration information in your note-taking software.
2. Then import the project in your code.
3. Initialize a client in the code based on different note brands.
```go
// Notion
config := &NotionProviderConfig{
    AppSecret:  "your secret",
    Env:        Dev_Environment,  // dev environment
    DatabaseId: "your database id",
}
// wolai
// config := &WolaiProviderConfig{
//    .....
//}

client := NewNaaCSClient(config)
```
4. Then define a struct corresponding to the configuration file for serialization.

```go
type TestSetting struct {
	K1 string   `json:"k1"`
	K2 []string `json:"k2"`
}

cfg := &TestSetting{}
client.GetConfig("project", cfg)

print cfg


> {"k1":"v1","k2":["v2.2","v2.3"]}

```

## Usage Method for Each Note Brand
### Notion
1. You need to create an integration in Notion: https://www.notion.so/profile/integrations
2. Create a database and obtain the database ID. The method to obtain it is: https://developers.notion.com/docs/working-with-databases#adding-pages-to-a-database
3. Bind the above common integration to the page where the database is located: https://developers.notion.com/docs/create-a-notion-integration
### Wolai
1. Create a Wolai application and obtain the app key and app secret: https://www.wolai.com/dev
2. Create a database and obtain the database block ID. The method to obtain it is: https://www.wolai.com/wolai/2kRSq4mVwxCUUcUhrgnQgp

