namespace go model

struct IDL{
    1: i64 id // idl record id
    2: i64 idl_repository_id // repository id where stores the idl
    3: i64 service_repository_id // repository id where stores the service
    4: string main_idl_path // repo ref path of idl
    5: string commit_hash // idl file commit hash
    6: list<ImportIDL> import_idls
    7: string service_name // service name
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