namespace go agent

struct GenerateCodeReq{
    1: i64 idl_id
}
struct GenerateCodeRes{
    1: i32 code
    2: string msg
}

struct SyncRepositoryByIdReq{
    1: list<i64> ids
}
struct SyncRepositoryByIdRes{
    1: i32 code
    2: string msg
}

struct UpdateRepositoryStatusReq{
    1: i64 id
    2: string status
}
struct UpdateRepositoryStatusRes{
    1: i32 code
    2: string msg
}

struct SyncIDLsByIdReq{
    1: list<i64> ids (api.body="ids,required")
}
struct SyncIDLsByIdRes{
    1: i32 code
    2: string msg
}

service AgentService {
    GenerateCodeRes GenerateCode(1: GenerateCodeReq req)
    SyncRepositoryByIdRes SyncRepositoryById(1: SyncRepositoryByIdReq req)
    UpdateRepositoryStatusRes UpdateRepositoryStatus(1: UpdateRepositoryStatusReq req)
    SyncIDLsByIdRes SyncIDLsById(1: SyncIDLsByIdReq req)
}