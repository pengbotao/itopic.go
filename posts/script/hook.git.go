package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// go run main.go --host :8088 --path /data/wwwroot --key 123456

var (
	host     string
	hookPath string
	hookKey  string
	hookCmd  string
)

func main() {
	flag.StringVar(&host, "host", "127.0.0.1:8088", "host")
	flag.StringVar(&hookPath, "path", "", "hook path")
	flag.StringVar(&hookKey, "key", "", "hook key")
	flag.StringVar(&hookCmd, "cmd", "/bin/sh", "sh command")
	flag.Parse()
	if hookPath == "" {
		fmt.Println("please specified hook path.")
		os.Exit(0)
	}
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
		return
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		k, ok := r.Form["k"]
		if hookKey != "" {
			if !ok || strings.Compare(strings.Join(k, ""), hookKey) != 0 {
				return
			}
		}
		str := fmt.Sprintf("cd %s && git checkout master && git reset --hard HEAD && git pull", hookPath)
		println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + str)
		go func() {
			out, err := exec.Command(hookCmd, "-c", str).Output()
			if err != nil {
				println(err.Error())
				return
			}
			println(string(out))
			return
		}()
		fmt.Fprintln(w, "Finish.")
	})
	fmt.Println("The hook server is running at http://" + host)
	fmt.Println("Quit the server with Control-C")
	http.ListenAndServe(host, nil)
}
