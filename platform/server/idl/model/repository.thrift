namespace go model

struct Repository{
    1: i64 id // repository record id
    2: i32 repository_type // repository type (1: gitlab, 2: github)
    3: string repository_domain // repository domain
    4: string repository_owner // repository owner
    5: string repository_name // repository name
    6: string repository_branch // repository branch
    7: i32 store_type // repository store type (1: stores idl file, 2: stores meta code)
    8: i64 token_id // token id which repo currently using
    9: i32 status // repository status (need sync or not)
    10: string last_update_time // repo file last update time
    11: string last_sync_time // repo last sync time
    12: bool is_deleted
    13: string create_time
    14: string update_time
}
