package pkgFile

import (
	"embed"
	"fmt"
	"io/fs"
)

// FetchSourcePath
// Do exactly in this way because we need to use custom paths in tests
var FetchSourcePath = "source/"

// FetchSourceFile - helper function to easily open embedded files
//
// Function is created for the rule file.FilesAreSameRule
// How does it work? - from the projectlinter side, we put the file in the "source" folder (architectural solution)
// Then we get it through go:embed and compare it with the file in project
func FetchSourceFile(files embed.FS, fileName string) (fs.File, error) {
	f, err := files.Open(FetchSourcePath + fileName)
	if err != nil {
		return nil, fmt.Errorf("projectlinter has no embedded file \"%s\"", fileName)
	}

	return f, nil
}
