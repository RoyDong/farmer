package main


import (
    "os"
    "os/exec"
    "fmt"
)


const (
    DirPerm = 0755
    LogPerm = 0644
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("do nothing")
        return
    }

    var e error
    action := os.Args[1]
    name := os.Args[2]
    if action == "start" {
        cmd := exec.Command(name, os.Args[3:]...)
        if cmd.Stdout, e = createLogfile(name, "stdout"); e != nil {
            println("main log " + e.Error())
        }
        if cmd.Stderr, e = createLogfile(name, "stderr"); e != nil {
            println("error log " + e.Error())
        }
        e = cmd.Start()
    } else if action == "stop" {
        cmd := exec.Command("killall", name)
        e = cmd.Run()
    } else if action == "restart" {
        println("run stop then start")
        return
    } else {
        fmt.Println("do nothing")
        return
    }

    if e != nil {
        println("error: " + e.Error())
    } else {
        fmt.Println("done")
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
