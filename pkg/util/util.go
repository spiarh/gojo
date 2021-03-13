package util

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// IsTTYAllocated returns true when a TTY is allocated.
func IsTTYAllocated() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

func ParseFQIN(fqin string) (string, string, string, error) {
	var registry, name, tag string
	var errMsg = fmt.Errorf("invalid image name: %s", fqin)

	_, err := url.Parse(fqin)
	if err != nil {
		return registry, name, tag, errors.Wrapf(err, errMsg.Error())
	}

	// Get tag first with ":" separator
	repo_tag := strings.Split(fqin, ":")
	if len(repo_tag) == 1 {
		return registry, name, tag, errMsg
	}
	// in case port is specified, e.g: localhost:5000
	repo := strings.Join(repo_tag[:len(repo_tag)-1], ":")
	tag = repo_tag[len(repo_tag)-1]

	// Get registry and name with "/" separator
	registry_name := strings.Split(repo, "/")
	if len(registry_name) == 1 {
		return registry, name, tag, errMsg
	}
	registry = strings.Join(registry_name[:len(registry_name)-1], "/")
	name = registry_name[len(registry_name)-1]

	return registry, name, tag, nil
}

func SanitizeVersion(v string) string {
	return strings.TrimPrefix(v, "v")
}

func ExtractFileAsByteFromTar(in io.Reader, filename string) ([]byte, error) {
	tarReader := tar.NewReader(in)
	buf := new(bytes.Buffer)
	found := false
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			if header.Name == filename {
				found = true
				_, err := io.Copy(buf, tarReader)
				if err != nil {
					return nil, err
				}
				break
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("file not found in archive: %s", filename)
	}

	return buf.Bytes(), nil
}

func WriteToFile(path string, data []byte, perm fs.FileMode) (err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); err != nil {
			err = closeErr
		}
	}()

	err = ioutil.WriteFile(path, data, perm)
	return
}

func MakeDir(path string, perm fs.FileMode) error {
	fileInfo, err := os.Stat(path)

	if fileInfo != nil {
		if !fileInfo.IsDir() {
			return fmt.Errorf("provided path is not a directory, %s", path)
		}
	}

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, perm)
		if errDir != nil {
			return errDir
		}

	}
	return nil
}

func GetRelPathFromPathInTree(root, file string) (string, error) {
	rootSlice := strings.Split(root, "/")
	fileSlice := strings.Split(file, "/")

	for i := range rootSlice {
		if rootSlice[i] != fileSlice[i] {
			return "", fmt.Errorf("paths are not in the same tree: %s, %s", root, file)
		}
	}

	tmp := strings.Split(file, "/")[len(rootSlice):]
	relPath := strings.Join(tmp, "/")

	return relPath, nil
}

func EnsureStringSliceDuplicates(stringSlice []string) error {
	keys := make(map[string]struct{})
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = struct{}{}
			continue
		}
		return fmt.Errorf("duplicate values")
	}
	return nil
}
