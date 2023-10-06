namespace go agent

include "agent.thrift"
include "repository.thrift"
include "idl.thrift"

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

    agent.UpdateTasksRes UpdateTasks(1: agent.UpdateTasksReq req)
    agent.GenerateCodeRes GenerateCode(1: agent.GenerateCodeReq req)
}