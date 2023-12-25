namespace go template

include "../model/template.thrift"

struct AddTemplateReq{
    1: string name (api.body="name,required",api.vd="len($)>0")
    2: i32 type (api.body="type,required")
}
struct AddTemplateResData{
    1: i64 id
}
struct AddTemplateRes{
    1: i32 code
    2: string msg
    3: AddTemplateResData data
}

struct DeleteTemplateReq{
    1: list<i64> ids (api.body="ids,required",api.vd="len($)>0")
}
struct DeleteTemplateRes{
    1: i32 code
    2: string msg
}

struct UpdateTemplateReq{
    1: i64 id (api.body="id,required")
    2: string name (api.body="name,required",api.vd="len($)>0")
}
struct UpdateTemplateRes{
    1: i32 code
    2: string msg
}

struct AddTemplateItemReq{
    1: i64 template_id
    2: string name (api.body="name,required",api.vd="len($)>0")
    3: string content (api.body="content,required")
}
struct AddTemplateItemResData{
    1: i64 id
}
struct AddTemplateItemRes{
    1: i32 code
    2: string msg
    3: AddTemplateItemResData data
}

struct DeleteTemplateItemReq{
    1: list<i64> ids (api.body="ids,required",api.vd="len($)>0")
}
struct DeleteTemplateItemRes{
    1: i32 code
    2: string msg
}

struct UpdateTemplateItemReq{
    1: i64 id (api.body="id,required")
    2: string name (api.body="name")
    3: string content (api.body="content")
}
struct UpdateTemplateItemRes{
    1: i32 code
    2: string msg
}

struct GetTemplatesReq{
    1: i32 page (api.query="page",api.vd="$>=0")
    2: i32 limit (api.query="limit",api.vd="$>=0")
    3: i32 order (api.query="order",api.vd="$>=0")
    4: string order_by (api.query="order_by")
    5: i32 type (api.query="type")
    6: string name (api.query="name")
}
struct GetTemplatesRes{
    1: i32 code
    2: string msg
    3: GetTemplatesResData data
}
struct GetTemplatesResData{
    1: list<template.TemplateWithInfo> templates
    2: i32 total
}

service TemplateService {
    AddTemplateRes AddTemplate(1: AddTemplateReq req) (api.post="/api/template")
    DeleteTemplateRes DeleteTemplate(1: DeleteTemplateReq req) (api.delete="/api/template")
    UpdateTemplateRes UpdateTemplate(1: UpdateTemplateReq req) (api.patch="/api/template")

    AddTemplateItemRes AddTemplateItem(1: AddTemplateItemReq req) (api.post="/api/template/item")
    DeleteTemplateItemRes DeleteTemplateItem(1: DeleteTemplateItemReq req) (api.delete="/api/template/item")
    UpdateTemplateItemRes UpdateTemplateItem(1: UpdateTemplateItemReq req) (api.patch="/api/template/item")

    GetTemplatesRes GetTemplates(1: GetTemplatesReq req) (api.get="/api/template")
}