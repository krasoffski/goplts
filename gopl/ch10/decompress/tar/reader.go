package zip

import (
	"archive/tar"
	"io"
	"os"

	"github.com/krasoffski/goplts/gopl/ch10/decompress"
)

// NewReader create a zip reader from the File.
func NewReader(f *os.File) ([]*decompress.Entry, error) {
	tr := tar.NewReader(f)

	// storing compressed file readers in the slice of interfaces
	entries := make([]*decompress.Entry, 0)
	// zip file contains a number of individual files/readers
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		var isDir bool
		if header.Typeflag == tar.TypeDir {
			isDir = true
		}

		// NOTE: does not work due to tar reader internal state.
		e := decompress.Entry{
			Reader: tr,
			Name:   header.Name,
			IsDir:  isDir,
			Mode:   os.FileMode(header.Mode),
		}
		entries = append(entries, &e)
	}
	return entries, nil
}

func init() {
	decompress.RegisterFormat("ustar", 257, NewReader)
}
