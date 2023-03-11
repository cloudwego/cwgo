# cwgo

[中文](./README_CN.md) | English

cwgo is an all-in-one code generation tool for CloudWeGo. It integrates the advantages of various components to improve
the developer experience. The main features of cwgo tool are as follows:

## Tool Characteristics

- User-friendly generation method

  cwgo tool provides both interactive command line and static command line methods. The interactive command line can
  generate code at a low cost, without worrying about passing parameters or executing multiple commands, which is
  suitable for most users; while advanced users with specific needs can still use regular static commands to construct
  generation commands.

- Supports generating project templates

  cwgo tool supports generating MVC project layouts. Users only need to complete their own business logic in the
  corresponding locations according to the functionality of different directories, focusing on business logic.

- Supports generating server and client code

  cwgo supports generating Kitex and Hertz's server and client code, providing encapsulation for clients. Users can use
  it out of the box to call downstream services, saving them from cumbersome steps such as packaging clients.

- Supports generating database code

  cwgo tool supports generating database CURD (Create Update Read Delete) codes. Users no longer need to package
  tedious CURD codes themselves, thereby improving work efficiency.

- Support fallback to Kitex and Hz tools

  If you were a user of Kitex or Hz before, you can still use the cwgo tool. With its rollback function support backward
  compatibility with these tools.

## Install cwgo Tool

```
# Go version prior to 1.15 
GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get github.com/cloudwego/cwgo@latest


# Go version after using above than v1.16
GOPROXY=https://goproxy.cn/,direct go install github.com/cloudwego/cwgo@latest
```

## Detailed Documentation

### [Quick Start](https://www.cloudwego.io/docs/cwgo/getting-started/)

### Command Line Tool

Contains detailed documentation on how cwgo CLI works, see this [document](https://www.cloudwego.io/docs/cwgo/tutorials/cli/)

### Layout

This documents explains Layouts discussing how layout are generated, see this [document](https://www.cloudwego.io/docs/cwgo/tutorials/layout/).

### Client

The document details on how Generated Clients that have been Encapsulated Can be used, see this [document](https://www.cloudwego.io/docs/cwgo/tutorials/client/)

### DB

Details information containing HOW TO Use cwgo TOOL GEN To Generate Curd Codes, see this [document](https://www.cloudwego.io/docs/cwgo/tutorials/db/)

### Template Extension

Instructions on how to customize templates can be found in the [document](https://www.cloudwego.io/docs/cwgo/tutorials/templete-extension/).

## Open Source License

cwgo has got Apache License2.0 licensing
format, [Apache License](https://github.com/cloudswego/cwgo/blob/main/LICENSE). Its dependent
third-party component open-source licenses will include Licenses.


## Contact Us

- Email: conduct@cloudwego.io
- How to become a member: [COMMUNITY MEMBERSHIP](https://github.com/cloudwego/community/blob/main/COMMUNITY_MEMBERSHIP.md)
- Issues: [Issues](https://github.com/cloudwego/cwgo/issues)
- Slack: Join our [Slack channel](https://join.slack.com/t/cloudwego/shared_invite/zt-tmcbzewn-UjXMF3ZQsPhl7W3tEDZboA)
- Feishu group (Register for [Feishu](https://www.larksuite.com/en-US/download) and join the group)

  ![LarkGroup](images/lark_group.png)

- WeChat: CloudWeGo community

  ![WechatGroup](images/wechat_group_en.png)
