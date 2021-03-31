// utils.go: Utility functions for jablkomods
// Cale Overstreet
// 2020/12/29
// Used globally throughout jablkomods

package jablkomods

import (
	"strings"
	"archive/zip"
	"path/filepath"
	"os"
	"io"
	"fmt"
	"os/exec"
	"time"
	"strconv"

	"github.com/ccoverstreet/Jablko/src/jlog"
)

func BuildJablkoMod(buildDir string) error {
	buildCMD := exec.Command("go", "build", "-buildmode=plugin", "-o", "jablkomod.so", ".")		
	buildCMD.Dir = buildDir
	jlog.Println(buildCMD)
	outBytes, err := buildCMD.CombinedOutput()

	jlog.Printf("Build Output of \"%s\":\n%s", buildDir, string(outBytes))

	return err
}

func GithubSourceToURL(sourceStr string) (string, error){
	splitSource := strings.Split(sourceStr, "/")
	splitRepo := strings.Split(splitSource[2], "-")

	// Check if split repo is correct
	if len(splitRepo) < 2 {
		return "", fmt.Errorf("Repo name does not include version.")	
	}

	if splitRepo[1] == "master" {
		return "https://" + splitSource[0] + "/" + splitSource[1] + "/" + splitRepo[0] + "/archive/" + splitRepo[1] + ".zip", nil
	}

	return "https://" + splitSource[0] + "/" + splitSource[1] + "/" + splitRepo[0] + "/archive/v" + splitRepo[1] + ".zip", nil
}

func CreateUUID() string {
	curTime := int(time.Now().UnixNano() / 1000)

	idStr := strconv.Itoa(curTime)
	idBytes := []byte(idStr)

	return string(idBytes)
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
