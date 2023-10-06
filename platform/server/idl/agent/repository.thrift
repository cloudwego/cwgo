namespace go agent

struct AddRepositoryReq{
    1: i32 repository_type
    2: string repository_url
    3: string token
}
struct AddRepositoryRes{
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