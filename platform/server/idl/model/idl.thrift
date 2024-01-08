namespace go model

include "repository.thrift"

struct IDL{
    1: i64 id // idl record id
    2: i64 idl_repository_id // repository id where stores the idl
    3: i64 service_repository_id // repository id where stores the meta
    4: string main_idl_path // repo ref path of idl
    5: string commit_hash // idl file commit hash
    6: list<ImportIDL> import_idls
    7: string service_name // meta name
    8: string last_sync_time // idl last sync time
    9: i32 status // idl status
    10: bool is_deleted
    11: string create_time
    12: string update_time
}

struct ImportIDL{
    1: string idl_path
    2: string commit_hash
}

struct IDLWithRepositoryInfo{
    1: i64 id // idl record id
    2: i64 idl_repository_id // repository id where stores the idl
    3: repository.Repository idl_repository // idl repository info
    4: i64 service_repository_id // repository id where stores the meta
    5: repository.Repository service_repository // meta repository info
    6: string main_idl_path // repo ref path of idl
    7: string commit_hash // idl file commit hash
    8: list<ImportIDL> import_idls
    9: string service_name // meta name
    10: string last_sync_time // idl last sync time
    11: i32 status // idl status
    12: bool is_deleted
    13: string create_time
    14: string update_time
}
