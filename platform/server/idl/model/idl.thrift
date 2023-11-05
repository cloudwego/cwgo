namespace go model

struct IDL{
    1: i64 id // idl record id
    2: i64 repository_id // repository id where stores the idl
    3: string main_idl_path // repo ref path of idl
    4: string commit_hash // idl file commit hash
    5: list<ImportIDL> import_idls
    6: string service_name // service name
    7: string last_sync_time // idl last sync time
    8: bool is_deleted
    9: string create_time
    10: string update_time
}

struct ImportIDL{
    1: string idl_path
    2: string commit_hash
}