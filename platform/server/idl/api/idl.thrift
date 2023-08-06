namespace go idl

struct IDL{
    1: i64 repository_id
    2: string main_idl_path
    3: string idl_hash
    4: string service_name
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

struct SyncIDLByIdReq{
    1: list<i64> ids (api.body="ids,required")
}
struct SyncIDLByIdRes{
    1: i32 code
    2: string msg
}

service IDLService {
    AddIDLRes AddIDL(1: AddIDLReq req) (api.post="/idl")
    DeleteIDLsRes DeleteIDL(1: DeleteIDLsReq req) (api.delete="/idl")
    UpdateIDLRes UpdateIDL(1: UpdateIDLReq req) (api.patch="/idl")
    GetIDLsRes GetIDL(1: GetIDLsReq req) (api.get="/idl")

    SyncIDLByIdRes SyncIDL(1: SyncIDLByIdReq req) (api.post="/idl/sync")
}