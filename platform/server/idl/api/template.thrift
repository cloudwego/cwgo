namespace go template

struct Template{
    1: i64 id
    2: string name
    3: i8 type
}

struct TemplateItem{
    1: i64 id
    2: i64 template_id
    3: string name
    4: string content
}

struct AddTemplateReq{
    1: string name (api.body="name,required")
    2: i8 type (api.body="type,required") // 1: hz, 2: kitex
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

struct GetTemplateReq{
    1: i32 page (api.query="page",api.vd="$>=0")
    2: i32 limit (api.query="limit",api.vd="$>=0")
}
struct GetTemplateRes{
    1: i32 code
    2: string msg
    3: GetTemplateResData data
}
struct GetTemplateResData{
    1: list<Template> templates
}

struct AddTemplateItemReq{
    1: string name (api.body="name,required")
    2: string content (api.body="content,required")
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
}
struct GetTemplateItemsRes{
    1: i32 code
    2: string msg
    3: GetTemplateItemsResData data
}
struct GetTemplateItemsResData{
    1: list<TemplateItem> template_items
}

service TemplateService {
    AddTemplateReq AddTemplate(1: AddTemplateReq req) (api.post="/template")
    DeleteTemplateReq DeleteTemplate(1: DeleteTemplateReq req) (api.delete="/template")
    UpdateTemplateReq UpdateTemplate(1: UpdateTemplateReq req) (api.patch="/template")
    GetTemplateReq GetTemplate(1: GetTemplateItemsReq req) (api.get="/template")

    AddTemplateItemReq AddTemplateItem(1: AddTemplateItemReq req) (api.post="/template/item")
    DeleteTemplateReq DeleteTemplateItem(1: DeleteTemplateItemReq req) (api.delete="/template/item")
    UpdateTemplateItemReq UpdateTemplateItem(1: UpdateTemplateItemReq req) (api.patch="/template/item")
    GetTemplateItemsReq GetTemplateItem(1: GetTemplateItemsReq req) (api.get="/template/item")
}