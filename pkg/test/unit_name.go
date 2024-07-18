package pkgTest

import (
	"github.com/dozer111/projectlinter-core/util/path_provider"
	"strings"
)

func UnitName() string {
	pr := path_provider.NewPathProvider("")
	path := pr.PathToCaller()
	pathSlice := strings.Split(path, "/")
	res := pathSlice[len(pathSlice)-3]

	return res
}
