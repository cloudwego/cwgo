package consts

const (
	MainRef      = "main"
	GitHubDomain = "github.com"
)

const (
	RepositoryTypeNumGitLab = iota + 1
	RepositoryTypeNumGithub
	RepositoryTypeNum = iota
)

var RepositoryTypeNumMap = map[int]struct{}{
	RepositoryTypeNumGitLab: {},
	RepositoryTypeNumGithub: {},
}

const (
	RepositoryStoreTypeNumIdl = iota + 1
	RepositoryStoreTypeNumService
	RepositoryStoreTypeNum = iota
)

const (
	RepositoryStatusNumInactive = iota + 1
	RepositoryStatusNumActive
)

var RepositoryStatusNumMap = map[int]struct{}{
	RepositoryStatusNumActive:   {},
	RepositoryStatusNumInactive: {},
}
