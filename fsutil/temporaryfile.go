package fsutil

import (
	"fmt"
	"os"
	"io/ioutil"
)

// ErrPath is returned as the panic() context when a Must function is invoked
// as a wrapped error..
type ErrTemporaryFile struct {
	message string
	dir string
	prefix string
	err error
}

func newErrTemporaryFile(message string, dir string, prefix string, err error) *ErrTemporaryFile {
	return &ErrTemporaryFile{
		message: message,
		prefix: prefix,
		dir: dir,
		err: err,
	}
}

func (this ErrTemporaryFile) WrappedErrors() []error {
	return []error{this.err}
}

func (this ErrTemporaryFile) Error() string {
	return this.message
}

func (this ErrTemporaryFile) Dir() string {
	return this.dir
}

func (this ErrTemporaryFile) Prefix() string {
	return this.prefix
}

// Creates a temporary file with the given content and returns the filename
func TemporaryFileWithContent(dir string, prefix string, content []byte, mode os.FileMode) (string, error) {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return "", newErrTemporaryFile("Failed creating temporary file:", dir, prefix, err)
	}
	filename := f.Name()
	f.Close()
	if err := ioutil.WriteFile(f.Name(), content, mode); err != nil {
		return "", newErrTemporaryFile("Failed writing data to temporary file:", dir, prefix, err)
	}
	return filename, nil
}

// Creates a temporary file with the given content and returns the filename.
// Panics if any error occurs.
func MustTemporaryFileWithContent(dir string, prefix string, content []byte, mode os.FileMode) string {
	fname, err := TemporaryFileWithContent(dir, prefix, content, mode)
	if err != nil {
		panic(err)
	}
	return fname
}