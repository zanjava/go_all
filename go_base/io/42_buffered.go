package io

import "os"

type BufferedFileWriter struct {
	fout           *os.File // 文件句柄
	buffer         []byte   // 缓冲区
	bufferEndIndex int      // 缓冲区偏移量
}

func NewBufferedFileWriter(fout *os.File, bufferSize int) *BufferedFileWriter {
	return &BufferedFileWriter{
		fout:           fout,
		buffer:         make([]byte, bufferSize),
		bufferEndIndex: 0,
	}
}

func (w *BufferedFileWriter) Write(data []byte) (int, error) {
	dataLen := len(data)
	if dataLen == 0 {
		return 0, nil
	}

	// 如果数据长度大于缓冲区大小，则直接写入文件
	if dataLen >= len(w.buffer) {
		w.Flush()
		if _, err := w.fout.Write(data); err != nil {
			return 0, err
		}
	}

	// 如果缓冲区剩余空间不足以容纳新数据，则先刷新缓冲区
	if w.bufferEndIndex+dataLen > len(w.buffer) {
		if err := w.Flush(); err != nil {
			return 0, err
		}
	}

	copy(w.buffer[w.bufferEndIndex:], data)
	w.bufferEndIndex += dataLen
	return dataLen, nil
}

func (f *BufferedFileWriter) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}

func (w *BufferedFileWriter) Flush() error {
	if w.bufferEndIndex == 0 {
		return nil
	}

	if _, err := w.fout.Write(w.buffer[:w.bufferEndIndex]); err != nil {
		return err
	}

	w.bufferEndIndex = 0
	return nil
}

func (w *BufferedFileWriter) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	return w.fout.Close()
}

var data string = "This is a test line.\n"

// 直接写文件
func WriteFileDirectly(outFile string) {
	fout, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	for i := 0; i < 100000; i++ {
		fout.WriteString(data)
	}
}

// 使用缓冲区写文件
func WriteWithBuffer(outFile string) {
	fout, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	writer := NewBufferedFileWriter(fout, 4096)
	defer writer.Close()
	for i := 0; i < 100000; i++ {
		writer.WriteString(data)
	}
}
