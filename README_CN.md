# cwgo 

中文 | [English](./README.md)

cwgo 是 CloudWeGo All in one 代码生成工具，整合了各个组件的优势，提高开发者提体验。
cwgo 工具可以方便生成工程化模版，其主要功能特点如下：

## 工具特点
- 用户友好生成方式

  cwgo 工具同时提供了交互式命令行和静态命令行两种方式。交互式命令行可以低成本生成代码，不用再去关心传递什么参数，也不用执行多次命令，
适合大部分用户；而对高级功能有需求的用户，仍可使用常规的静态命令行构造生成命令。
- 支持生成工程化模板

  cwgo 工具支持生成 MVC 项目 Layout，用户只需要根据不同目录的功能，在相应的位置完成自己的业务代码即可，聚焦业务逻辑。
- 支持生成 Server、Client 代码

  cwgo 工具支持生成 Kitex、Hertz 的 Server 和 Client 代码，提供了对 Client 的封装。用户可以开箱即用的调用下游，免去封装 Client 的繁琐步骤
- 支持生成数据库代码

  cwgo 工具支持生成数据库 CURD 代码。用户无需再自行封装繁琐的 CURD 代码，提高用户的工作效率。
- 支持回退为 Kitex、Hz 工具

  如果之前是 Kitex、Hz 的用户，仍然可以使用 cwgo 工具。cwgo 工具支持回退功能，可以当作 Kitex、Hz 使用，真正实现一个工具生成所有。

## 安装 cwgo 工具
```
# Go 1.15 及之前版本
GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get github.com/cloudwego/cwgo@latest

# Go 1.16 及以后版本
GOPROXY=https://goproxy.cn/,direct go install github.com/cloudwego/cwgo@latest
```

## 详细文档
### [快速开始](https://www.cloudwego.io/zh/docs/cwgo/getting-started/)
### 命令行工具
包含命令行工具形式及使用，详见[文档](https://www.cloudwego.io/zh/docs/cwgo/tutorials/cli/)
### Layout
Layout 生成及 Layout 介绍，详见[文档](https://www.cloudwego.io/zh/docs/cwgo/tutorials/layout/)
### Client
包含封装后的 Client 的生成和使用，详见[文档](https://www.cloudwego.io/zh/docs/cwgo/tutorials/client/)
### DB
包含如何使用 cwgo 工具生成 CURD 代码，详见[文档](https://www.cloudwego.io/zh/docs/cwgo/tutorials/db/)
### 模版扩展
包含如何自定义模板，详见[文档](https://www.cloudwego.io/zh/docs/cwgo/tutorials/templete-extension/)


## 开源许可

cwgo 基于[Apache License 2.0](https://github.com/cloudwego/cwgo/blob/main/LICENSE) 许可证，其依赖的三方组件的开源许可见 [Licenses](https://github.com/cloudwego/cwgo/blob/main/licenses)。

## 联系我们
- Email: conduct@cloudwego.io
- 如何成为 member: [COMMUNITY MEMBERSHIP](https://github.com/cloudwego/community/blob/main/COMMUNITY_MEMBERSHIP.md)
- Issues: [Issues](https://github.com/cloudwego/cwgo/issues)
- Slack: 加入我们的 [Slack 频道](https://join.slack.com/t/cloudwego/shared_invite/zt-tmcbzewn-UjXMF3ZQsPhl7W3tEDZboA)
- 飞书用户群（[注册飞书](https://www.larksuite.com/zh_cn/download)进群）

  ![LarkGroup](images/lark_group_cn.png)

## Landscapes

<p align="center">
<img src="https://landscape.cncf.io/images/left-logo.svg" width="150"/>&nbsp;&nbsp;<img src="https://landscape.cncf.io/images/right-logo.svg" width="200"/>
<br/><br/>
CloudWeGo 丰富了 <a href="https://landscape.cncf.io/">CNCF 云原生生态</a>。
</p>
