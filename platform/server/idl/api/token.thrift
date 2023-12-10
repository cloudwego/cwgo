namespace go token

include "../model/token.thrift"

struct AddTokenReq{
    1: i32 repository_type (api.body="repository_type,required")
    2: string repository_domain (api.body="repository_domain,required")
    3: string token (api.body="token,required")
}
struct AddTokenResData{
    1: string owner
    2: string expiration_time
}
struct AddTokenRes{
    1: i32 code
    2: string msg
    3: AddTokenResData data
}

struct DeleteTokenReq{
    1: list<i64> ids (api.body="ids,required",api.vd="len($)>0")
}
struct DeleteTokenRes{
    1: i32 code
    2: string msg
}

struct GetTokenReq{
    1: i32 page (api.query="page",api.vd="$>=0")
    2: i32 limit (api.query="limit",api.vd="$>=0")
    3: i32 order (api.query="order",api.vd="$>=0")
    4: string order_by (api.query="order_by")
    5: i32 repository_type (api.query="repository_type")
    6: string repository_domain (api.query="repository_domain")
    7: string owner (api.query="owner")
}
struct GetTokenResData{
    1: list<token.Token> tokens
    2: i32 total
}
struct GetTokenRes{
    1: i32 code
    2: string msg
    3: GetTokenResData data
}

service TokenService {
    AddTokenRes AddToken(1: AddTokenReq req) (api.post="/api/token")
    DeleteTokenRes DeleteToken(1: DeleteTokenReq req) (api.delete="/api/token")
    GetTokenRes GetToken(1: GetTokenReq req) (api.get="/api/token")
}