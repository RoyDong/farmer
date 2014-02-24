package main


import (
    "os"
    "os/exec"
    "fmt"
)


const (
    DirPerm = 0755
    LogPerm = 0655
)

func main() {
    if len(os.Args) < 3 {
        println("do nothing")
        return
    }

    var e error
    action := os.Args[1]
    name := os.Args[2]
    if action == "start" {
        cmd := exec.Command(name, os.Args[1:]...)
        if cmd.Stdout, e = createLogfile(name, "main"); e != nil {
            println(e)
        }
        if cmd.Stderr, e = createLogfile(name, "error"); e != nil {
            println(e)
        }
        e = cmd.Start()
    } else if action == "stop" {
        cmd := exec.Command("killall", name)
        e = cmd.Run()
    } else if action == "restart" {
        println("you fucking lazy, run stop then start")
        return
    } else {
        println("do nothing")
        return
    }

    if e != nil {
        println("error: " + e.Error())
    } else {
        println("done")
    }
}

func createLogfile(n, t string) (*os.File, error) {
    var f *os.File
    var e error
    e = os.MkdirAll("/var/log/" + n, DirPerm)
    if e != nil {
        return nil, e
    }
    f, e = os.OpenFile(fmt.Sprintf("/var/log/%s/%s.log", n, t), os.O_CREATE|os.O_WRONLY|os.O_APPEND, LogPerm)
    if e != nil {
        return nil, e
    }
    return f, nil
}
