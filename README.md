# cwgo

[中文](./README_CN.md) | English

cwgo is an all-in-one code generation tool for CloudWeGo. It integrates the advantages of various components to improve
the developer experience. The main features of cwgo tool are as follows:

## Tool Characteristics

- Support for generating engineering templates

  The cwgo tool supports the generation of MVC project layout. Users only need to complete their own business code in the corresponding position according to the functions of different directories, focusing on business logic.

- Support generating Server and Client code

  The cwgo tool supports generating Server and Client codes of Kitex and Hertz, and provides an encapsulation of Client. Users can call downstream out of the box, eliminating the cumbersome steps of encapsulating the Client.

- Support for generating relational database code

  The cwgo tool supports generating relational database CURD code. Users no longer need to encapsulate the cumbersome CURD code by themselves, which improves the user's work efficiency.

- Support for generating document database code

  The cwgo tool supports generating document database CURD code based on IDL (Thrift/protobuf), and currently supports Mongodb. Users no longer need to encapsulate the cumbersome CURD code by themselves, which improves the user's work efficiency.

- Support for generating command line automatic completion scripts

  The cwgo tool supports generating command line completion scripts to improve the efficiency of user command line writing.

- Support fallback to Kitex, Hz tools

  If you were a Kitex or Hz user before, you can still use the cwgo tool. The cwgo tool supports the fallback function and can be used as Kitex and Hz, truly realizing a tool to generate all.

## Install cwgo Tool

```bash
# Go 1.15 and earlier version
GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get github.com/cloudwego/cwgo@latest

# Go 1.16 and later version
GOPROXY=https://goproxy.cn/,direct go install github.com/cloudwego/cwgo@latest
```

## Detailed Documentation

### [Quick Start](https://www.cloudwego.io/docs/cwgo/getting-started/)

### Command Line Tool

Contains detailed documentation on how cwgo CLI works, see this [document](https://www.cloudwego.io/docs/cwgo/tutorials/cli/)

### Template Extension

Instructions on how to customize templates can be found in the [document](https://www.cloudwego.io/docs/cwgo/tutorials/templete-extension/).

### Layout

This documents explains Layouts discussing how layout are generated, see this [document](https://www.cloudwego.io/docs/cwgo/tutorials/layout/).

### Client

The document details on how Generated Clients that have been Encapsulated Can be used, see this [document](https://www.cloudwego.io/docs/cwgo/tutorials/client/)

### DB

Details information containing HOW TO Use cwgo TOOL  To Generate relational Curd Codes, see this [document](https://www.cloudwego.io/docs/cwgo/tutorials/db/)

### Doc

Including how to use the cwgo tool to generate document database CURD code, see this [document](https://www.cloudwego.cn/docs/cwgo/tutorials/doc/).

### Api-list

supports getting the relationship between routes and (route registration) code by analyzing Hertz project code, see this [document](https://www.cloudwego.io/docs/cwgo/tutorials/api-list)

### Server

Including how to generate RPC Server and HTTP Server code, see this [document](https://www.cloudwego.cn/docs/cwgo/tutorials/server/).

### auto-completion

Including how to enable command line auto completion function, see this [document](https://www.cloudwego.cn/docs/cwgo/tutorials/auto-completion/).

## Open Source License

cwgo is based on Apache License 2.0, [Apache License](https://github.com/cloudswego/cwgo/blob/main/LICENSE). Its dependent
third-party component open-source licenses will include Licenses.

## Contact Us

- Email: conduct@cloudwego.io
- How to become a member: [COMMUNITY MEMBERSHIP](https://github.com/cloudwego/community/blob/main/COMMUNITY_MEMBERSHIP.md)
- Issues: [Issues](https://github.com/cloudwego/cwgo/issues)
- Slack: Join our [Slack channel](https://join.slack.com/t/cloudwego/shared_invite/zt-tmcbzewn-UjXMF3ZQsPhl7W3tEDZboA)
- Feishu group (Register for [Feishu](https://www.larksuite.com/en-US/download) and join the group)

  ![LarkGroup](images/lark_group.png)
