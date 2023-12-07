package consts

const (
	MainRef = "main"

	InvalidToken = "invalid_token"
)

const (
	RepositoryTypeNumGitLab = iota + 1
	RepositoryTypeNumGithub
	RepositoryTypeNum = iota
)

var (
	RepositoryTypeNumMap = map[int]struct{}{
		RepositoryTypeNumGitLab: {},
		RepositoryTypeNumGithub: {},
	}
)

const (
	RepositoryStoreTypeNumIdl = iota + 1
	RepositoryStoreTypeNumService
	RepositoryStoreTypeNum = iota
)

const (
	RepositoryStatusNumActive = iota + 1
	RepositoryStatusNumInactive
)

var (
	RepositoryStatusNumMap = map[int]struct{}{
		RepositoryStatusNumActive:   {},
		RepositoryStatusNumInactive: {},
	}
)
