package zip

import (
	"archive/zip"
	"os"

	"github.com/krasoffski/goplts/gopl/ch10/decompress"
)

// NewReader create a zip reader from the File.
func NewReader(f *os.File) ([]*decompress.Entry, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return nil, err
	}
	// storing compressed file readers in the slice of interfaces
	entries := make([]*decompress.Entry, 0, len(zr.File))
	// zip file contains a number of individual files/readers
	for _, zf := range zr.File {
		rc, err := zf.Open() // entry is represent compressed file
		if err != nil {
			return nil, err
		}
		// NOTE: need verify ability to skip reader creation for directories.
		e := decompress.Entry{
			Reader: rc,
			Header: zf.Name,
			IsDir:  zf.FileInfo().IsDir(),
			Mode:   zf.Mode(),
		}
		entries = append(entries, &e)
	}
	return entries, nil
}

func init() {
	decompress.RegisterFormat("PK", 0, NewReader)
}
