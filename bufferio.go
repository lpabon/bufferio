// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bufferio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

var (
	ErrOverrun = errors.New("buffer overrun")
	ErrEOF     = errors.New("end of file")
)

type BufferIO struct {
	buf []byte
	off int64
}

func NewBufferIO(b []byte) *BufferIO {
	return &BufferIO{buf: b}
}

func NewBufferIOMake(nbytes int) *BufferIO {
	return &BufferIO{buf: make([]byte, nbytes)}
}

func (b *BufferIO) WriteAt(p []byte, off int64) (n int, err error) {
	if off >= b.Size() {
		return 0, ErrOverrun
	}
	bytes_copied := copy(b.buf[off:], p)
	return bytes_copied, nil
}

func (b *BufferIO) Write(p []byte) (n int, err error) {
	n, err = b.WriteAt(p, b.off)
	if err == nil {
		b.off += int64(n)
	}
	return n, err
}

func (b *BufferIO) WriteData(order binary.ByteOrder, data interface{}) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, order, data)
	if err != nil {
		return err
	}
	_, err = b.Write(buf.Bytes())
	return err
}

func (b *BufferIO) WriteDataLE(data interface{}) error {
	return b.WriteData(binary.LittleEndian, data)
}

func (b *BufferIO) WriteDataBE(data interface{}) error {
	return b.WriteData(binary.BigEndian, data)
}

func (b *BufferIO) ReadAt(p []byte, off int64) (n int, err error) {
	if off >= b.Size() {
		return 0, ErrEOF
	}
	bytes_copied := copy(p, b.buf[off:])
	return bytes_copied, nil
}

func (b *BufferIO) Read(p []byte) (n int, err error) {
	n, err = b.ReadAt(p, b.off)
	if err == nil {
		b.off += int64(n)
	}
	return n, err
}

func (b *BufferIO) ReadData(order binary.ByteOrder, data interface{}) error {
	buf := bytes.NewReader(b.buf[b.off:]) // this can probably be done with BufferIO
	return binary.Read(buf, order, data)
}

func (b *BufferIO) ReadDataLE(data interface{}) error {
	return b.ReadData(binary.LittleEndian, data)
}
func (b *BufferIO) ReadDataBE(data interface{}) error {
	return b.ReadData(binary.BigEndian, data)
}

func (b *BufferIO) Seek(offset int64, whence int) (int64, error) {
	var position int64
	switch whence {
	case os.SEEK_SET:
		position = offset
	case os.SEEK_CUR:
		position = b.off + offset
	case os.SEEK_END:
		return 0, ErrOverrun
	default:
		return 0, errors.New("invalid whence")
	}

	if position >= b.Size() {
		return 0, ErrOverrun
	}
	if position < 0 {
		return 0, errors.New("negative position")
	}

	b.off = position
	return position, nil
}

func (b *BufferIO) Bytes() []byte {
	return b.buf
}

func (b *BufferIO) Reset() {
	b.off = 0
}

func (b *BufferIO) Size() int64 {
	return int64(len(b.buf))
}
