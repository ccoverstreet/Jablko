package github

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Retrieving tag data
// https://api.github.com/repos/ccoverstreet/Jablko/tags

func DownloadZipRepo(zipURL string) error {
	base := filepath.Base(zipURL)

	log.Printf("%v", zipURL)

	resp, err := http.Get(zipURL)
	if err != nil {
		return err
	}

	if !strings.Contains(base, ".zip") {
		base = base + ".zip"
	}

	tmpZip, err := os.Create("tmp/" + base)
	if err != nil {
		return err
	}
	defer tmpZip.Close()

	size, err := io.Copy(tmpZip, resp.Body)
	log.Printf("%d", size)
	defer resp.Body.Close()

	fmt.Println(size)

	if err != nil {
		return err
	}

	return nil
}

func UnpackZipRepo(zipURL string, dest string) error {
	base := filepath.Base(zipURL)
	log.Printf("%v", base)

	r, err := zip.OpenReader("tmp/" + base)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	// Path prefix is removed from all files in Zip and replaced with dest
	pathPrefix := strings.Split(r.File[0].FileHeader.Name, "/")[0]

	// Writes file from zip to destination file
	fExtractFile := func(zf *zip.File) error {
		zfHandle, err := zf.Open()
		if err != nil {
			return err
		}

		fileName := strings.Replace(zf.Name, pathPrefix, dest, -1)
		if zf.FileInfo().IsDir() {
			os.MkdirAll(fileName, zf.Mode())
		} else {
			f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, zfHandle)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// Iterate through all contents except for parent dir
	for i := 1; i < len(r.File); i++ {
		err := fExtractFile(r.File[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDefaultBranch(repoName string) (string, error) {
	type defaultBranchHolder struct {
		DefaultBranch string `json:"default_branch"`
	}

	trimmedName := strings.Replace(repoName, "github.com/", "", 1)
	url := "https://api.github.com/repos/" + trimmedName
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var resBody defaultBranchHolder

	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return "", err
	}

	return resBody.DefaultBranch, nil
}

func RetrieveSource(jmodName string, commit string) error {
	repoName := jmodName
	versionTag := commit
	//trimmedRepoName := strings.Replace(repoName, "github.com/", "", 1)

	defaultBranch, err := GetDefaultBranch(repoName)
	log.Printf("%v", defaultBranch)
	if err != nil {
		return err
	}
	url := "https://" + jmodName + "/archive/" + versionTag + ".zip"

	err = DownloadZipRepo(url)
	if err != nil {
		return err
	}

	err = UnpackZipRepo(url, jmodName)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

func DeleteSource(jmodPath string) error {
	return os.RemoveAll(jmodPath)
}
