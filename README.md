# go-docker

## Environment

you need to have Go's runtime installed. Please install it by referring to the [Official Site] (https://golang.org/) in advance.

## What is this

This program is a simple example of connecting to a docker container over the network.

There are two connection examples. The first is a connection using net.Conn.

The second is a connection using [gorilla](https://github.com/gorilla/websocket) (websocket).