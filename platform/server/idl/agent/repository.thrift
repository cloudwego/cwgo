namespace go agent

include "../model/repository.thrift"

struct AddRepositoryReq{
    1: i32 repository_type
    2: string repository_url
    3: string token
    4: i32 store_type
}
struct AddRepositoryRes{
    1: i32 code
    2: string msg
}

struct DeleteRepositoriesReq{
    1: list<string> ids
}
struct DeleteRepositoriesRes{
    1: i32 code
    2: string msg
}

struct UpdateRepositoryReq{
    1: i64 id
    2: string token (api.body="token")
    3: i32 status (api.body="status")
}
struct UpdateRepositoryRes{
    1: i32 code
    2: string msg
}

struct GetRepositoriesReq{
    1: i32 page
    2: i32 limit
    3: i32 order
    4: string order_by
}
struct GetRepositoriesRes{
    1: i32 code
    2: string msg
    3: GetRepositoriesResData data
}
struct GetRepositoriesResData{
    1: list<repository.Repository> repositories
}

struct SyncRepositoryByIdReq{
    1: list<i64> ids
}
struct SyncRepositoryByIdRes{
    1: i32 code
    2: string msg
}
