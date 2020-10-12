package database

import (
	"fmt"
)

type DuplicateFileError struct {
	path string
	err  error
}

func (d DuplicateFileError) Error() string {
	return fmt.Sprintf("File already on disk -> %s :%v", d.path, d.err)
}
func NewDuplicateFileError(path string, err error) *DuplicateFileError {
	return &DuplicateFileError{
		path: path,
		err:  err,
	}
}
