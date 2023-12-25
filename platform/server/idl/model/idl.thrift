namespace go model

include "repository.thrift"
include "template.thrift"

struct IDL{
    1: i64 id // idl record id
    2: i64 idl_repository_id // repository id where stores the idl
    3: i64 service_repository_id // repository id where stores the service
    4: string main_idl_path // repo ref path of idl
    5: string commit_hash // idl file commit hash
    6: list<ImportIDL> import_idls
    7: string service_name // service name
    8: string last_sync_time // idl last sync time
    9: i64 template_id
    10: i32 status // idl status
    11: bool is_deleted
    12: string create_time
    13: string update_time
}

struct ImportIDL{
    1: string idl_path
    2: string commit_hash
}

struct IDLWithInfo{
    1: i64 id // idl record id
    2: i64 idl_repository_id // repository id where stores the idl
    3: repository.Repository idl_repository // idl repository info
    4: i64 service_repository_id // repository id where stores the service
    5: repository.Repository service_repository // service repository info
    6: string main_idl_path // repo ref path of idl
    7: string commit_hash // idl file commit hash
    8: list<ImportIDL> import_idls
    9: string service_name // service name
    10: i64 template_id
    11: template.Template template
    12: i32 status // idl status
    13: string last_sync_time // idl last sync time
    14: bool is_deleted
    15: string create_time
    16: string update_time
}