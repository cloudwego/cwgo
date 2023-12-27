# cwgo

[中文](./README_CN.md) | English

cwgo is an all-in-one code generation tool for CloudWeGo. It integrates the advantages of various components to improve
the developer experience. The main features of cwgo tool are as follows:

## Tool Characteristics


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
# Go 1.15 and earlier version
GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get github.com/cloudwego/cwgo@latest

# Go 1.16 and later version
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

### How to enable auto-completion
#### Supported in Bash
First download Bash [autocomplete script](https://github.com/urfave/cli/blob/v2-maint/autocomplete/bash_autocomplete), assuming it is downloaded to the autocomplete/ folder in the root directory of the project (you can Define location)
##### Temporary support for completion
```shell
PROG=cwgo

source ./autocomplete/bash_autocomplete
```
##### Permanent support for completion
```shell
sudo cp autocomplete/bash_autocomplete /etc/bash_completion.d/cwgo

source /etc/bash_completion.d/cwgo
```
#### Supported in Zsh
First download Zsh [autocomplete script](https://github.com/urfave/cli/blob/v2-maint/autocomplete/zsh_autocomplete), assuming it is downloaded to the autocomplete/ folder in the root directory of the project (you can Define location)
##### Temporary support for completion
```shell
PROG=cwgo

source ./autocomplete/zsh_autocomplete
```
#### Supported in PowerShell
First download PowerShell [autocomplete script](https://github.com/urfave/cli/blob/v2-maint/autocomplete/powershell_autocomplete.ps1), assuming it is downloaded to the autocomplete/ folder in the root directory of the project ( Customizable location)
##### Temporary support for completion
First, rename the downloaded powershell_autocomplete.ps1 to cwgo.ps1.

Then execute:
```shell
& autocomplete/cwgo.ps1
```

##### Permanent support for completion
open the $profile.

Add a line inside:
```shell
& path/to/autocomplete/cwgo.ps1
```
Note that the name and path of the ps1 script must be correctly configured here, and then permanent auto-completion can be performed.

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
