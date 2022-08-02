package socket

import(
	"io"
	"os"
	"fmt"
	"net"
)

func Client() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("connection error: %v\n", err.Error())
		return
	}
	go func() {io.Copy(conn, os.Stdin)}()
	io.Copy(os.Stdout, conn)
}