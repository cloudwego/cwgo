namespace go agent

include "../model/model.thrift"

struct UpdateTasksReq{
    1: list<model.Task> tasks
}
struct UpdateTasksRes{
    1: i32 code
    2: string msg
}

struct GenerateCodeReq{
    1: i64 idl_id
}
struct GenerateCodeRes{
    1: i32 code
    2: string msg
}
