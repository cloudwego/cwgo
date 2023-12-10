package parser

type Parser interface {
	GetDependentFilePaths(mainIdlPath string) ([]string, error) // Obtain string slices that rely on main idl
}
