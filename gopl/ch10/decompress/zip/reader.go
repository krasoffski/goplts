package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	"github.com/krasoffski/goplts/gopl/ch10/decompress"
)

// NewReader create a zip reader from the File.
func NewReader(f *os.File) ([]io.ReadCloser, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return nil, err
	}
	readers := make([]io.ReadCloser, 0, len(zr.File))
	for _, zf := range zr.File {
		fmt.Fprintln(os.Stderr, zf.Name)
		r, err := zf.Open() // zip internal file reader
		if err != nil {
			return nil, err
		}
		readers = append(readers, r)
	}
	return readers, nil
}

func init() {
	decompress.RegisterFormat("PK", 0, NewReader)
}
