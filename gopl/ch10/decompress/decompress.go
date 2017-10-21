package decompress

import (
	"bufio"
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
type Decompressor func(*os.File) ([]*Entry, error) // not memory efficient

type format struct {
	signature    Signature
	decompressor Decompressor
}

var formats []*format

// Match reports whether magic matches b. Magic may contain "?" wildcards.
func match(magic string, b []byte) bool {
	if len(magic) != len(b) {
		return false
	}
	for i, c := range b {
		if magic[i] != c && magic[i] != '?' {
			return false
		}
	}
	return true
}

func sniff2(f *os.File) (*format, error) {
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

// Sniff determines the format of r's data.
func sniff(r io.Reader) (*format, error) {
	reader := bufio.NewReader(r)
	for _, f := range formats {
		b, err := reader.Peek(len(f.signature.Magic))
		if err == nil && match(f.signature.Magic, b) {
			return f, nil
		}
	}
	return nil, ErrFormat
}

// Entry represents basic compressed file within archive.
type Entry struct {
	Reader io.Reader
	Name   string
	IsDir  bool
	Mode   os.FileMode
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
func NewReader(file *os.File) ([]*Entry, error) {
	f, err := sniff2(file)
	if err != nil {
		return nil, err
	}
	return f.decompressor(file)
}
