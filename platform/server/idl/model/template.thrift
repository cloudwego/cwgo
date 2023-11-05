namespace go model

struct Template{
    1: i64 id
    2: string name
    3: i32 type // 1: hz, 2: kitex
    4: bool is_deleted
    5: string create_time
    6: string update_time
}

struct TemplateItem{
    1: i64 id
    2: i64 template_id
    3: string name
    4: string content
    5: bool is_deleted
    6: string create_time
    7: string update_time
}