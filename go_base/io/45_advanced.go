package io

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// 受限的reader
func LimitedReader() {

	reader := strings.NewReader("Hello, World!")
	// 创建一个限制大小的读取器
	limitedReader := io.LimitReader(reader, 3) // 只读取前3字节
	// 进行读取操作
	data := make([]byte, 5)
	n, err := limitedReader.Read(data)
	if err != nil {
		return
	}
	fmt.Printf("读取到的数据: %s\n", data[:n])

	if _, err := limitedReader.Read(data); err == io.EOF {
		fmt.Println("已到达读取限制或结束")
	}

}

// multiReader 可以将多个文件合并为一个文件
func MultiReader() {
	reader1 := strings.NewReader("Hello, ")
	reader2 := strings.NewReader("World!")
	// 创建一个多读取器
	multiReader := io.MultiReader(reader1, reader2)
	io.Copy(os.Stdout, multiReader) // 将内容写入标准输出
}

// multiWriter 可以将一条日志输出到多个文件
func MultiWriter() {
	var (
		w1 bytes.Buffer
		w2 bytes.Buffer
	)
	multiWriter := io.MultiWriter(&w1, &w2)
	multiWriter.Write([]byte("Hello, World!"))

	// 输出写入的内容
	fmt.Println("w1:", w1.String())
	fmt.Println("w2:", w2.String())
}

// teeReader
func TeeReader() {
	reader := strings.NewReader("Hello, World!")
	// 创建一个缓冲区
	var buf bytes.Buffer
	// 创建一个 TeeReader，将读取的数据同时写入 buf
	teeReader := io.TeeReader(reader, &buf)

	io.Copy(os.Stdout, teeReader)
	fmt.Printf("缓冲区内容: %s\n", buf.String())
}

// pipeIO
func PipeIO() {
	// 创建一个管道
	reader, writer := io.Pipe()

	// 启动一个 goroutine 来写入数据
	go func() {
		defer writer.Close()
		writer.Write([]byte("Hello from PipeIO!"))
	}()

	// 从管道中读取数据
	readerData := make([]byte, 1024)
	n, err := reader.Read(readerData)
	if err != nil {
		fmt.Println("读取错误:", err)
		return
	}
	fmt.Println("从管道读取的数据:", string(readerData[:n]))
}
