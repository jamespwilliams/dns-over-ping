package main

import "github.com/jamespwilliams/doicmp"

func main() {
	panic(doicmp.NewServer().Serve())
}
