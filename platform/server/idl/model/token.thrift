namespace go model

struct Token{
    1: i64 id // id
    2: i32 repository_type // repository type (1: gitlab, 2: github)
    3: string repository_domain // repository api domain
    4: string owner // repository owner
    5: string token // repository token
    6: i32 status // token status (0: expired, 1: valid)
    7: string expiration_time // token expiration time
    8: bool is_deleted // is deleted
    9: string create_time // create time
    10: string update_time // update time
}