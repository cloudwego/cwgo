namespace go idl

struct IDL{
    1: i64 id
    2: i64 repository_id
    3: string main_idl_path
    4: string content
    5: string service_name
    6: string last_sync_time
    7: string create_time
    8: string update_time
}

struct AddIDLReq{
    1: i64 repository_id (api.body="repository_id,required")
    2: string main_idl_path (api.body="main_idl_path,required")
    3: string service_name (api.body="service_name,required")
}
struct AddIDLRes{
    1: i32 code
    2: string msg
}

struct DeleteIDLsReq{
    1: list<i64> ids (api.body="ids,required")
}
struct DeleteIDLsRes{
    1: i32 code
    2: string msg
}

struct UpdateIDLReq{
    1: i64 id (api.body="id,required")
    2: i64 repository_id (api.body="repository_id")
    3: string main_idl_path (api.body="main_idl_path")
    4: string service_name (api.body="service_name")
}
struct UpdateIDLRes{
    1: i32 code
    2: string msg
}

struct GetIDLsReq{
    1: i32 page (api.query="page",api.vd="$>=0")
    2: i32 limit (api.query="limit",api.vd="$>=0")
}
struct GetIDLsRes{
    1: i32 code
    2: string msg
    3: GetIDLsResData data
}
struct GetIDLsResData{
    1: list<IDL> idls
}

struct SyncIDLsByIdReq{
    1: list<i64> ids (api.body="ids,required")
}
struct SyncIDLsByIdRes{
    1: i32 code
    2: string msg
}

service IDLService {
    AddIDLRes AddIDL(1: AddIDLReq req) (api.post="/idl")
    DeleteIDLsRes DeleteIDL(1: DeleteIDLsReq req) (api.delete="/idl")
    UpdateIDLRes UpdateIDL(1: UpdateIDLReq req) (api.patch="/idl")
    GetIDLsRes GetIDLs(1: GetIDLsReq req) (api.get="/idl")

    SyncIDLsByIdRes SyncIDLs(1: SyncIDLsByIdReq req) (api.post="/idl/sync")
}