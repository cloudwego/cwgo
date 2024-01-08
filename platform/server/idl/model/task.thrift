namespace go model

enum TaskType {
    Sync = 1
}

struct Task {
    1: string ID
    2: TaskType Type
    3: string ScheduleTime
    4: i64 IdlID
}
