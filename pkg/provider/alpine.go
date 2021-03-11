package provider

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/lcavajani/gojo/pkg/util"
)

const (
	alpineAPKIndexArchiveName = "APKINDEX.tar.gz"
	alpineAPKIndexFilename    = "APKINDEX"
	alpineDefaultMirror       = "http://dl-cdn.alpinelinux.org"
	alpineOS                  = "alpine"
)

type Alpine struct {
	log zerolog.Logger

	arch      string
	mirror    string
	repo      string
	versionId string
	pkgName   string
}

func NewAlpine(mirror, arch, versionId, repo, pkgName string) *Alpine {
	return &Alpine{
		log:       log.With().Str("provider", string(ProviderAlpine)).Logger(),
		mirror:    mirror,
		arch:      arch,
		versionId: versionId,
		repo:      repo,
		pkgName:   pkgName,
	}
}

func (a *Alpine) buildURL() (*url.URL, error) {
	// http://dl-cdn.alpinelinux.org/alpine/v3.13/main/x86_64/APKINDEX.tar.gz
	v := "v" + a.versionId
	path := path.Join(alpineOS, v, a.repo, a.arch, alpineAPKIndexArchiveName)
	u := fmt.Sprintf("%s/%s", a.mirror, path)
	return url.Parse(u)
}

func (a *Alpine) getAPKIndexArchive() (io.ReadCloser, error) {
	apkIndexURL, err := a.buildURL()
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	resp, err := client.Get(apkIndexURL.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error getting apk index file: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (a *Alpine) getAPKIndexFromArchive(archive io.ReadCloser) ([]byte, error) {
	gzReader, err := gzip.NewReader(archive)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	apkIndex, err := util.ExtractFileAsByteFromTar(gzReader, alpineAPKIndexFilename)
	if err != nil {
		return nil, err
	}

	return apkIndex, nil
}

func (a *Alpine) GetLatest() (string, error) {
	apkIndexArchive, err := a.getAPKIndexArchive()
	if err != nil {
		return "", err
	}
	defer apkIndexArchive.Close()

	apkIndex, err := a.getAPKIndexFromArchive(apkIndexArchive)
	if err != nil {
		return "", err
	}

	pkg, err := parseAPKIndex(apkIndex, a.pkgName)
	if err != nil {
		return "", err
	}
	fmt.Println(pkg.version)

	return pkg.version, nil
}

type AlpinePackageMeta struct {
	name    string
	version string
	arch    string
}

func (a *AlpinePackageMeta) isValid() error {
	errFunc := func(field string) error { return fmt.Errorf("Alpine package meta field is empty: %s", field) }
	if a.name == "" {
		return errFunc("Name")
	}
	if a.arch == "" {
		return errFunc("Arch")
	}
	if a.version == "" {
		return errFunc("Version")
	}
	return nil
}

func parseAPKIndex(apkIndex []byte, pkgName string) (*AlpinePackageMeta, error) {
	sc := bufio.NewScanner(bytes.NewReader(apkIndex))
	a := AlpinePackageMeta{}

Loop:
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "P:"):
			a.name = strings.Split(line, ":")[1]
		case strings.HasPrefix(line, "V:"):
			a.version = strings.Split(line, ":")[1]
		case strings.HasPrefix(line, "A:"):
			a.arch = strings.Split(line, ":")[1]
		case line == "":
			if a.name == pkgName {
				break Loop
			}
			// Start the new block with an empty struct
			a = AlpinePackageMeta{}
		}
	}

	if err := sc.Err(); err != nil {
		return nil, errors.Wrap(err, "Invalid input")
	}

	if a.name == pkgName {
		if err := a.isValid(); err != nil {
			return nil, err
		}
		return &a, nil
	}

	return nil, fmt.Errorf("Package not found: %s", pkgName)
}
