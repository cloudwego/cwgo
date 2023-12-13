namespace go base

struct PingReq{
}
struct PingRes{
    1: i32 code
    2: string msg
}

service BaseService {
    PingRes Ping(1: PingReq req)(api.get="/api/ping")
}