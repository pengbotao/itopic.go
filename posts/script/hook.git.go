package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	host     = "127.0.0.1:8002"
	hookPath = "/data/www/gopath/src/"
)

func hook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	project, ok := r.Form["project"]
	if !ok {
		w.Write([]byte("Hello World"))
		return
	}
	path := hookPath + strings.Join(project, "")
	str := fmt.Sprintf("cd %s && git checkout master && git reset --hard HEAD && git pull", path)
	println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + str)
	name := "/bin/sh"
	if strings.Compare(runtime.GOOS, "windows") == 0 {
		name = "C:/Program Files/Git/git-bash.exe"
	}
	go func() {
		out, err := exec.Command(name, "-c", str).Output()
		if err != nil {
			println(err.Error())
			return
		}
		println(string(out))
		return
	}()
	fmt.Fprintln(w, "Hello World.")
}

func main() {
	http.HandleFunc("/", hook)
	fmt.Println("The topic server is running (" + runtime.GOOS + ") at http://" + host)
	fmt.Println("Quit the server with Control-C")
	http.ListenAndServe(host, nil)
}
