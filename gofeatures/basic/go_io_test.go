package basic

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestReadIo(t *testing.T) {
	r := strings.NewReader("Hello, ioReader")
	// 以8Bytes大小作为buffer进行读取
	b := make([]byte, 8)
	for {
		// read直接传入用于接收的buffer数组，n是指这次用到了多少（最后一次可能<size)
		n, err := r.Read(b)
		fmt.Printf("n = %v, err = %v, b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		// 对于Reader来说，reach到EOR是当作异常给出
		if err == io.EOF {
			break
		}
	}
}
