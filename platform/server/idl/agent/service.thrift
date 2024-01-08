namespace go agent

include "repository.thrift"
include "idl.thrift"
include "template.thrift"
include "token.thrift"
include "task.thrift"

struct PingReq {
    1: string msg
}

struct PingResp {
    1: string msg
}

service AgentService {
    PingResp Ping(1: PingReq req)

    repository.AddRepositoryResp AddRepository(1: repository.AddRepositoryReq req)
    repository.DeleteRepositoriesResp DeleteRepositories(1: repository.DeleteRepositoriesReq req)
    repository.UpdateRepositoryResp UpdateRepository(1: repository.UpdateRepositoryReq req)
    repository.GetRepositoriesResp GetRepositories(1: repository.GetRepositoriesReq req)

    idl.AddIDLResp AddIDL(1: idl.AddIDLReq req)
    idl.DeleteIDLsResp DeleteIDL(1: idl.DeleteIDLsReq req)
    idl.UpdateIDLResp UpdateIDL(1: idl.UpdateIDLReq req)
    idl.GetIDLsResp GetIDLs(1: idl.GetIDLsReq req)
    idl.SyncIDLsByIdResp SyncIDLsById(1: idl.SyncIDLsByIdReq req)

    template.AddTemplateResp AddTemplate(1: template.AddTemplateReq req)
    template.DeleteTemplateResp DeleteTemplate(1: template.DeleteTemplateReq req)
    template.UpdateTemplateResp UpdateTemplate(1: template.UpdateTemplateReq req)
    template.GetTemplatesResp GetTemplates(1: template.GetTemplatesReq req)
    template.AddTemplateItemResp AddTemplateItem(1: template.AddTemplateItemReq req)
    template.DeleteTemplateItemResp DeleteTemplateItem(1: template.DeleteTemplateItemReq req)
    template.UpdateTemplateItemResp UpdateTemplateItem(1: template.UpdateTemplateItemReq req)
    template.GetTemplateItemsResp GetTemplateItems(1: template.GetTemplateItemsReq req)

    task.UpdateTaskResp UpdateTask(1: task.UpdateTaskReq req)

    token.AddTokenResp AddToken(1: token.AddTokenReq req)
    token.DeleteTokenResp DeleteToken(1: token.DeleteTokenReq req)
    token.GetTokenResp GetToken(1: token.GetTokenReq req)
}
