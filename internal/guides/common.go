package guides

import (
	"errors"
	"strings"
)

var (
	ErrDirDoesNotContainAnyConfigFile  = errors.New("this directory does not contain any config (.yaml|.yml) file")
	ErrDirDoesNotContainAnyDiagramFile = errors.New("this directory does not contain any diagram (.xml) file")
)

func replaceDoubleSlash(str string) string {
	return strings.ReplaceAll(str, "//", "/")
}
