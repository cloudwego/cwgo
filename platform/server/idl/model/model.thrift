namespace go model

struct Task{
    1: string id
    2: i32 type
    3: string schedule_time
    4: string data
}

struct IDL{
    1: i64 id
    2: i64 repository_id
    3: string main_idl_path
    4: string content
    5: string service_name
    6: string last_sync_time
    7: bool is_deleted
    8: string create_time
    9: string update_time
}

struct Repository{
    1: i64 id
    2: i32 repository_type
    3: i32 store_type
    4: string repository_url
    5: string token
    6: string status
    7: string last_update_time
    8: string last_sync_time
    9: bool is_deleted
    10: string create_time
    11: string update_time
}

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