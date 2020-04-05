package resource

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	fpath "path/filepath"
	"plugin"
	"regexp"

	"github.com/pkg/errors"
)

var securityPlugin *plugin.Plugin

var bannedExtensions = map[string]struct{}{}

var allowedResourceTypes = map[string]struct{}{
	"gamemode": struct{}{},
	"map":      struct{}{},
	"script":   struct{}{},
	"misc":     struct{}{},
}

func init() {
	secpath := os.Getenv("MTAHUB_SECURITY_PLUGIN")
	if secpath == "" {
		return
	}

	var err error
	securityPlugin, err = plugin.Open(secpath)
	if err != nil {
		panic(err.Error())
	}

	beS, err := securityPlugin.Lookup("BannedExtensions")
	if err != nil {
		panic(err.Error())
	}

	beMp, ok := beS.(*map[string]struct{})
	if !ok {
		panic("could not load BannedExtensions")
	}

	bannedExtensions = *beMp
}

// CheckResourceZip decodes an input zip and checks if the resource is ok
func CheckResourceZip(f io.ReaderAt, size int64) (meta *XmlMeta, ok bool, reason string, err error) {
	r, err := zip.NewReader(f, size)
	if err != nil {
		return
	}

	var metaZipfile *zip.File
	fmt.Printf("bannedExtensions %#v\n", bannedExtensions)

	for _, file := range r.File {
		fileInfo := file.FileInfo()

		filename := fileInfo.Name()
		if filename == "meta.xml" && !fileInfo.IsDir() {
			metaZipfile = file
			continue
		}

		_, isBanned := bannedExtensions[fpath.Ext(filename)]
		if isBanned {
			return nil, false, fmt.Sprintf("contains blocked file %#v", filename), nil
		}
	}

	if metaZipfile == nil {
		return nil, false, "missing meta.xml file", nil
	}

	metafile, err := metaZipfile.Open()
	if err != nil {
		return nil, false, "", errors.Wrap(err, "could not open meta.xml")
	}
	defer metafile.Close()

	meta, errReason := checkMeta(metafile)
	if meta == nil {
		return nil, false, errReason, nil
	}

	fmt.Printf(`ok {"version": "%s", "type": "%s", "name": "%s"}`+"\n", meta.Infos[0].Version, meta.Infos[0].Type, meta.Infos[0].Name)
	fmt.Printf("%#v\n", meta)

	return meta, true, "", nil
}

func checkMeta(file io.ReadCloser) (meta *XmlMeta, reason string) {
	meta = &XmlMeta{}

	decoder := xml.NewDecoder(file)
	err := decoder.Decode(meta)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode meta.xml").Error()
	}

	// Require exactly 1 info node
	if len(meta.Infos) != 1 {
		return nil, "meta.xml should have exactly 1 info field"
	}

	// Check <info type
	info := meta.Infos[0]
	if info.Type == "" {
		return nil, "meta.xml is missing the 'type' field"
	} else if _, ok := allowedResourceTypes[info.Type]; !ok {
		return nil, "meta.xml has an invalid 'type' field"
	}

	// Require <info version=
	if info.Version == "" {
		return nil, "meta.xml is missing the version field for <info>"
	}

	// Require that version is well formed
	if ok, _ := regexp.MatchString(`^(\d\.\d\.\d|\d\.\d|\d)$`, info.Version); !ok {
		return nil, "meta.xml contains a malformed version field (should be in the form X, X.X, or X.X.X)"
	}

	return
}
