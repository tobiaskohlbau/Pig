package module

import (
	"os"
	"bufio"
	"io/ioutil"
	"fmt"
	"os/exec"
	"io"
)

const (
	GIT = "/usr/bin/git"
)

type Repo struct {
	Name string
	Path string
	Remote string
	Branch string
}

type RepoList []Repo

func readToString(f *os.File) string {
	defer f.Close()
	reader := bufio.NewReader(f)
	contents, _ := ioutil.ReadAll(reader)
	return string(contents)
}

func Parse(url string) RepoList {
	file, err := os.Open(url)
	if err != nil {
		panic(err)
	}
	inputfile := readToString(file)
	tmp := Tree{Name:"testTree", text: inputfile, lex: lex("test", inputfile)}
	tmp.parse()
	repoList := RepoList{}
	for _,node := range tmp.Root.Nodes {
		switch node := node.(type) {
		case *ModuleNode:
			repoList = append(repoList,
				Repo{Name: node.Name.String(), Path: node.Path.String(), Remote: node.Remote.String(), Branch: node.Branch.String()})
		}
	}
	return repoList
}

func (m *Repo) pull(path string) {
	cmd := exec.Command(GIT, "pull", "-v", "--progress", "origin", m.Branch)
	cmd.Dir = path
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	cmd.Wait()
}

func (m *Repo) clone(path string) {
	cmd := exec.Command(GIT, "clone", "-v", "--progress", m.Remote, path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	cmd.Wait()
}

func (m *Repo) Sync(folder string, tasks chan int) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(pwd + "/" + folder + "/" + m.Path + "/.git"); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Repository %s doesn't exist clone it right now\n", m.Name)
			m.clone(pwd + "/" + folder + "/" + m.Path)
		}
	} else {
		fmt.Printf("Repository %s exist sync it right now\n", m.Name)
		m.pull(pwd + "/" + folder + "/" + m.Path)
	}
	tasks<-1
}