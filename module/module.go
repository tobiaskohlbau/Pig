package module

import (
	"os"
	"bufio"
	"io/ioutil"
	"fmt"
	"os/exec"
	"io"
	"runtime"
)

const (
	GIT_LINUX = "/usr/bin/git"
	GIT_WINDOWS = "C:\\Program Files (x86)\\Git\\bin\\git.exe"
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

func (m *Repo) fetch(path string) {
	var GIT string
	if runtime.GOOS=="windows"{
   		GIT = GIT_WINDOWS
	} else {
		GIT = GIT_LINUX
	}
	cmd := exec.Command(GIT, "fetch", "--progress", "origin")
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

func (m *Repo) pull(path string) {
	var GIT string
	if runtime.GOOS=="windows"{
   		GIT = GIT_WINDOWS
	} else {
		GIT = GIT_LINUX
	}
	cmd := exec.Command(GIT, "pull", "--progress", "origin", m.Branch)
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
	var GIT, remote string
	if runtime.GOOS=="windows"{
   		GIT = GIT_WINDOWS
   		remote = "/" + m.Remote 
	} else {
		GIT = GIT_LINUX
		remote = m.Remote
	}
	cmd := exec.Command(GIT, "clone", "--progress", remote, path)
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

func (m *Repo) checkout(path string) {
	var GIT string
	if runtime.GOOS=="windows"{
   		GIT = GIT_WINDOWS 
	} else {
		GIT = GIT_LINUX
	}
	cmd := exec.Command(GIT, "checkout", m.Branch)
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
		fmt.Printf("Repository %s exist fetch data right now\n", m.Name)
		m.fetch(pwd + "/" + folder + "/" + m.Path)
		fmt.Printf("Repository %s exist sync it right now\n", m.Name)
		m.pull(pwd + "/" + folder + "/" + m.Path)
	}

	fmt.Printf("Repository %s exist checkout %s branch now\n", m.Name, m.Branch)
	m.checkout(pwd + "/" + folder + "/" + m.Path)

	tasks<-1
}