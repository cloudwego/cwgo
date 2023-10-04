namespace go registry

struct RegisterReq{
    1: string service_id (api.form="service_id,required")
}
struct RegisterRes{
    1: i32 code
    2: string msg
}

struct UnregisterReq{
    1: string service_id (api.form="service_id,required")
}
struct UnRegisterRes{
    1: i32 code
    2: string msg
}

struct UpdateReq{
    1: string service_id (api.form="service_id,required")
}
struct UpdateRes{
    1: i32 code
    2: string msg
}

service RegistryService {
    RegisterRes Register(1: RegisterReq req) (api.get="/registry/register")
    UnRegisterRes Unregister(1: UnregisterReq req) (api.get="/registry/unregister")
    UpdateRes Update(1: UpdateReq req) (api.get="/registry/update")
}