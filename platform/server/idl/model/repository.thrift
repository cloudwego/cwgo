namespace go model

struct Repository{
    1: i64 id // repository record id
    2: i32 repository_type // repository type (1: gitlab, 2: github)
    3: i32 store_type // repository store type (1: stores idl file, 2: stores service code)
    4: string repository_url // repository url (full url)
    5: string token // repository token
    6: i32 status // repository status (need sync or not)
    7: string last_update_time // repo file last update time
    8: string last_sync_time // repo last sync time
    9: bool is_deleted
    10: string create_time
    11: string update_time
}