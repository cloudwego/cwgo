namespace go template

include "../model/template.thrift"

struct AddTemplateReq{
    1: string name (api.body="name,required",api.vd="len($)>0")
    2: i32 type (api.body="type,required")
}
struct AddTemplateRes{
    1: i32 code
    2: string msg
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

struct GetTemplatesReq{
    1: i32 page (api.query="page",api.vd="$>=0")
    2: i32 limit (api.query="limit",api.vd="$>=0")
    3: i32 order (api.query="order",api.vd="$>=0")
    4: string order_by (api.query="order_by")
}
struct GetTemplatesRes{
    1: i32 code
    2: string msg
    3: GetTemplatesResData data
}
struct GetTemplatesResData{
    1: list<template.Template> templates
}

struct AddTemplateItemReq{
    1: i64 template_id
    2: string name (api.body="name,required",api.vd="len($)>0")
    3: string content (api.body="content,required")
}
struct AddTemplateItemRes{
    1: i32 code
    2: string msg
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
    2: string name (api.body="name",api.vd="len($)>0")
    3: string content (api.body="content")
}
struct UpdateTemplateItemRes{
    1: i32 code
    2: string msg
}

struct GetTemplateItemsReq{
    1: i64 id (api.query="id",api.vd="$>=0")
    2: i32 page (api.query="page",api.vd="$>=0")
    3: i32 limit (api.query="limit",api.vd="$>=0")
    4: i32 order (api.query="order",api.vd="$>=0")
    5: string order_by (api.query="order_by")
}
struct GetTemplateItemsRes{
    1: i32 code
    2: string msg
    3: GetTemplateItemsResData data
}
struct GetTemplateItemsResData{
    1: list<template.TemplateItem> template_items
}

service TemplateService {
    AddTemplateRes AddTemplate(1: AddTemplateReq req) (api.post="/api/template")
    DeleteTemplateRes DeleteTemplate(1: DeleteTemplateReq req) (api.delete="/api/template")
    UpdateTemplateRes UpdateTemplate(1: UpdateTemplateReq req) (api.patch="/api/template")
    GetTemplatesRes GetTemplates(1: GetTemplatesReq req) (api.get="/api/template")

    AddTemplateItemRes AddTemplateItem(1: AddTemplateItemReq req) (api.post="/api/template/item")
    DeleteTemplateItemRes DeleteTemplateItem(1: DeleteTemplateItemReq req) (api.delete="/api/template/item")
    UpdateTemplateItemRes UpdateTemplateItem(1: UpdateTemplateItemReq req) (api.patch="/api/template/item")
    GetTemplateItemsRes GetTemplateItems(1: GetTemplateItemsReq req) (api.get="/api/template/item")
}