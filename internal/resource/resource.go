package resource

import (
	"archive/zip"
	"io"
	"strconv"
)

// CheckResourceZip decodes an input zip and checks if the resource is ok
func CheckResourceZip(f io.ReaderAt, size int64) (ok bool, reason string, err error) {
	r, err := zip.NewReader(f, size)
	if err != nil {
		return
	}

	return false, strconv.Itoa(len(r.File)) + " files", nil
}
