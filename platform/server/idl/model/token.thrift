namespace go model

struct Token{
    1: i64 id // id
    2: i32 repository_type // repository type (1: gitlab, 2: github)
    3: string repository_domain // repository api domain
    4: string owner // token owner
    5: i64 owner_id // token owner id
    6: i32 token_type // token type (1: personal, 2: organization)
    7: string token // repository token
    8: i32 status // token status (1: expired, 2: valid)
    9: string expiration_time // token expiration time
    10: bool is_deleted // is deleted
    11: string create_time // create time
    12: string update_time // update time
}