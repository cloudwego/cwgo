# DB
cwgo 集成了 gorm/gen 用于帮助用户生成 Model 代码以及基础的CURD代码。

# 基础命令

使用 `cwgo model -h` 查看使用详情

```
NAME:
   cwgo model - generate DB model

                Examples:
                  # Generate DB model code
                  cwgo  model --db_type mysql --dsn "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"


USAGE:
   cwgo model [command options] [arguments...]

OPTIONS:
   --c value                          Specify the config file path
   --dsn value                        Specify the database source name. (https://gorm.io/docs/connecting_to_the_database.html)
   --db_type value                    Specify database type. (mysql or sqlserver or sqlite or postgres) (default: mysql)
   --out_dir value                    Specify output directory (default: biz/dao/query)
   --out_file value                   Specify output filename (default: gen.go)
   --tables value [ --tables value ]  Specify databases tables
   --unittest                         Specify generate unit test (default: false)
   --only_model                       Specify only generate model code (default: false)
   --model_pkg value                  Specify model package name
   --nullable                         Specify generate with pointer when field is nullable (default: false)
   --type_tag                         Specify generate field with gorm column type tag (default: false)
   --index_tag                        Specify generate field with gorm index tag (default: false)
   --help, -h                         show help (default: false)
```




## 详细参数

```
   --c value                          指定配置文件路径
   --dsn value                        指定数据库DSN
   --db_type value                    指定数据库类型(mysql or sqlserver or sqlite or postgres) (默认 mysql)
   --out_dir value                    指定输出目录路径，默认为 biz/dao/query
   --out_file value                   指定输出文件名，默认为 gen.go
   --tables value                     指定数据库表，默认为全表
   --unittest                         指定是否生成单测，默认为 false
   --only_model                       指定是否生成仅 model，默认为 false 
   --model_pkg value                  指定 model 的包名
   --nullable                         指定生成字段是否为指针当字段为 nullable，默认为 false
   --type_tag                         指定字段是否生成 gorm 的 type tag，默认为 false
   --index_tag                        指定字段是否生成 gorm 的 index tag，默认为 false
```




## 用法示例

```
cwgo  model --db_type mysql --dsn "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
```

or

```
cwgo model --c ./gen.yml
```

```yaml
version: "0.1"
database:
  # consult[https://gorm.io/docs/connecting_to_the_database.html]"
  dsn : "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
  # input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]
  db  : "mysql"
  # enter the required data table or leave it blank.You can input : orders,users,goods
  tables  : "table1,table2"
  # only generate model
  onlyModel: false
  # specify a directory for output
  outPath :  "./dao/query"
  # query code file name, default: gen.go
  outFile :  ""
  # generate unit test for query code
  withUnitTest  : false
  # generated model code's package name
  modelPkgName  : ""
  # generate with pointer when field is nullable
  fieldNullable : false
  # generate field with gorm index tag
  fieldWithIndexTag : false
  # generate field with gorm column type tag
  fieldWithTypeTag  : false
```