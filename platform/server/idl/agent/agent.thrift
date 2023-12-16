namespace go agent

include "../model/task.thrift"

struct UpdateTasksReq{
    1: list<task.Task> tasks
}
struct UpdateTasksRes{
    1: i32 code
    2: string msg
}
