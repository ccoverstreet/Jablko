// utils.go: Utility functions for jablkomods
// Cale Overstreet
// 2020/12/29
// Used globally throughout jablkomods

package jablkomods

import (
	"strings"
)

func GithubSourceToURL(sourceStr string) string {
	splitSource := strings.Split(sourceStr, "/")
	splitRepo := strings.Split(splitSource[2], "-")
	if splitRepo[1] == "master" {
		return "https://" + splitSource[0] + "/" + splitSource[1] + "/" + splitRepo[0] + "/archive/" + splitRepo[1] + ".zip"
	}

	return "https://" + splitSource[0] + "/" + splitSource[1] + "/" + splitRepo[0] + "/archive/v" + splitRepo[1] + ".zip"
}

func GithubSourceToInstallDir(sourceStr string) string {
	splitSource := strings.Split(sourceStr, "/")

	if strings.HasSuffix(sourceStr, "master") {
		return splitSource[0] + "/" + splitSource[1] + "/" + splitSource[2] + "-" + splitSource[3]
	} else {
		return splitSource[0] + "/" + splitSource[1] + "/" + splitSource[2] + "-" + splitSource[3][1:len(splitSource[3])]
	}
}
