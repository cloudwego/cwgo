namespace go agent

include "agent.thrift"
include "repository.thrift"
include "idl.thrift"
include "template.thrift"
include "token.thrift"

service AgentService {
    repository.AddRepositoryRes AddRepository(1: repository.AddRepositoryReq req)
    repository.DeleteRepositoriesRes DeleteRepositories(1: repository.DeleteRepositoriesReq req)
    repository.UpdateRepositoryRes UpdateRepository(1: repository.UpdateRepositoryReq req)
    repository.GetRepositoriesRes GetRepositories(1: repository.GetRepositoriesReq req)
    repository.SyncRepositoryByIdRes SyncRepositoryById(1: repository.SyncRepositoryByIdReq req)

    idl.AddIDLRes AddIDL(1: idl.AddIDLReq req)
    idl.DeleteIDLsRes DeleteIDL(1: idl.DeleteIDLsReq req)
    idl.UpdateIDLRes UpdateIDL(1: idl.UpdateIDLReq req)
    idl.GetIDLsRes GetIDLs(1: idl.GetIDLsReq req)
    idl.SyncIDLsByIdRes SyncIDLsById(1: idl.SyncIDLsByIdReq req)

    template.AddTemplateRes AddTemplate(1: template.AddTemplateReq req)
    template.DeleteTemplateRes DeleteTemplate(1: template.DeleteTemplateReq req)
    template.UpdateTemplateRes UpdateTemplate(1: template.UpdateTemplateReq req)
    template.GetTemplatesRes GetTemplates(1: template.GetTemplatesReq req)
    template.AddTemplateItemRes AddTemplateItem(1: template.AddTemplateItemReq req)
    template.DeleteTemplateItemRes DeleteTemplateItem(1: template.DeleteTemplateItemReq req)
    template.UpdateTemplateItemRes UpdateTemplateItem(1: template.UpdateTemplateItemReq req)
    template.GetTemplateItemsRes GetTemplateItems(1: template.GetTemplateItemsReq req)

    agent.UpdateTasksRes UpdateTasks(1: agent.UpdateTasksReq req)
    agent.GenerateCodeRes GenerateCode(1: agent.GenerateCodeReq req)

    token.AddTokenRes AddToken(1: token.AddTokenReq req)
    token.DeleteTokenRes DeleteToken(1: token.DeleteTokenReq req)
    token.GetTokenRes GetToken(1: token.GetTokenReq req)
}