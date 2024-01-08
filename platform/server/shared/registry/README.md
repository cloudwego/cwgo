## 说明

### 接口

```go
type IRegistry interface {
    registry.Registry
    GetAgents() ([]*meta.Agent, error)
}
```

`IRegistry` 在 kitex 的 `registry.Registry` 的基础上扩展了一个 `GetAgents` 方法，用于获取所有的 agent 信息。

### 注册和注销的实现原理

```go
go m.subscribe(ctx, info, r)

for msg := range ch {
    split := strings.Split(msg.Payload, "/")
    registerType, srvName, addr := split[0], split[1], split[2]
    
    switch registerType {
    case register:
        if m.mentorMap[srvName] == nil {
            m.mentorMap[srvName] = mapset.NewSet[string]()
        }
        m.mentorMap[srvName].Add(addr)
    case deregister:
        m.mentorMap[srvName].Remove(addr)
    default:
    }
}
```
- 订阅了一个 redis 的 channel，通过 channel 接收到 agent 的信息, 来确认是 register 还是 deregister.
然后进行相应的操作

```go
hashTable, err := prepareRegistryMeta(info)
args := []any{
	hashTable.field, hashTable.value, 60,
	generateMsg(register, info.ServiceName, info.Addr.String()),
}
// HSET cwgo:ping "127.0.0.1:8081" registryHashMap
// PUBLISH cwgo:ping "register/ping/127.0.0.1:8081"
// EXPIRE cwgo:ping 60 
err = registerScript.Run(ctx, r.rdb, []string{hashTable.key}, args).Err()
```

- 将 agent 的信息存储到 redis 中，key 为 `cwgo:ping`，field 为 agent 的地址，value 为 agent 的信息，过期时间为 60s.
并且将 agent 的信息通过 redis 的 channel 发布出去

```go
// HDEL cwgo:ping 
// PUBLISH cwgo:ping "deregister/ping/127.0.0.1:8081
err = deregisterScript.Run(ctx, r.rdb, []string{hashTable.key}, args).Err()
```
- 将 agent 的信息从 redis 中删除，并且将 agent 的信息通过 redis 的 channel 发布出去


### resolver agent 信息的实现原理

```go
fvs := r.rdb.HGetAll(ctx, desc).Val()
for _, hashTable := range fvs {
	var rInfo registryInfo
	err := sonic.Unmarshal([]byte(hashTable), &rInfo)
	instances = append(instances, discovery.NewInstance("tcp", rInfo.Addr, weight, rInfo.Tags))
}
```

- 通过 `HGetAll` 获取到所有的 agent 信息，然后将信息转换成 `discovery.Instance` 的格式，返回给调用方
