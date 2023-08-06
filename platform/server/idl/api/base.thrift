namespace go base

include "../base/user.thrift"

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

service ApiService {

    RegisterRes Register(1: RegisterReq req) (api.post="/register")
    LoginRes Login(1: LoginReq req) (api.post="/login");
}