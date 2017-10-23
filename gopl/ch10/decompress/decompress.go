package decompress

import (
	"errors"
	"io"
	"os"
)

// ErrFormat is unsupported format error
var ErrFormat = errors.New("decompress: unknown format")

// Signature contains the magic string and the offset for look up.
type Signature struct {
	Magic  string
	Offset int
}

// Decompressor creates a reader form the file parameter
type Decompressor func(*os.File) (MultiPartFile, error) // not memory efficient

type format struct {
	signature    Signature
	decompressor Decompressor
}

var formats []*format

// sniff determines the format of r's data.
func sniff(f *os.File) (*format, error) {
	// get magic from file
	for _, format := range formats {
		b := make([]byte, len(format.signature.Magic))
		_, err := f.ReadAt(b, int64(format.signature.Offset))
		if err != nil {
			return nil, err
		}
		// reset file position
		_, err = f.Seek(0, os.SEEK_SET)
		if err != nil {
			return nil, err
		}
		// return true if magic matches
		if format.signature.Magic == string(b) {
			return format, nil
		}
	}
	return nil, ErrFormat
}

// Entry represents basic compressed file within archive.
type Entry struct {
	Reader io.Reader
	Header string
	IsDir  bool
	Mode   os.FileMode
}

// MultiPartFile describes minimal interface for interating over compressed files.
// Using the similar approach as for tar module.
type MultiPartFile interface {
	Next() (*Entry, error)
}

// RegisterFormat registers decompressor for the magic number.
func RegisterFormat(magic string, offset int, d Decompressor) {
	f := format{
		signature:    Signature{Magic: magic, Offset: offset},
		decompressor: d,
	}
	formats = append(formats, &f)
}

// NewReader creates a decompressing reader.
func NewReader(file *os.File) (MultiPartFile, error) {
	f, err := sniff(file)
	if err != nil {
		return nil, err
	}
	return f.decompressor(file)
}
