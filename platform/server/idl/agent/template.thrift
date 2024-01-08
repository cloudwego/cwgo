namespace go agent

include "../model/template.thrift"

struct AddTemplateReq{
    1: string name
    2: i32 type
}
struct AddTemplateResp{
    1: i32 code
    2: string msg
}

struct DeleteTemplateReq{
    1: list<i64> ids
}
struct DeleteTemplateResp{
    1: i32 code
    2: string msg
}

struct UpdateTemplateReq{
    1: i64 id
    2: string name
}
struct UpdateTemplateResp{
    1: i32 code
    2: string msg
}

struct GetTemplatesReq{
    1: i32 page
    2: i32 limit
    3: i32 order
    4: string order_by
}
struct GetTemplatesResp{
    1: i32 code
    2: string msg
    3: GetTemplatesRespData data
}
struct GetTemplatesRespData{
    1: list<template.Template> templates
}

struct AddTemplateItemReq{
    1: i64 template_id
    2: string name
    3: string content
}
struct AddTemplateItemResp{
    1: i32 code
    2: string msg
}

struct DeleteTemplateItemReq{
    1: list<i64> ids
}
struct DeleteTemplateItemResp{
    1: i32 code
    2: string msg
}

struct UpdateTemplateItemReq{
    1: i64 id
    2: string name
    3: string content
}
struct UpdateTemplateItemResp{
    1: i32 code
    2: string msg
}

struct GetTemplateItemsReq{
    1: i64 template_id
    2: i32 page
    3: i32 limit
    4: i32 order
    5: string order_by
}
struct GetTemplateItemsResp{
    1: i32 code
    2: string msg
    3: GetTemplateItemsRespData data
}
struct GetTemplateItemsRespData{
    1: list<template.TemplateItem> template_items
}
