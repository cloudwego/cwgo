namespace go registry

struct RegisterReq{
    1: string service_id (api.query="service_id,required")
    2: string host (api.query="host,required")
    3: i32 port (api.query="port,required")
}
struct RegisterRes{
    1: i32 code
    2: string msg
}

struct DeregisterReq{
    1: string service_id (api.query="service_id,required")
}
struct DeRegisterRes{
    1: i32 code
    2: string msg
}

struct UpdateReq{
    1: string service_id (api.query="service_id,required")
}
struct UpdateRes{
    1: i32 code
    2: string msg
}

service RegistryService {
    RegisterRes Register(1: RegisterReq req) (api.get="/registry/register")
    DeRegisterRes Deregister(1: DeregisterReq req) (api.get="/registry/dnregister")
    UpdateRes Update(1: UpdateReq req) (api.get="/registry/update")
}