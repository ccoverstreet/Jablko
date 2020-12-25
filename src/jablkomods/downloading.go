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
	"archive/zip"
	"path/filepath"
	"os"
	"fmt"
	"strings"
)

func DownloadJablkoMod(repoPath string) error {
	// Download zip from github
	resp, err := http.Get("https://" + repoPath + "/archive/master.zip")
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Bad HTTP status code received\n")
	}

	defer resp.Body.Close()

	// Create file
	err = os.MkdirAll("./tmp/" + repoPath, 0755)
	if err != nil {
		return err
	}

	out, err := os.Create("./tmp/" + repoPath +  "/source.zip")
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return err
	}

	err = os.MkdirAll("./" + repoPath, 0755)
	if err != nil {
		return err
	}

	filenames, err := Unzip("./tmp/" + repoPath +  "/source.zip", "./" + repoPath)
	if err != nil {
		return err
	}

	topLevelDirRaw := strings.Split(filenames[0], "/")
	topLevelDir := topLevelDirRaw[len(topLevelDirRaw) - 1]
	repoPathSplit := strings.Split(repoPath, "/")
	authorDir := repoPathSplit[0] + "/" + repoPathSplit[1]

	// Move directory up one level and rename correctly
	err = os.Rename("./" + repoPath + "/" + topLevelDir, "./" + authorDir + "/" + topLevelDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll("./" + repoPath)
	if err != nil {
		return err
	}

	err = os.Rename("./" + authorDir + "/" + topLevelDir, "./" + repoPath)
	if err != nil {
		return err
	}

	return err
}

func Unzip(src string, dest string) ([]string, error) {
	// From golangcode.com

    var filenames []string

    r, err := zip.OpenReader(src)
    if err != nil {
        return filenames, err
    }
    defer r.Close()

    for _, f := range r.File {

        // Store filename/path for returning and using later on
        fpath := filepath.Join(dest, f.Name)

        // Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return filenames, fmt.Errorf("%s: illegal file path", fpath)
        }

        filenames = append(filenames, fpath)

        if f.FileInfo().IsDir() {
            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        // Make File
        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return filenames, err
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return filenames, err
        }

        rc, err := f.Open()
        if err != nil {
            return filenames, err
        }

        _, err = io.Copy(outFile, rc)

        // Close the file without defer to close before next iteration of loop
        outFile.Close()
        rc.Close()

        if err != nil {
            return filenames, err
        }
    }
    return filenames, nil
}
