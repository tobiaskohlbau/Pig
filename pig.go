package main

import (
	"tu-chemnitz.de/dst/Pig/module"
	"flag"
	"os"
	"fmt"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: pig remote_url folder\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("no enough arguments")
		os.Exit(1)
	}	

	modules := module.Parse(args[0])
	fmt.Println(args[1])
	tasks := make(chan int)
	for _,repo := range modules {
		go func(repo module.Repo) {
			repo.Sync(args[1], tasks)
		}(repo)
	}
	for i := 0; i < len(modules); i++ {
		<-tasks
	}
}