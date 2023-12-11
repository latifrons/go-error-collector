package main

import (
	"github.com/latifrons/goerrorcollector/consumer/cmd"
	_ "google.golang.org/grpc"
)

func main() {
	cmd.Execute()
}
