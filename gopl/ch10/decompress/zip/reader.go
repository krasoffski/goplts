package zip

import (
	"archive/zip"
	"io"
	"os"

	"github.com/krasoffski/goplts/gopl/ch10/decompress"
)

type zipReader struct {
	reader *zip.Reader
	index  int
}

func (zr *zipReader) Next() (*decompress.Entry, error) {
	if zr.index >= len(zr.reader.File) {
		return nil, io.EOF
	}

	zf := zr.reader.File[zr.index]
	zr.index++ // index for next iteration, increasing here for next file attempt

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
	return &e, nil
}

// NewReader create a zip reader from the File.
func NewReader(f *os.File) (decompress.MultiPartFile, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return nil, err
	}
	return &zipReader{reader: zr}, nil
}

func init() {
	decompress.RegisterFormat("PK", 0, NewReader)
}
