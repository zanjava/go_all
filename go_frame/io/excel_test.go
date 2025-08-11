package io_test

import (
	"go/frame/io"
	"testing"
)

func TestExcel(t *testing.T) {
	io.ReadWriteExcel("../data/学生信息表.xlsx")
}

// go test -v ./io -run=^TestExcel$ -count=1
