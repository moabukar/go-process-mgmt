package main

import (
	"fmt"
	"github.com/ahmetb/pstree"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"github.com/google/uuid"
	"syscall"
)

var proc *os.Process
var incarnation string

func main() {
	incarnation = uuid.New().String()
	http.HandleFunc("/kill", kill)
	http.HandleFunc("/start", kill)
	http.HandleFunc("/run", run)
	http.HandleFunc("/ps", ps)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func home(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "container incarnation id: %s\n\n", incarnation)
	fmt.Fprintln(w, "/run")
	fmt.Fprintln(w, "/ps")
	fmt.Fprintln(w, "/kill (optional ?pid=, otherwise kills subprocess’s pgrp via negative PID)")
}
func run(w http.ResponseWriter, _ *http.Request) {
	if proc != nil {
		fmt.Fprintf(w, "process still running. go to /ps or /kill")
		return
	}
	cmd := exec.Command("/bin/sh", "-c", "sleep 1000000 && sleep 5")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(w, "failed to start process %d: %v", err)
		return
	}
	fmt.Fprintf(w, "started process pid=%d", cmd.Process.Pid)
	proc = cmd.Process
}

func kill(w http.ResponseWriter, req *http.Request) {
	var id int
	if pid := req.URL.Query().Get("pid"); pid != "" {
		id, _ = strconv.Atoi(pid)
	} else if proc != nil {
		id = -proc.Pid
		proc = nil
	} else {
		fmt.Fprintf(w, "visit /run first to start process")
		return
	}

	if err := syscall.Kill(id, syscall.SIGKILL); err != nil {
		fmt.Fprintf(w, "failed to kill %d: %+v", id, err)
		return
	}
	fmt.Fprintf(w, "killed %d", id)
}

func ps(w http.ResponseWriter, req *http.Request) {
	if runtime.GOOS != "linux" {
		fmt.Fprintf(w, "pstree not available on %q", runtime.GOOS)
		return
	}
	pids, err := pstree.New()
	if err != nil {
		fmt.Fprintf(w, "failed to get pstree: %+v", err)
	}
	var display func(io.Writer, int, int)
	display = func(out io.Writer, pid int, indent int) {
		proc := pids.Procs[pid]
		pp := fmt.Sprintf("pid=%d [ppid=%d,pgrp=%d] (%c) %s", proc.Stat.Pid, proc.Stat.Ppid, proc.Stat.Pgrp, proc.Stat.State, proc.Name)
		prefix := strings.Repeat("  ", indent)
		fmt.Fprintf(out, "%s%s\n", prefix, pp)
		for _, cid := range pids.Procs[pid].Children {
			display(out, cid, indent+1)
		}
	}
	display(w, 1, 0)
}
