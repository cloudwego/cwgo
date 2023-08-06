namespace go repository

struct Repository{
    1: string id
    2: string repository_url
    3: string last_update_time
    4: string last_sync_time
    5: string token
    6: string status
}

struct AddRepositoryReq{
    1: string repository_url (api.body="repository_url,required")
    2: string token (api.body="token")
}
struct AddRepositroyRes{
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
    1: string token (api.body="token")
}
struct UpdateRepositoryRes{
    1: i32 code
    2: string msg
}

struct GetRepositoriesReq{
    1: i32 page (api.query="page",api.vd="$>=0")
    2: i32 limit (api.query="limit",api.vd="$>=0")
}
struct GetRepositoriesRes{
    1: i32 code
    2: string msg
    3: GetRepositoriesResData data
}
struct GetRepositoriesResData{
    1: list<Repository> repositories
}

struct SyncRepositoryByIdReq{
    1: list<i64> ids (api.body="ids,required")
}
struct SyncRepositoryByIdRes{
    1: i32 code
    2: string msg
}

service IDLService {
    AddRepositroyRes AddIDL(1: AddRepositoryReq req) (api.post="/repo")
    DeleteRepositoriesRes DeleteIDL(1: DeleteRepositoriesReq req) (api.delete="/repo")
    UpdateRepositoryRes UpdateIDL(1: UpdateRepositoryReq req) (api.patch="/repo")
    GetRepositoriesRes GetIDL(1: GetRepositoriesReq req) (api.get="/repo")

    SyncRepositoryByIdRes SyncIDL(1: SyncRepositoryByIdReq req) (api.post="/repo/sync")
}