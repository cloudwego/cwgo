namespace go agent

include "../model/template.thrift"

struct AddTemplateReq{
    1: string name
    2: i32 type
}
struct AddTemplateRes{
    1: i32 code
    2: string msg
}

struct DeleteTemplateReq{
    1: list<i64> ids
}
struct DeleteTemplateRes{
    1: i32 code
    2: string msg
}

struct UpdateTemplateReq{
    1: i64 id
    2: string name
}
struct UpdateTemplateRes{
    1: i32 code
    2: string msg
}

struct GetTemplatesReq{
    1: i32 page
    2: i32 limit
    3: i32 order
    4: string order_by
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
    2: string name
    3: string content
}
struct AddTemplateItemRes{
    1: i32 code
    2: string msg
}

struct DeleteTemplateItemReq{
    1: list<i64> ids
}
struct DeleteTemplateItemRes{
    1: i32 code
    2: string msg
}

struct UpdateTemplateItemReq{
    1: i64 id
    2: string name
    3: string content
}
struct UpdateTemplateItemRes{
    1: i32 code
    2: string msg
}

struct GetTemplateItemsReq{
    1: i32 id
}
struct GetTemplateItemsRes{
    1: i32 code
    2: string msg
    3: GetTemplateItemsResData data
}
struct GetTemplateItemsResData{
    1: list<template.TemplateItem> template_items
}