package docker

import (
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"strings"
)

// DockerfileInfo contains information used to generate a Dockerfile for a template.
type DockerfileInfo struct {
	Mirrors    string   // Mirrors holds the mirror URLs
	GoMainFrom string   // GoMainFrom is the path to the Go main file
	GoFile     string   // GoFile is the Go file Path (main file)
	ExeFile    string   // ExeFile is the Path of the executable file
	BaseImage  string   // BaseImage is the base image for the Dockerfile
	Port       int      // Port is the port number
	EtcDirs    []string // EtcDirs contains the directory containing the etc files
	Argument   string   // Argument contains the list of build arguments for the Dockerfile
	Version    string   // Version is the version number of the application
	Timezone   string   // Timezone is the timezone for the application
}

// FillInfo populates a DockerfileInfo struct with data from the DockerArgument configuration.
func FillInfo(c *config.DockerArgument) *DockerfileInfo {
	// Initialize mirror URL builder
	var mirrorBuilder strings.Builder
	// Check if mirrors are provided and construct the mirror string
	if len(c.Mirrors) > 0 {
		for _, mirror := range c.Mirrors {
			// Ensure each mirror URL starts with the correct prefix
			if !strings.HasPrefix(mirror, consts.HttpsPrefix) {
				mirrorBuilder.WriteString(consts.HttpsPrefix)
			}
			mirrorBuilder.WriteString(mirror)
			mirrorBuilder.WriteString(consts.MirrorSep)
		}
		// Append "direct" at the end of the mirror list
		mirrorBuilder.WriteString(consts.MirrorDirect)
	}

	// Initialize arguments slice for the Dockerfile
	var args []string

	// Check if there are any etc files to include
	var hasEtcFile, etcPaths = checkEtcPath(c.EtcDirs)
	if hasEtcFile {
		args = append(args, "-f")
		for _, path := range etcPaths {
			var err error
			args, err = utils.GetAllFile(path, args)
			if err != nil {
				return nil
			}
		}
	}

	// Append any additional arguments provided in the DockerArgument
	args = append(args, c.Arguments...)
	var argsStr string
	for _, str := range args {
		if !strings.HasPrefix(str, "\"") {
			str = "\"" + str
		}
		if !strings.HasSuffix(str, "\"") {
			str = str + "\"" // Wrap each argument in double quotes for Dockerfile compatibility
		}
		argsStr = argsStr + ", " + str // Separate arguments with spaces for Dockerfile compatibility
	}

	// Return a populated DockerfileInfo struct
	return &DockerfileInfo{
		Mirrors:    mirrorBuilder.String(), // Constructed mirror URLs
		GoMainFrom: c.Main,                 // Path to the main Go file
		GoFile:     c.Main,                 // Main Go file Path
		ExeFile:    c.ExeName,              // Executable file Path
		BaseImage:  c.BaseImage,            // Base Docker image
		Port:       int(c.Port),            // Port number
		EtcDirs:    etcPaths,               // Etc directory containing
		Argument:   argsStr,                // List of Docker build arguments
		Version:    c.Version,              // Version number of the application
		Timezone:   c.TZ,                   // Timezone setting
	}
}

// checkEtcPath checks if there are any files in the "etc" directory.
func checkEtcPath(paths []string) (bool, []string) {
	etdPaths := make([]string, 0)
	for _, dir := range paths {
		// 检查路径是否存在
		if exists, err := utils.PathExist(dir); err != nil {
			return false, nil // 直接返回，如果路径不存在或发生错误
		} else if !exists {
			continue
		}

		// 路径有效，添加到 etdPaths 列表
		etdPaths = append(etdPaths, dir)
	}

	return len(etdPaths) > 0, etdPaths
}
