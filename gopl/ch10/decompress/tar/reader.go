package tar

import (
	"archive/tar"
	"os"

	"github.com/krasoffski/goplts/gopl/ch10/decompress"
)

type tarReader struct {
	reader *tar.Reader
}

func (tr *tarReader) Next() (*decompress.Entry, error) {
	header, err := tr.reader.Next()
	if err != nil {
		return nil, err
	}

	var isDir bool
	if header.Typeflag == tar.TypeDir {
		isDir = true
	}

	e := decompress.Entry{
		Reader: tr.reader,
		Header: header.Name,
		IsDir:  isDir,
		Mode:   os.FileMode(header.Mode),
	}
	return &e, nil
}

// NewReader create a tar reader from the File.
func NewReader(f *os.File) (decompress.MultiPartFile, error) {
	tr := tar.NewReader(f)
	return &tarReader{reader: tr}, nil
}

func init() {
	decompress.RegisterFormat("ustar", 257, NewReader)
}
