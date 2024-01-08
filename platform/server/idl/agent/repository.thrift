namespace go agent

include "../model/repository.thrift"

struct AddRepositoryReq{
    1: i32 repository_type
    2: string repository_domain
    3: string repository_owner
    4: string repository_name
    5: string branch
    6: i32 store_type
}
struct AddRepositoryResp{
    1: i32 code
    2: string msg
}

struct DeleteRepositoriesReq{
    1: list<i64> ids
}
struct DeleteRepositoriesResp{
    1: i32 code
    2: string msg
}

struct UpdateRepositoryReq{
    1: i64 id
    2: string branch
    3: i32 status
}
struct UpdateRepositoryResp{
    1: i32 code
    2: string msg
}

struct GetRepositoriesReq{
    1: i32 page
    2: i32 limit
    3: i32 order
    4: string order_by
    5: i32 repository_type
    6: i32 store_type
    7: string repository_domain
    8: string repository_owner
    9: string repository_name
}
struct GetRepositoriesResp{
    1: i32 code
    2: string msg
    3: GetRepositoriesRespData data
}
struct GetRepositoriesRespData{
    1: list<repository.Repository> repositories
    2: i32 total
}
