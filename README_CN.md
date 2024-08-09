# cwgo

中文 | [English](./README.md)

cwgo 是 CloudWeGo All in one 代码生成工具，整合了 kitex 和 hz 工具的优势，以提高开发者的编码效率和使用体验。其主要功能特点如下：

## 工具特点

- 支持生成工程化模板

  cwgo 工具支持生成 MVC 项目 Layout，用户只需要根据不同目录的功能，在相应的位置完成自己的业务代码即可，聚焦业务逻辑。

- 支持生成 Server、Client 代码

  cwgo 工具支持生成 Kitex、Hertz 的 Server 和 Client 代码，提供了对 Client 的封装。用户可以开箱即用的调用下游，免去封装 Client 的繁琐步骤。

- 支持生成关系型数据库代码

  cwgo 工具支持生成关系型数据库 CURD 代码。用户无需再自行封装繁琐的 CURD 代码，提高用户的工作效率。

- 支持生成文档类数据库代码

  cwgo 工具支持基于 IDL (thrift/protobuf) 生成文档类数据库 CURD 代码，目前支持 MongoDB。用户无需再自行封装繁琐的 CURD 代码，提高用户的工作效率。

- 支持生成命令行自动补全脚本

  cwgo 工具支持生成命令行自动补全脚本，提高用户命令行编写的效率。

- 支持分析 Hertz 项目路由和(路由注册)代码的关系

  cwgo 支持通过分析 Hertz 项目代码获取路由和(路由注册)代码的关系。

- 支持回退为 kitex、Hz 工具

  如果之前是 kitex、Hz 工具的用户，仍然可以使用 cwgo 工具。cwgo 工具支持回退功能，可以当作 kitex、Hz 使用，真正实现一个工具生成所有。

## 安装 cwgo 工具

```shell
# Go 1.15 及之前版本
GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get github.com/cloudwego/cwgo@latest

# Go 1.16 及以后版本
GOPROXY=https://goproxy.cn/,direct go install github.com/cloudwego/cwgo@latest
```

## 详细文档

### [快速开始](https://www.cloudwego.io/zh/docs/cwgo/getting-started/)

### 命令行工具

包含命令行工具形式及使用，详见[文档](https://www.cloudwego.io/zh/docs/cwgo/tutorials/cli/)

### 模板拓展

包含用户自定义模板的使用，详见[文档](https://www.cloudwego.cn/zh/docs/cwgo/tutorials/templete-extension/)

### Layout

Layout 生成及 Layout 介绍，详见[文档](https://www.cloudwego.io/zh/docs/cwgo/tutorials/layout/)

### Client

包含封装后的 Client 的生成和使用，详见[文档](https://www.cloudwego.io/zh/docs/cwgo/tutorials/client/)

### DB

包含如何使用 cwgo 工具生成关系型数据库 CURD 代码，详见[文档](https://www.cloudwego.io/zh/docs/cwgo/tutorials/db/)

### Doc 

包含如何使用 cwgo 工具生成文档型数据库 CURD 代码，详见[文档](https://www.cloudwego.cn/zh/docs/cwgo/tutorials/doc/)

### Api-list

支持分析 Hertz 项目代码获取路由和(路由注册)代码的关系，详见[文档](https://www.cloudwego.io/zh/docs/cwgo/tutorials/api-list)

### Server

包含如何生成 RPC Server、HTTP Server 代码，详见[文档](https://www.cloudwego.cn/zh/docs/cwgo/tutorials/server/)

### 命令行自动补全

包含如何启用命令行自动补全功能，详见[文档](https://www.cloudwego.cn/zh/docs/cwgo/tutorials/auto-completion/)

## 开源许可

cwgo 基于[Apache License 2.0](https://github.com/cloudwego/cwgo/blob/main/LICENSE) 许可证，其依赖的三方组件的开源许可见 [Licenses](https://github.com/cloudwego/cwgo/blob/main/licenses)。

## 联系我们

- Email: conduct@cloudwego.io
- 如何成为 member: [COMMUNITY MEMBERSHIP](https://github.com/cloudwego/community/blob/main/COMMUNITY_MEMBERSHIP.md)
- Issues: [Issues](https://github.com/cloudwego/cwgo/issues)
- Discord: 加入我们的 [Discord 频道](https://discord.gg/jceZSE7DsW)
- 飞书用户群（[注册飞书](https://www.larksuite.com/zh_cn/download)进群）

  ![LarkGroup](images/lark_group_cn.png)

## Landscapes

<p align="center">
<img src="https://landscape.cncf.io/images/cncf-landscape-horizontal-color.svg" width="150"/>&nbsp;&nbsp;<img src="https://www.cncf.io/wp-content/uploads/2023/04/cncf-main-site-logo.svg" width="200"/>
<br/><br/>
CloudWeGo 丰富了 <a href="https://landscape.cncf.io/">CNCF 云原生生态</a>。
</p>
