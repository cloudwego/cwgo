namespace go agent

include "../model/idl.thrift"

struct AddIDLReq{
    1: i64 repository_id
    2: string main_idl_path
    3: string service_name
    4: string service_repository_name
}
struct AddIDLResp{
    1: i32 code
    2: string msg
    3: AddIDLRespData data
}
struct AddIDLRespData{
    1: i64 idl_id
}

struct DeleteIDLsReq{
    1: list<i64> ids
}
struct DeleteIDLsResp{
    1: i32 code
    2: string msg
}

struct UpdateIDLReq{
    1: i64 id
    2: i64 repository_id
    3: string main_idl_path
    4: i32 status
    5: string service_name
}
struct UpdateIDLResp{
    1: i32 code
    2: string msg
}

struct GetIDLsReq{
    1: i32 page
    2: i32 limit
    3: i32 order
    4: string order_by
    5: string service_name
}
struct GetIDLsResp{
    1: i32 code
    2: string msg
    3: GetIDLsRespData data
}
struct GetIDLsRespData{
    1: list<idl.IDLWithRepositoryInfo> idls
    2: i32 total
}

struct SyncIDLsByIdReq{
    1: list<i64> ids
}

struct SyncIDLsByIdResp{
    1: i32 code
    2: string msg
}
