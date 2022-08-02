package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"math/rand"

	"github.com/gorilla/websocket"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

var (
	conf = &container.Config {
		AttachStdin: true, 
		AttachStdout: true,
		AttachStderr: true,
		Tty: true,
		Image: "alpine",
	}
	hconf = &container.HostConfig {
		AutoRemove: true,
	}
	spec = &specs.Platform{
		Architecture: "amd64",
		OS: "linux",
	}
	econf = types.ExecConfig{
		AttachStdin: true, 
		AttachStdout: true,
		AttachStderr: true,
		Tty: true,
		Cmd: []string{"/bin/sh"},
	}
	upgrader = websocket.Upgrader{
   		ReadBufferSize:  1024,
  		WriteBufferSize: 1024,
	}
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello, %v\n", req.FormValue("name"))
}

func wsHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	go handle(context.Background(), conn.UnderlyingConn())
}

func handle(ctx context.Context, conn net.Conn) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("[+] Create Container.")
	name := strconv.Itoa(rand.Int())
	createdBody, err := cli.ContainerCreate(ctx, conf, hconf, nil, spec, name)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("[+] Start Container.")
	if err = cli.ContainerStart(ctx, createdBody.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
		return
	}
	log.Println("[+] Exec command on Container.")
	exec, err := cli.ContainerExecCreate(ctx, name, econf)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()
	go func() { io.Copy(conn, resp.Conn) }()
	io.Copy(resp.Conn, conn)
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":80", nil)
}