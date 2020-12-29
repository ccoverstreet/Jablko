// downloading.go: Download procedure for jablkomods
// Cale Overstreet
// 2020/12/25
// Handles the downloading of jablkomods from Github repos. This
// should be able to download at any point during the runtime.
// Currently, it is called when required during the Initialize 
// function inside of src/jablkomods/jablkomods.go.

package jablkomods

import (
	"net/http"
	"io"
	"bufio"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
)

func DownloadJablkoMod(repoPath string) error {
	// Download zip, branch based on provider
	if strings.HasPrefix(repoPath, "github.com") {
		return downloadGithub(repoPath)
	}

	return fmt.Errorf("Provider not supported")
}

func downloadGithub(sourcePath string) error {
	downloadURL := GithubSourceToURL(sourcePath)

	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Bad HTTP status code received\n")
	}

	defer resp.Body.Close()

	// Create file
	err = os.MkdirAll("./tmp/" + sourcePath, 0755)
	if err != nil {
		return err
	}

	out, err := os.Create("./tmp/" + sourcePath +  "/source.zip")
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return err
	}

	err = os.MkdirAll("./" + sourcePath, 0755)
	if err != nil {
		return err
	}

	splitSource := strings.Split(sourcePath, "/")
	installPath := GithubSourceToInstallDir(sourcePath)

	_, err = Unzip("./tmp/" + sourcePath +  "/source.zip", "./" + splitSource[0] + "/" + splitSource[1])
	if err != nil {
		return err
	}

	// Replace module line in go.mod to include version
	goModFile, err := os.Open(installPath + "/go.mod")
	if err != nil {
		return err
	}

	modScanner := bufio.NewScanner(goModFile)
	modScanner.Split(bufio.ScanLines)
	var modLines []string

	for modScanner.Scan() {
		tempLine := modScanner.Text()
		if strings.HasPrefix(tempLine, "module") {
			modLines = append(modLines, "module " + sourcePath)
		} else {
			modLines = append(modLines, tempLine)
		}
	}

	goModFile.Close()

	// Write sanitized file to go.mod
	err = ioutil.WriteFile(sourcePath + "/go.mod", []byte(strings.Join(modLines, "\n")),0666)

	return err
}
