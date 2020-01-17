package resource

import (
	"io"
	"io/ioutil"
)

// CheckResourceZip decodes an input zip and checks if the resource is ok
func CheckResourceZip(f io.Reader) (ok bool, reason string, err error) {
	_, err = ioutil.ReadAll(f)
	if err != nil {
		return
	}

	return true, "unimplemented", nil
}
