package docker

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

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
)

func Server() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10000* time.Millisecond)
	defer cancel()


	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}
		go handle(ctx, cli, conn)
	}

}

func handle(ctx context.Context, cli *client.Client, conn net.Conn) {
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
	defer conn.Close()
	go func () { io.Copy(conn, resp.Conn) }()
	io.Copy(resp.Conn, conn)
}