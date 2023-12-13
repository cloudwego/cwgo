namespace go idl

include "../model/idl.thrift"

struct AddIDLReq{
    1: i64 repository_id (api.body="repository_id,required")
    2: string main_idl_path (api.body="main_idl_path,required",api.vd="len($)>0")
    3: string service_name (api.body="service_name,required",api.vd="len($)>0")
    4: string service_repository_name (api.body='service_repository_name')
}
struct AddIDLRes{
    1: i32 code
    2: string msg
}

struct DeleteIDLsReq{
    1: list<i64> ids (api.body="ids,required",api.vd="len($)>0")
}
struct DeleteIDLsRes{
    1: i32 code
    2: string msg
}

struct UpdateIDLReq{
    1: i64 id (api.body="id,required")
    2: i32 status (api.body="status")
    3: string service_name (api.body="service_name")
}
struct UpdateIDLRes{
    1: i32 code
    2: string msg
}

struct GetIDLsReq{
    1: i32 page (api.query="page",api.vd="$>=0")
    2: i32 limit (api.query="limit",api.vd="$>=0")
    3: i32 order (api.query="order",api.vd="$>=0")
    4: string order_by (api.query="order_by")
    5: string service_name (api.query="service_name")
}
struct GetIDLsRes{
    1: i32 code
    2: string msg
    3: GetIDLsResData data
}
struct GetIDLsResData{
    1: list<idl.IDLWithRepositorInfo> idls
    2: i32 total
}

struct SyncIDLsByIdReq{
    1: list<i64> ids (api.body="ids,required",api.vd="len($)>0")
}
struct SyncIDLsByIdRes{
    1: i32 code
    2: string msg
}

service IdlService {
    AddIDLRes AddIDL(1: AddIDLReq req) (api.post="/api/idl")
    DeleteIDLsRes DeleteIDL(1: DeleteIDLsReq req) (api.delete="/api/idl")
    UpdateIDLRes UpdateIDL(1: UpdateIDLReq req) (api.patch="/api/idl")
    GetIDLsRes GetIDLs(1: GetIDLsReq req) (api.get="/api/idl")

    SyncIDLsByIdRes SyncIDLs(1: SyncIDLsByIdReq req) (api.post="/api/idl/sync")
}