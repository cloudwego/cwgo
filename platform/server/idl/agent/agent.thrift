namespace go agent

struct GenerateCodeReq{
    1: i64 idl_id
}
struct GenerateCodeRes{
    1: i32 code
    2: string msg
}

service AgentService {
    GenerateCodeRes GenerateCode(1: GenerateCodeReq req) (api.post="/")
}