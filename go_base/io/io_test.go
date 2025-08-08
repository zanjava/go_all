package io_test

import (
	"fmt"
	"go/base/io"
	"testing"
	"time"
)

func TestWriteFile(t *testing.T) {
	io.WriteFile()
}

func TestWriteFileWithBuffer(t *testing.T) {
	io.WriteFileWithBuffer()
}

func TestReadFile(t *testing.T) {
	io.ReadFile()
}

func TestReadFileWithBuffer(t *testing.T) {
	io.ReadFileWithBuffer()
}

// TestBufferedFileWriter tests the BufferedFileWriter functionality.
func TestBufferedFileWriter(t *testing.T) {
	t1 := time.Now()
	io.WriteFileDirectly("../../data/no_buffer.txt")
	t2 := time.Now()
	io.WriteWithBuffer("../../data/with_buffer.txt")
	t3 := time.Now()
	fmt.Printf("直接写文件耗时: %dms\n", t2.Sub(t1).Milliseconds())
	fmt.Printf("使用缓冲区写文件耗时: %dms\n", t3.Sub(t2).Milliseconds())
}

func TestCreateFile(t *testing.T) {
	io.CreateFile("../../data/test.txt")
}
func TestWalkDir(t *testing.T) {
	if err := io.WalkDir("../../data"); err != nil {
		t.Errorf("WalkDir failed: %v", err)
	}
}

func TestSplitFile(t *testing.T) {
	// Split the file into 4 parts
	infile := "../../data/img/output.png"
	n := 4
	io.SplitFile(infile, "../../data/img/图像分割", n)
}

func TestMergeFile(t *testing.T) {
	// Merge the split files back into one
	outfile := "../../data/img/merged_output.png"
	io.MergeFile("../../data/img/图像分割", outfile)
}

func TestLimitedReader(t *testing.T) {
	io.LimitedReader()
}

func TestMultiReader(t *testing.T) {
	io.MultiReader()
}

func TestMultiWriter(t *testing.T) {
	io.MultiWriter()
}

func TestTeeReader(t *testing.T) {
	io.TeeReader()
}

func TestPipeIO(t *testing.T) {
	io.PipeIO()
}

func TestCopyFile(t *testing.T) {
	io.CopyFile("../../data/no_buffer.txt", "../../data/no_buffer_copy.txt")
}

func TestCompress(t *testing.T) {
	io.Compress("../../data/no_buffer.txt", "../../data/no_buffer.txt.gz", io.ZLIB)
}

func TestDecompress(t *testing.T) {
	io.Decompress("../../data/no_buffer.txt.gz", "../../data/no_buffer1.txt", io.ZLIB)
}

func TestJsonSerialize(t *testing.T) {
	io.JsonSerialize()
}

func TestNewLogger(t *testing.T) {
	logger := io.NewLogger("../../data/biz.log")
	io.Log(logger)
}

func TestNewSLogger(t *testing.T) {
	logger := io.NewSLogger("../../data/sbiz.log")
	io.SLog(logger)
}

func TestSysCall(t *testing.T) {
	io.SysCall()
}

func TestUseRegex(t *testing.T) {
	io.UseRegex()
}
