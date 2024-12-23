package docker

import (
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"path/filepath"
	"strings"
)

// DockerfileInfo contains information used to generate a Dockerfile for a template.
type DockerfileInfo struct {
	Mirrors   string // Mirrors holds the mirror URLs
	HasMirror bool   // HasMirror is a boolean indicating if mirrors are set

	GoMainFrom  string   // GoMainFrom is the path to the Go main file
	GoFile      string   // GoFile is the Go file Path (main file)
	ExeFile     string   // ExeFile is the Path of the executable file
	BaseImage   string   // BaseImage is the base image for the Dockerfile
	HasPort     bool     // HasPort is a boolean indicating whether a port is defined
	Port        int      // Port is the port number
	HasEtcFile  bool     // HasEtcFile is a boolean indicating if any etc files exist
	EtcDirs     []string // EtcDirs contains the directory containing the etc files
	Argument    string   // Argument contains the list of build arguments for the Dockerfile
	Version     string   // Version is the version number of the application
	HasTimezone bool     // HasTimezone indicates if a timezone is set
	Timezone    string   // Timezone is the timezone for the application
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
	if c.EtcDirs[0] != "." {
		c.EtcDirs = append(c.EtcDirs, ".")
	}
	var hasEtcFile, etcFiles = checkEtcFile(c.EtcDirs)
	if hasEtcFile {
		args = append(args, "-f")
		args = append(args, etcFiles...)
	}

	// Append any additional arguments provided in the DockerArgument
	args = append(args, c.Arguments...)
	//TODO:为每个arg环绕双引号
	//TODO:测试dockerfile可行性和各个参数有效性
	argsStr := strings.Join(args, ", ")
	if argsStr != "" {
		argsStr = ", " + argsStr
	}

	// Return a populated DockerfileInfo struct
	return &DockerfileInfo{
		Mirrors:     mirrorBuilder.String(), // Constructed mirror URLs
		HasMirror:   len(c.Mirrors) > 0,     // Whether mirrors are provided
		GoMainFrom:  c.Main,                 // Path to the main Go file
		GoFile:      c.Main,                 // Main Go file Path
		ExeFile:     c.ExeName,              // Executable file Path
		BaseImage:   c.BaseImage,            // Base Docker image
		HasPort:     c.Port > 0,             // Whether a port is defined
		Port:        int(c.Port),            // Port number
		HasEtcFile:  hasEtcFile,             // Whether etc files are present
		EtcDirs:     etcFiles,               // Etc directory containing
		Argument:    argsStr,                // List of Docker build arguments
		Version:     c.Version,              // Version number of the application
		HasTimezone: len(c.Tz) > 0,          // Whether a timezone is set
		Timezone:    c.Tz,                   // Timezone setting
	}
}

// checkEtcFile checks if there are any files in the "etc" directory.
func checkEtcFile(paths []string) (bool, []string) {
	etdPaths := make([]string, 0)
	for _, dir := range paths {
		var path string
		if !strings.HasSuffix(dir, "etc") {
			// Build the path to the "etc" directory
			path = filepath.Join(dir, "etc")
		} else {
			path = dir
		}
		// 尝试获取绝对路径
		absPath, err := filepath.Abs(path)
		if err != nil {
			return false, nil // 直接返回，避免继续执行
		}

		// 检查路径是否存在
		if exists, err := utils.PathExist(absPath); err != nil {
			return false, nil // 直接返回，如果路径不存在或发生错误
		} else if !exists {
			continue
		}

		// 路径有效，添加到 etdPaths 列表
		etdPaths = append(etdPaths, path)

	}

	return len(etdPaths) > 0, etdPaths
}
