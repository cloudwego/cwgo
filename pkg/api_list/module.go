package api_list

import (
	"fmt"

	"github.com/cloudwego/cwgo/pkg/common/utils"
)

func getModuleName(path string) (string, error) {
	module, _, ok := utils.SearchGoMod(path, false)
	if !ok {
		return "", fmt.Errorf("path: %s not found go.mod", path)
	}

	return module, nil
}
