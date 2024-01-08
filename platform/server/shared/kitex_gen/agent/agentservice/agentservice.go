// Code generated by Kitex v0.8.0. DO NOT EDIT.

package agentservice

import (
	"context"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	task "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/task"
	frugal "github.com/cloudwego/frugal"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	"reflect"
	"sync"
)

func serviceInfo() *kitex.ServiceInfo {
	return agentServiceServiceInfo
}

var agentServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "AgentService"
	handlerType := (*agent.AgentService)(nil)
	methods := map[string]kitex.MethodInfo{
		"Ping":               kitex.NewMethodInfo(pingHandler, newAgentServicePingArgs, newAgentServicePingResult, false),
		"AddRepository":      kitex.NewMethodInfo(addRepositoryHandler, newAgentServiceAddRepositoryArgs, newAgentServiceAddRepositoryResult, false),
		"DeleteRepositories": kitex.NewMethodInfo(deleteRepositoriesHandler, newAgentServiceDeleteRepositoriesArgs, newAgentServiceDeleteRepositoriesResult, false),
		"UpdateRepository":   kitex.NewMethodInfo(updateRepositoryHandler, newAgentServiceUpdateRepositoryArgs, newAgentServiceUpdateRepositoryResult, false),
		"GetRepositories":    kitex.NewMethodInfo(getRepositoriesHandler, newAgentServiceGetRepositoriesArgs, newAgentServiceGetRepositoriesResult, false),
		"AddIDL":             kitex.NewMethodInfo(addIDLHandler, newAgentServiceAddIDLArgs, newAgentServiceAddIDLResult, false),
		"DeleteIDL":          kitex.NewMethodInfo(deleteIDLHandler, newAgentServiceDeleteIDLArgs, newAgentServiceDeleteIDLResult, false),
		"UpdateIDL":          kitex.NewMethodInfo(updateIDLHandler, newAgentServiceUpdateIDLArgs, newAgentServiceUpdateIDLResult, false),
		"GetIDLs":            kitex.NewMethodInfo(getIDLsHandler, newAgentServiceGetIDLsArgs, newAgentServiceGetIDLsResult, false),
		"SyncIDLsById":       kitex.NewMethodInfo(syncIDLsByIdHandler, newAgentServiceSyncIDLsByIdArgs, newAgentServiceSyncIDLsByIdResult, false),
		"AddTemplate":        kitex.NewMethodInfo(addTemplateHandler, newAgentServiceAddTemplateArgs, newAgentServiceAddTemplateResult, false),
		"DeleteTemplate":     kitex.NewMethodInfo(deleteTemplateHandler, newAgentServiceDeleteTemplateArgs, newAgentServiceDeleteTemplateResult, false),
		"UpdateTemplate":     kitex.NewMethodInfo(updateTemplateHandler, newAgentServiceUpdateTemplateArgs, newAgentServiceUpdateTemplateResult, false),
		"GetTemplates":       kitex.NewMethodInfo(getTemplatesHandler, newAgentServiceGetTemplatesArgs, newAgentServiceGetTemplatesResult, false),
		"AddTemplateItem":    kitex.NewMethodInfo(addTemplateItemHandler, newAgentServiceAddTemplateItemArgs, newAgentServiceAddTemplateItemResult, false),
		"DeleteTemplateItem": kitex.NewMethodInfo(deleteTemplateItemHandler, newAgentServiceDeleteTemplateItemArgs, newAgentServiceDeleteTemplateItemResult, false),
		"UpdateTemplateItem": kitex.NewMethodInfo(updateTemplateItemHandler, newAgentServiceUpdateTemplateItemArgs, newAgentServiceUpdateTemplateItemResult, false),
		"GetTemplateItems":   kitex.NewMethodInfo(getTemplateItemsHandler, newAgentServiceGetTemplateItemsArgs, newAgentServiceGetTemplateItemsResult, false),
		"UpdateTask":         kitex.NewMethodInfo(updateTaskHandler, newAgentServiceUpdateTaskArgs, newAgentServiceUpdateTaskResult, false),
		"AddToken":           kitex.NewMethodInfo(addTokenHandler, newAgentServiceAddTokenArgs, newAgentServiceAddTokenResult, false),
		"DeleteToken":        kitex.NewMethodInfo(deleteTokenHandler, newAgentServiceDeleteTokenArgs, newAgentServiceDeleteTokenResult, false),
		"GetToken":           kitex.NewMethodInfo(getTokenHandler, newAgentServiceGetTokenArgs, newAgentServiceGetTokenResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "agent",
		"ServiceFilePath": `idl/agent/service.thrift`,
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.8.0",
		Extra:           extra,
	}
	return svcInfo
}

func pingHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServicePingArgs)
	realResult := result.(*agent.AgentServicePingResult)
	success, err := handler.(agent.AgentService).Ping(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServicePingArgs() interface{} {
	return agent.NewAgentServicePingArgs()
}

func newAgentServicePingResult() interface{} {
	return agent.NewAgentServicePingResult()
}

func addRepositoryHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceAddRepositoryArgs)
	realResult := result.(*agent.AgentServiceAddRepositoryResult)
	success, err := handler.(agent.AgentService).AddRepository(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceAddRepositoryArgs() interface{} {
	return agent.NewAgentServiceAddRepositoryArgs()
}

func newAgentServiceAddRepositoryResult() interface{} {
	return agent.NewAgentServiceAddRepositoryResult()
}

func deleteRepositoriesHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceDeleteRepositoriesArgs)
	realResult := result.(*agent.AgentServiceDeleteRepositoriesResult)
	success, err := handler.(agent.AgentService).DeleteRepositories(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceDeleteRepositoriesArgs() interface{} {
	return agent.NewAgentServiceDeleteRepositoriesArgs()
}

func newAgentServiceDeleteRepositoriesResult() interface{} {
	return agent.NewAgentServiceDeleteRepositoriesResult()
}

func updateRepositoryHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceUpdateRepositoryArgs)
	realResult := result.(*agent.AgentServiceUpdateRepositoryResult)
	success, err := handler.(agent.AgentService).UpdateRepository(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceUpdateRepositoryArgs() interface{} {
	return agent.NewAgentServiceUpdateRepositoryArgs()
}

func newAgentServiceUpdateRepositoryResult() interface{} {
	return agent.NewAgentServiceUpdateRepositoryResult()
}

func getRepositoriesHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceGetRepositoriesArgs)
	realResult := result.(*agent.AgentServiceGetRepositoriesResult)
	success, err := handler.(agent.AgentService).GetRepositories(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceGetRepositoriesArgs() interface{} {
	return agent.NewAgentServiceGetRepositoriesArgs()
}

func newAgentServiceGetRepositoriesResult() interface{} {
	return agent.NewAgentServiceGetRepositoriesResult()
}

func addIDLHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceAddIDLArgs)
	realResult := result.(*agent.AgentServiceAddIDLResult)
	success, err := handler.(agent.AgentService).AddIDL(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceAddIDLArgs() interface{} {
	return agent.NewAgentServiceAddIDLArgs()
}

func newAgentServiceAddIDLResult() interface{} {
	return agent.NewAgentServiceAddIDLResult()
}

func deleteIDLHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceDeleteIDLArgs)
	realResult := result.(*agent.AgentServiceDeleteIDLResult)
	success, err := handler.(agent.AgentService).DeleteIDL(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceDeleteIDLArgs() interface{} {
	return agent.NewAgentServiceDeleteIDLArgs()
}

func newAgentServiceDeleteIDLResult() interface{} {
	return agent.NewAgentServiceDeleteIDLResult()
}

func updateIDLHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceUpdateIDLArgs)
	realResult := result.(*agent.AgentServiceUpdateIDLResult)
	success, err := handler.(agent.AgentService).UpdateIDL(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceUpdateIDLArgs() interface{} {
	return agent.NewAgentServiceUpdateIDLArgs()
}

func newAgentServiceUpdateIDLResult() interface{} {
	return agent.NewAgentServiceUpdateIDLResult()
}

func getIDLsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceGetIDLsArgs)
	realResult := result.(*agent.AgentServiceGetIDLsResult)
	success, err := handler.(agent.AgentService).GetIDLs(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceGetIDLsArgs() interface{} {
	return agent.NewAgentServiceGetIDLsArgs()
}

func newAgentServiceGetIDLsResult() interface{} {
	return agent.NewAgentServiceGetIDLsResult()
}

func syncIDLsByIdHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceSyncIDLsByIdArgs)
	realResult := result.(*agent.AgentServiceSyncIDLsByIdResult)
	success, err := handler.(agent.AgentService).SyncIDLsById(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceSyncIDLsByIdArgs() interface{} {
	return agent.NewAgentServiceSyncIDLsByIdArgs()
}

func newAgentServiceSyncIDLsByIdResult() interface{} {
	return agent.NewAgentServiceSyncIDLsByIdResult()
}

func addTemplateHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceAddTemplateArgs)
	realResult := result.(*agent.AgentServiceAddTemplateResult)
	success, err := handler.(agent.AgentService).AddTemplate(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceAddTemplateArgs() interface{} {
	return agent.NewAgentServiceAddTemplateArgs()
}

func newAgentServiceAddTemplateResult() interface{} {
	return agent.NewAgentServiceAddTemplateResult()
}

func deleteTemplateHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceDeleteTemplateArgs)
	realResult := result.(*agent.AgentServiceDeleteTemplateResult)
	success, err := handler.(agent.AgentService).DeleteTemplate(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceDeleteTemplateArgs() interface{} {
	return agent.NewAgentServiceDeleteTemplateArgs()
}

func newAgentServiceDeleteTemplateResult() interface{} {
	return agent.NewAgentServiceDeleteTemplateResult()
}

func updateTemplateHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceUpdateTemplateArgs)
	realResult := result.(*agent.AgentServiceUpdateTemplateResult)
	success, err := handler.(agent.AgentService).UpdateTemplate(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceUpdateTemplateArgs() interface{} {
	return agent.NewAgentServiceUpdateTemplateArgs()
}

func newAgentServiceUpdateTemplateResult() interface{} {
	return agent.NewAgentServiceUpdateTemplateResult()
}

func getTemplatesHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceGetTemplatesArgs)
	realResult := result.(*agent.AgentServiceGetTemplatesResult)
	success, err := handler.(agent.AgentService).GetTemplates(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceGetTemplatesArgs() interface{} {
	return agent.NewAgentServiceGetTemplatesArgs()
}

func newAgentServiceGetTemplatesResult() interface{} {
	return agent.NewAgentServiceGetTemplatesResult()
}

func addTemplateItemHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceAddTemplateItemArgs)
	realResult := result.(*agent.AgentServiceAddTemplateItemResult)
	success, err := handler.(agent.AgentService).AddTemplateItem(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceAddTemplateItemArgs() interface{} {
	return agent.NewAgentServiceAddTemplateItemArgs()
}

func newAgentServiceAddTemplateItemResult() interface{} {
	return agent.NewAgentServiceAddTemplateItemResult()
}

func deleteTemplateItemHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceDeleteTemplateItemArgs)
	realResult := result.(*agent.AgentServiceDeleteTemplateItemResult)
	success, err := handler.(agent.AgentService).DeleteTemplateItem(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceDeleteTemplateItemArgs() interface{} {
	return agent.NewAgentServiceDeleteTemplateItemArgs()
}

func newAgentServiceDeleteTemplateItemResult() interface{} {
	return agent.NewAgentServiceDeleteTemplateItemResult()
}

func updateTemplateItemHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceUpdateTemplateItemArgs)
	realResult := result.(*agent.AgentServiceUpdateTemplateItemResult)
	success, err := handler.(agent.AgentService).UpdateTemplateItem(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceUpdateTemplateItemArgs() interface{} {
	return agent.NewAgentServiceUpdateTemplateItemArgs()
}

func newAgentServiceUpdateTemplateItemResult() interface{} {
	return agent.NewAgentServiceUpdateTemplateItemResult()
}

func getTemplateItemsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceGetTemplateItemsArgs)
	realResult := result.(*agent.AgentServiceGetTemplateItemsResult)
	success, err := handler.(agent.AgentService).GetTemplateItems(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceGetTemplateItemsArgs() interface{} {
	return agent.NewAgentServiceGetTemplateItemsArgs()
}

func newAgentServiceGetTemplateItemsResult() interface{} {
	return agent.NewAgentServiceGetTemplateItemsResult()
}

func updateTaskHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceUpdateTaskArgs)
	realResult := result.(*agent.AgentServiceUpdateTaskResult)
	success, err := handler.(agent.AgentService).UpdateTask(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceUpdateTaskArgs() interface{} {
	return agent.NewAgentServiceUpdateTaskArgs()
}

func newAgentServiceUpdateTaskResult() interface{} {
	return agent.NewAgentServiceUpdateTaskResult()
}

func addTokenHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceAddTokenArgs)
	realResult := result.(*agent.AgentServiceAddTokenResult)
	success, err := handler.(agent.AgentService).AddToken(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceAddTokenArgs() interface{} {
	return agent.NewAgentServiceAddTokenArgs()
}

func newAgentServiceAddTokenResult() interface{} {
	return agent.NewAgentServiceAddTokenResult()
}

func deleteTokenHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceDeleteTokenArgs)
	realResult := result.(*agent.AgentServiceDeleteTokenResult)
	success, err := handler.(agent.AgentService).DeleteToken(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceDeleteTokenArgs() interface{} {
	return agent.NewAgentServiceDeleteTokenArgs()
}

func newAgentServiceDeleteTokenResult() interface{} {
	return agent.NewAgentServiceDeleteTokenResult()
}

func getTokenHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*agent.AgentServiceGetTokenArgs)
	realResult := result.(*agent.AgentServiceGetTokenResult)
	success, err := handler.(agent.AgentService).GetToken(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAgentServiceGetTokenArgs() interface{} {
	return agent.NewAgentServiceGetTokenArgs()
}

func newAgentServiceGetTokenResult() interface{} {
	return agent.NewAgentServiceGetTokenResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Ping(ctx context.Context, req *agent.PingReq) (r *agent.PingResp, err error) {
	var _args agent.AgentServicePingArgs
	_args.Req = req
	var _result agent.AgentServicePingResult
	if err = p.c.Call(ctx, "Ping", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AddRepository(ctx context.Context, req *agent.AddRepositoryReq) (r *agent.AddRepositoryResp, err error) {
	var _args agent.AgentServiceAddRepositoryArgs
	_args.Req = req
	var _result agent.AgentServiceAddRepositoryResult
	if err = p.c.Call(ctx, "AddRepository", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteRepositories(ctx context.Context, req *agent.DeleteRepositoriesReq) (r *agent.DeleteRepositoriesResp, err error) {
	var _args agent.AgentServiceDeleteRepositoriesArgs
	_args.Req = req
	var _result agent.AgentServiceDeleteRepositoriesResult
	if err = p.c.Call(ctx, "DeleteRepositories", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateRepository(ctx context.Context, req *agent.UpdateRepositoryReq) (r *agent.UpdateRepositoryResp, err error) {
	var _args agent.AgentServiceUpdateRepositoryArgs
	_args.Req = req
	var _result agent.AgentServiceUpdateRepositoryResult
	if err = p.c.Call(ctx, "UpdateRepository", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetRepositories(ctx context.Context, req *agent.GetRepositoriesReq) (r *agent.GetRepositoriesResp, err error) {
	var _args agent.AgentServiceGetRepositoriesArgs
	_args.Req = req
	var _result agent.AgentServiceGetRepositoriesResult
	if err = p.c.Call(ctx, "GetRepositories", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AddIDL(ctx context.Context, req *agent.AddIDLReq) (r *agent.AddIDLResp, err error) {
	var _args agent.AgentServiceAddIDLArgs
	_args.Req = req
	var _result agent.AgentServiceAddIDLResult
	if err = p.c.Call(ctx, "AddIDL", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteIDL(ctx context.Context, req *agent.DeleteIDLsReq) (r *agent.DeleteIDLsResp, err error) {
	var _args agent.AgentServiceDeleteIDLArgs
	_args.Req = req
	var _result agent.AgentServiceDeleteIDLResult
	if err = p.c.Call(ctx, "DeleteIDL", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateIDL(ctx context.Context, req *agent.UpdateIDLReq) (r *agent.UpdateIDLResp, err error) {
	var _args agent.AgentServiceUpdateIDLArgs
	_args.Req = req
	var _result agent.AgentServiceUpdateIDLResult
	if err = p.c.Call(ctx, "UpdateIDL", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetIDLs(ctx context.Context, req *agent.GetIDLsReq) (r *agent.GetIDLsResp, err error) {
	var _args agent.AgentServiceGetIDLsArgs
	_args.Req = req
	var _result agent.AgentServiceGetIDLsResult
	if err = p.c.Call(ctx, "GetIDLs", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) SyncIDLsById(ctx context.Context, req *agent.SyncIDLsByIdReq) (r *agent.SyncIDLsByIdResp, err error) {
	var _args agent.AgentServiceSyncIDLsByIdArgs
	_args.Req = req
	var _result agent.AgentServiceSyncIDLsByIdResult
	if err = p.c.Call(ctx, "SyncIDLsById", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AddTemplate(ctx context.Context, req *agent.AddTemplateReq) (r *agent.AddTemplateResp, err error) {
	var _args agent.AgentServiceAddTemplateArgs
	_args.Req = req
	var _result agent.AgentServiceAddTemplateResult
	if err = p.c.Call(ctx, "AddTemplate", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteTemplate(ctx context.Context, req *agent.DeleteTemplateReq) (r *agent.DeleteTemplateResp, err error) {
	var _args agent.AgentServiceDeleteTemplateArgs
	_args.Req = req
	var _result agent.AgentServiceDeleteTemplateResult
	if err = p.c.Call(ctx, "DeleteTemplate", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateTemplate(ctx context.Context, req *agent.UpdateTemplateReq) (r *agent.UpdateTemplateResp, err error) {
	var _args agent.AgentServiceUpdateTemplateArgs
	_args.Req = req
	var _result agent.AgentServiceUpdateTemplateResult
	if err = p.c.Call(ctx, "UpdateTemplate", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetTemplates(ctx context.Context, req *agent.GetTemplatesReq) (r *agent.GetTemplatesResp, err error) {
	var _args agent.AgentServiceGetTemplatesArgs
	_args.Req = req
	var _result agent.AgentServiceGetTemplatesResult
	if err = p.c.Call(ctx, "GetTemplates", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AddTemplateItem(ctx context.Context, req *agent.AddTemplateItemReq) (r *agent.AddTemplateItemResp, err error) {
	var _args agent.AgentServiceAddTemplateItemArgs
	_args.Req = req
	var _result agent.AgentServiceAddTemplateItemResult
	if err = p.c.Call(ctx, "AddTemplateItem", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteTemplateItem(ctx context.Context, req *agent.DeleteTemplateItemReq) (r *agent.DeleteTemplateItemResp, err error) {
	var _args agent.AgentServiceDeleteTemplateItemArgs
	_args.Req = req
	var _result agent.AgentServiceDeleteTemplateItemResult
	if err = p.c.Call(ctx, "DeleteTemplateItem", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateTemplateItem(ctx context.Context, req *agent.UpdateTemplateItemReq) (r *agent.UpdateTemplateItemResp, err error) {
	var _args agent.AgentServiceUpdateTemplateItemArgs
	_args.Req = req
	var _result agent.AgentServiceUpdateTemplateItemResult
	if err = p.c.Call(ctx, "UpdateTemplateItem", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetTemplateItems(ctx context.Context, req *agent.GetTemplateItemsReq) (r *agent.GetTemplateItemsResp, err error) {
	var _args agent.AgentServiceGetTemplateItemsArgs
	_args.Req = req
	var _result agent.AgentServiceGetTemplateItemsResult
	if err = p.c.Call(ctx, "GetTemplateItems", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateTask(ctx context.Context, req *task.UpdateTaskReq) (r *task.UpdateTaskResp, err error) {
	var _args agent.AgentServiceUpdateTaskArgs
	_args.Req = req
	var _result agent.AgentServiceUpdateTaskResult
	if err = p.c.Call(ctx, "UpdateTask", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AddToken(ctx context.Context, req *agent.AddTokenReq) (r *agent.AddTokenResp, err error) {
	var _args agent.AgentServiceAddTokenArgs
	_args.Req = req
	var _result agent.AgentServiceAddTokenResult
	if err = p.c.Call(ctx, "AddToken", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteToken(ctx context.Context, req *agent.DeleteTokenReq) (r *agent.DeleteTokenResp, err error) {
	var _args agent.AgentServiceDeleteTokenArgs
	_args.Req = req
	var _result agent.AgentServiceDeleteTokenResult
	if err = p.c.Call(ctx, "DeleteToken", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetToken(ctx context.Context, req *agent.GetTokenReq) (r *agent.GetTokenResp, err error) {
	var _args agent.AgentServiceGetTokenArgs
	_args.Req = req
	var _result agent.AgentServiceGetTokenResult
	if err = p.c.Call(ctx, "GetToken", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

var pretouchOnce sync.Once

func pretouch() {
	pretouchOnce.Do(func() {
		var err error
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServicePingArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServicePingResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceAddRepositoryArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceAddRepositoryResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceDeleteRepositoriesArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceDeleteRepositoriesResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceUpdateRepositoryArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceUpdateRepositoryResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceGetRepositoriesArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceGetRepositoriesResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceAddIDLArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceAddIDLResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceDeleteIDLArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceDeleteIDLResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceUpdateIDLArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceUpdateIDLResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceGetIDLsArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceGetIDLsResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceSyncIDLsByIdArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceSyncIDLsByIdResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceAddTemplateArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceAddTemplateResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceDeleteTemplateArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceDeleteTemplateResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceUpdateTemplateArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceUpdateTemplateResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceGetTemplatesArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceGetTemplatesResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceAddTemplateItemArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceAddTemplateItemResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceDeleteTemplateItemArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceDeleteTemplateItemResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceUpdateTemplateItemArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceUpdateTemplateItemResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceGetTemplateItemsArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceGetTemplateItemsResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceUpdateTaskArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceUpdateTaskResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceAddTokenArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceAddTokenResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceDeleteTokenArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceDeleteTokenResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceGetTokenArgs()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		err = frugal.Pretouch(reflect.TypeOf(agent.NewAgentServiceGetTokenResult()))
		if err != nil {
			goto PRETOUCH_ERR
		}
		return
	PRETOUCH_ERR:
		println("Frugal pretouch in AgentService failed: " + err.Error())
	})
}
