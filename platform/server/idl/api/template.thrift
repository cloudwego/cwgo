namespace go template

include "../base/model.thrift"

struct AddTemplateReq{
    1: string name (api.body="name,required")
    2: i32 type (api.body="type,required")
}
struct AddTemplateRes{
    1: i32 code
    2: string msg
}

struct DeleteTemplateReq{
    1: list<i64> ids (api.body="ids,required")
}
struct DeleteTemplateRes{
    1: i32 code
    2: string msg
}

struct UpdateTemplateReq{
    1: i64 id (api.body="id,required")
    2: string name (api.body="name,required")
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
    1: list<model.Template> templates
}

struct AddTemplateItemReq{
    1: i64 template_id
    2: string name (api.body="name,required")
    3: string content (api.body="content,required")
}
struct AddTemplateItemRes{
    1: i32 code
    2: string msg
}

struct DeleteTemplateItemReq{
    1: list<i64> ids (api.body="ids,required")
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

struct GetTemplateItemsReq{
    1: i32 page (api.query="page",api.vd="$>=0")
    2: i32 limit (api.query="limit",api.vd="$>=0")
    3: i32 order (api.query="order",api.vd="$>=0")
    4: string order_by (api.query="order_by")
}
struct GetTemplateItemsRes{
    1: i32 code
    2: string msg
    3: GetTemplateItemsResData data
}
struct GetTemplateItemsResData{
    1: list<model.TemplateItem> template_items
}

service TemplateService {
    AddTemplateRes AddTemplate(1: AddTemplateReq req) (api.post="/template")
    DeleteTemplateRes DeleteTemplate(1: DeleteTemplateReq req) (api.delete="/template")
    UpdateTemplateRes UpdateTemplate(1: UpdateTemplateReq req) (api.patch="/template")
    GetTemplatesRes GetTemplates(1: GetTemplateItemsReq req) (api.get="/template")

    AddTemplateItemRes AddTemplateItem(1: AddTemplateItemReq req) (api.post="/template/item")
    DeleteTemplateRes DeleteTemplateItem(1: DeleteTemplateItemReq req) (api.delete="/template/item")
    UpdateTemplateItemRes UpdateTemplateItem(1: UpdateTemplateItemReq req) (api.patch="/template/item")
    GetTemplateItemsRes GetTemplateItems(1: GetTemplatesReq req) (api.get="/template/item")
}