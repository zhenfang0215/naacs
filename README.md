# Nacs - 将笔记变成配置服务器

![nacs](https://github.com/zhenfang0215/nacs/blob/main/doc/nacs-github-banner-img.png)

Nacs 是一款让你可以把笔记作为配置文件的工具。你是否有这样的困惑😕:
1. 你还在为配置文件散落在各种地方而感到不安全且不优雅吗? 
2. 你还在为自己的项目找不到合适的配置中心而苦恼吗？
3. 你还在为搭建一套繁重的配置中心而且不划算吗? 

现在你找到解决方案了,你可以把自己服务的配置在自己的笔记软件中,然后使用 nacs 在任意地方获取配置信息。这真的是🤣:
1. 轻量(不用部署服务)
2. 集中(都在你的笔记中)
3. 优雅(自认为)
4. 安全(转述笔记软件自己说的)

- [支持的笔记品牌](#支持的笔记品牌)
- [笔记需要设置的database格式](#笔记需要设置的 database 格式)
- [使用方式](#使用方式)
- [各家笔记使用方式](#各家笔记使用方式)
  - [Notion](#notion)
  - [Wolai](#wolai)


## 支持的笔记品牌
|品牌|官网地址|api文档|
|:--:|:--:|:--:|
| notion |[官网地址](https://www.notion.com/)|[api文档](https://developers.notion.com/docs/getting-started)|
| wolai |[官网地址](https://www.wolai.com/)|[api文档](https://www.wolai.com/7FB9PLeqZ1ni9FfD11WuUi)|


## 笔记需要设置的 database 格式
接下来是三个固定要求,请大佬轻喷: 目前只支持通过格式固定的 database 设置配置文件; 并且自己手动创建 database; 并且内部配置只能是 code 包括 json 格式的配置信息。 这个 database 和 table 的固定格式如下:
![database_table](https://github.com/zhenfang0215/nacs/blob/main/doc/database_sample.png)


## 使用方式
1. 在自己的笔记软件中配置好 json 配置信息
2. 然后在代码中引入该项目
```
import "github.com:zhenfang0215/nacs"
```
3. 代码中根据不同的笔记品牌初始化一个 client
```go
// Notion
config := &NotionProviderConfig{
    AppSecret:  "your secret",
    Env:        Dev_Environment,  // dev 环境
    DatabaseId: "your database id",
}
// wolai
// config := &WolaiProviderConfig{
//    .....
//}

client := NewNaaCSClient(config)
```
4. 然后定义一个和配置文件对应的结构体,用于序列化
```golang
type TestSetting struct {
	K1 string   `json:"k1"`
	K2 []string `json:"k2"`
}

cfg := &TestSetting{}
client.GetConfig("project", cfg)

print cfg


> {"k1":"v1","k2":["v2.2","v2.3"]}
```

## 各家笔记使用方式
#### Notion
1. 需要到 wolai 中创建一个集成: https://www.notion.so/profile/integrations
2. 创建一个 database, 并获取 database id,获取方式: https://developers.notion.com/docs/working-with-databases#adding-pages-to-a-database
3. 给 database 所在页面绑定上面常见的集成: https://developers.notion.com/docs/create-a-notion-integration

#### Wolai
1. 创建一个 wolai 应用,拿到 app key 和 app secret: https://www.wolai.com/dev
2. 创建一个 database, 并获取 database block id, 获取方式:https://www.wolai.com/wolai/2kRSq4mVwxCUUcUhrgnQgp
