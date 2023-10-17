package parser

type Parser interface {
	GetDependentFilePaths(mainIdlPath string) ([]string, error)
}
