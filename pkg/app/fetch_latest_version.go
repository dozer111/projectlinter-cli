package pkgApp

import (
	"github.com/Masterminds/semver/v3"
)

func FetchLatestVersion() *semver.Version {
	// TODO change this example to your real code
	// yourGitClient := NewYourGitClient
	//latestVersionFile, err := yourGitClient.FetchFile(
	//	"your_path_to_projectlinter/raw/cmd/projectlinter/version",
	//)
	//if err != nil {
	//	log.Panicf("cannot receive latest version: %v", err)
	//}

	//return semver.MustParse(string(latestVersionFile))
	return semver.MustParse("1.0.0")
}
