package audiometa

import (
	"io"
)

// MockReaderSeeker implements io.ReadSeeker for testing.
type MockReaderSeeker struct {
	data []byte
	pos  int64
}

func (r *MockReaderSeeker) Read(p []byte) (n int, err error) {
	n = copy(p, r.data[r.pos:])
	r.pos += int64(n)
	if r.pos >= int64(len(r.data)) {
		err = io.EOF
	}
	return
}

func (r *MockReaderSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.pos = offset
	case io.SeekCurrent:
		r.pos += offset
	case io.SeekEnd:
		r.pos = int64(len(r.data)) + offset
	}
	return r.pos, nil
}
