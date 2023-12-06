namespace go base

include "../model/user.thrift"

struct PingReq{
}
struct PingRes{
    1: i32 code
    2: string msg
}

struct RegisterReq{
    1: string username (api.form="username,required")
    2: string password (api.form="password,required")
}
struct RegisterRes{
    1: i32 code
    2: string msg
}

struct LoginReq{
    1: string username (api.form="username,required")
    2: string password (api.form="password,required")
}
struct LoginRes{
    1: i32 code;
    2: string msg;
    3: LoginResData data;
}
struct LoginResData{
    1: user.UserInfo user_info;
}

service BaseService {
    PingRes Ping(1: PingReq req)(api.get="/api/ping")

    RegisterRes Register(1: RegisterReq req) (api.post="/api/register")
    LoginRes Login(1: LoginReq req) (api.post="/api/login");
}