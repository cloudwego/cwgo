namespace go repository

include "../base/model.thrift"

struct AddRepositoryReq{
    1: i32 repository_type (api.body="repository_type,required")
    2: string repository_url (api.body="repository_url,required")
    3: string token (api.body="token")
}
struct AddRepositoryRes{
    1: i32 code
    2: string msg
}

struct DeleteRepositoriesReq{
    1: list<string> ids (api.body="ids,required")
}
struct DeleteRepositoriesRes{
    1: i32 code
    2: string msg
}

struct UpdateRepositoryReq{
    1: string id
    2: string token (api.body="token")
    3: string status (api.body="status")
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
}
struct GetRepositoriesRes{
    1: i32 code
    2: string msg
    3: GetRepositoriesResData data
}
struct GetRepositoriesResData{
    1: list<model.Repository> repositories
}

struct SyncRepositoryByIdReq{
    1: list<i64> ids (api.body="ids,required")
}
struct SyncRepositoryByIdRes{
    1: i32 code
    2: string msg
}

service RepositoryService {
    AddRepositoryRes AddRepository(1: AddRepositoryReq req) (api.post="/repo")
    DeleteRepositoriesRes DeleteRepository(1: DeleteRepositoriesReq req) (api.delete="/repo")
    UpdateRepositoryRes UpdateRepository(1: UpdateRepositoryReq req) (api.patch="/repo")
    GetRepositoriesRes GetRepositories(1: GetRepositoriesReq req) (api.get="/repo")

    SyncRepositoryByIdRes SyncRepository(1: SyncRepositoryByIdReq req) (api.post="/repo/sync")
}