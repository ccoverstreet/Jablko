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
	downloadURL, err := GithubSourceToURL(sourcePath)
	if err != nil {
		return err
	}

	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Bad HTTP status code received: %d", resp.StatusCode)
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

	_, err = Unzip("./tmp/" + sourcePath +  "/source.zip", "./" + splitSource[0] + "/" + splitSource[1])
	if err != nil {
		return err
	}

	// Replace module line in go.mod to include version
	goModFile, err := os.Open(sourcePath + "/go.mod")
	if err != nil {
		return err
	}

	modScanner := bufio.NewScanner(goModFile)
	modScanner.Split(bufio.ScanLines)
	var modLines []string

	// Go through go.mod file's lines and replace module name
	// for uniqueness
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
	if err != nil {
		return err
	}

	// Build newly downloaded and prepped module
	err = BuildJablkoMod(sourcePath)

	return err
}
