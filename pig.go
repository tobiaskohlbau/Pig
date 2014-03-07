package main

import (
	"tu-chemnitz.de/dst/Pig/module"
	"flag"
	"os"
	"fmt"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: Pig folder\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("no folder given")
		os.Exit(1)
	}	

	modules := module.Parse("/home/tobias/test.dst")
	tasks := make(chan int)
	for _,repo := range modules {
		go func(repo module.Repo) {
			repo.Sync(args[0], tasks)
		}(repo)
	}
	for i := 0; i < len(modules); i++ {
		<-tasks
	}
}