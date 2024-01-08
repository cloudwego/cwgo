namespace go agent

include "../model/token.thrift"

struct AddTokenReq{
    1: i32 repository_type
    2: string repository_domain
    3: string token
}
struct AddTokenResData{
    1: string owner
    2: string expiration_time
}
struct AddTokenResp{
    1: i32 code
    2: string msg
    3: AddTokenResData data
}

struct DeleteTokenReq{
    1: list<i64> ids
}
struct DeleteTokenResp{
    1: i32 code
    2: string msg
}

struct GetTokenReq{
    1: i32 page
    2: i32 limit
    3: i32 order
    4: string order_by
    5: string owner
    6: i32 repository_type
    7: string repository_domain
}
struct GetTokenResData{
    1: list<token.Token> tokens
    2: i32 total
}
struct GetTokenResp{
    1: i32 code
    2: string msg
    3: GetTokenResData data
}
