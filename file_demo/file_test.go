package file_demo

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	// 获得工作目录
	fmt.Println(os.Getwd())
	
	file, err := os.Open("testdata/my_file.txt")
	require.NoError(t, err)
	data := make([]byte, 64)
	n, err := file.Read(data)
	fmt.Println(n)
	require.NoError(t, err)
	
	n, err =  file.Write([]byte("hello world"))
	fmt.Println(n)
	// bad file descriptor 不可写
	fmt.Println(err)
	require.Error(t, err)
	file.Close()
	
	file, err = os.OpenFile("testdata/my_file.txt", os.O_APPEND | os.O_WRONLY | 
		os.O_CREATE, os.ModeAppend)
	require.NoError(t, err)
	n, err = file.WriteString("hello")
	fmt.Println(n)
	require.NoError(t, err)
	file.Close()

	file, err = os.Create("testdata/my_file_copy.txt")
	require.NoError(t, err)
	n, err = file.WriteString("hello, world")
	fmt.Println(n)
	require.NoError(t, err)
}