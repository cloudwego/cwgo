include "../model/task.thrift"

struct UpdateTaskReq{
    1: list<task.Task> Tasks
}

struct UpdateTaskResp{
    1: i32 Code
    2: string Msg
}
