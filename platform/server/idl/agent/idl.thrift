namespace go agent

struct AddIDLReq{
    1: i64 repository_id
    2: string main_idl_path
    3: string service_name
}
struct AddIDLRes{
    1: i32 code
    2: string msg
}

struct SyncIDLsByIdReq{
    1: list<i64> ids
}
struct SyncIDLsByIdRes{
    1: i32 code
    2: string msg
}