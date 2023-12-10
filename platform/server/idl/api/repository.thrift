namespace go repository

include "../model/repository.thrift"

struct AddRepositoryReq{
    1: i32 repository_type (api.body="repository_type,required")
    2: string repository_url (api.body="repository_url,required",api.vd="len($)>0")
    3: string branch (api.body="branch,required")
    4: i32 store_type (api.body="store_type,required")
}
struct AddRepositoryRes{
    1: i32 code
    2: string msg
}

struct DeleteRepositoriesReq{
    1: list<i64> ids (api.body="ids,required",api.vd="len($)>0")
}
struct DeleteRepositoriesRes{
    1: i32 code
    2: string msg
}

struct UpdateRepositoryReq{
    1: i64 id (api.body="id")
    2: string branch (api.body="branch")
    3: i32 status (api.body="status")
}
struct UpdateRepositoryRes{
    1: i32 code
    2: string msg
}

struct GetRepositoriesReq{
    1: i32 page (api.query="page",api.vd="$>=0")
    2: i32 limit (api.query="limit",api.vd="$>=0")
    3: i32 order (api.query="order",api.vd="$>=0")
    4: string order_by (api.query="order_by")
    5: i32 repository_type (api.query="repository_type")
    6: i32 store_type (api.query="store_type")
    7: string repository_domain (api.query="repository_domain")
    8: string repository_owner (api.query="repository_owner")
    9: string repository_name (api.query="repository_name")
}
struct GetRepositoriesRes{
    1: i32 code
    2: string msg
    3: GetRepositoriesResData data
}
struct GetRepositoriesResData{
    1: list<repository.Repository> repositories
    2: i32 total
}

struct SyncRepositoryByIdReq{
    1: list<i64> ids (api.body="ids,required",api.vd="len($)>0")
}
struct SyncRepositoryByIdRes{
    1: i32 code
    2: string msg
}

service RepositoryService {
    AddRepositoryRes AddRepository(1: AddRepositoryReq req) (api.post="/api/repo")
    DeleteRepositoriesRes DeleteRepository(1: DeleteRepositoriesReq req) (api.delete="/api/repo")
    UpdateRepositoryRes UpdateRepository(1: UpdateRepositoryReq req) (api.patch="/api/repo")
    GetRepositoriesRes GetRepositories(1: GetRepositoriesReq req) (api.get="/api/repo")

    SyncRepositoryByIdRes SyncRepository(1: SyncRepositoryByIdReq req) (api.post="/api/repo/sync")
}