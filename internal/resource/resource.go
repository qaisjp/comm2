package resource

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	fpath "path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

var bannedExtensions = map[string]struct{}{
	".exe": struct{}{},
	".com": struct{}{},
	".bat": struct{}{},
}

var allowedResourceTypes = map[string]struct{}{
	"gamemode": struct{}{},
	"map":      struct{}{},
	"script":   struct{}{},
	"misc":     struct{}{},
}

// CheckResourceZip decodes an input zip and checks if the resource is ok
func CheckResourceZip(f io.ReaderAt, size int64) (ok bool, reason string, err error) {
	r, err := zip.NewReader(f, size)
	if err != nil {
		return
	}

	var metaZipfile *zip.File

	for _, file := range r.File {
		fileInfo := file.FileInfo()

		filename := fileInfo.Name()
		if filename == "meta.xml" && !fileInfo.IsDir() {
			metaZipfile = file
			continue
		}

		_, isBanned := bannedExtensions[fpath.Ext(filename)]
		if isBanned {
			return false, fmt.Sprintf("contains blocked file %#v", filename), nil
		}
	}

	if metaZipfile == nil {
		return false, "missing meta.xml file", nil
	}

	metafile, err := metaZipfile.Open()
	if err != nil {
		return false, "", errors.Wrap(err, "could not open meta.xml")
	}
	defer metafile.Close()

	meta, errReason := checkMeta(metafile)
	if meta == nil {
		return false, errReason, nil
	}

	fmt.Printf(`ok {"version": "%s", "type": "%s", "name": "%s"}`+"\n", meta.Infos[0].Version, meta.Infos[0].Type, meta.Infos[0].Name)
	fmt.Printf("%#v\n", meta)

	return true, "", nil
}

func checkMeta(file io.ReadCloser) (meta *xmlMeta, reason string) {
	meta = &xmlMeta{}

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
