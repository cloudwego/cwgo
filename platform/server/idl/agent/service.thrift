namespace go agent

include "agent.thrift"
include "repository.thrift"
include "idl.thrift"

service AgentService {
    repository.AddRepositoryRes AddRepository(1: repository.AddRepositoryReq req)
    repository.UpdateRepositoryStatusRes UpdateRepositoryStatus(1: repository.UpdateRepositoryStatusReq req)
    repository.SyncRepositoryByIdRes SyncRepositoryById(1: repository.SyncRepositoryByIdReq req)

    idl.AddIDLRes AddIDL(1: idl.AddIDLReq req)
    idl.SyncIDLsByIdRes SyncIDLsById(1: idl.SyncIDLsByIdReq req)

    agent.UpdateTasksRes UpdateTasks(1: agent.UpdateTasksReq req)

    agent.GenerateCodeRes GenerateCode(1: agent.GenerateCodeReq req)
}