namespace go model

enum Type {
    sync_idl_data,
    sync_repo_data,
}

union Data {
    1: SyncIdlData syncIdlData,
    2: SyncRepoData syncRepoData,
}

struct Task {
    1: string Id,
    2: Type Type,
    3: string ScheduleTime,
    4: Data Data,
}

struct SyncIdlData {
    1: i64 IdlId,
}

struct SyncRepoData {
    1: i64 RepositoryId,
}
