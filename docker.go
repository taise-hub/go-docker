package main

import (
	"time"
	"fmt"
	"io"
	"os"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/docker/docker/pkg/stdcopy"
)


func main() {
	conf := &container.Config{
		AttachStdin: true, 
		AttachStdout: true,
		AttachStderr: true,
		Tty: true,
		Image: "alpine",
	}

	hconf := &container.HostConfig {
		AutoRemove: true,
	}

	spec := &specs.Platform{
		Architecture: "amd64",
		OS: "linux",
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10000* time.Millisecond)
	defer cancel()

	fmt.Println("[+] Create Container.")
	ccb, err := cli.ContainerCreate(ctx, conf, hconf, nil, spec, "test-container")
	if err != nil {
		panic(err)
	}

	fmt.Println("[+] Start Container.")
	if err = cli.ContainerStart(ctx, ccb.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	econf := types.ExecConfig{
		AttachStdin: true, 
		AttachStdout: true,
		AttachStderr: true,
		Tty: true,
		Cmd: []string{"/bin/sh"},
	}

	exec, err := cli.ContainerExecCreate(ctx, "test-container", econf)
	if err != nil {
		panic(err)
	}

	fmt.Println("[+] Exec command on Container.")
	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	go func() { _, _ = io.Copy(resp.Conn, os.Stdin) }()
	stdcopy.StdCopy(os.Stdout, os.Stderr, resp.Conn)
}