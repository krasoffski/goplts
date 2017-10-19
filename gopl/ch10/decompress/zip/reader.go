package zip

import (
	"archive/zip"
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
	zipFiles := make([]io.ReadCloser, len(zr.File))
	for _, zipFile := range zr.File {
		fileReader, err := zipFile.Open()
		if err != nil {
			return nil, err
		}
		zipFiles = append(zipFiles, fileReader)
	}
	return zipFiles, nil
}

func init() {
	decompress.RegisterFormat("PK", 0, NewReader)
	// decompress.RegisterFormat(archive.Magic{"PK", 0}, NewReader)
}
