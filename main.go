package main


import (
    "os"
    "log"
    "fmt"
    "flag"
    "strconv"
    "syscall"
    "launchpad.net/goyaml"
)


func init() {
    flag.Parse()
}

func main() {
    args := flag.Args()
    if len(args) < 2 {
        log.Fatal("not enough params")
    }

    var app *App
    var name string
    for name, app = range apps() {
        if name == args[1] { break }
    }

    if app == nil {
        log.Fatal("no app named", args[1])
    }

    var pid int
    if b, e := loadFile(app.Pid); e == nil {
        if i, e := strconv.ParseInt(string(b), 10, 0); e == nil {
            pid = int(i)

            if _,e := os.FindProcess(pid); e != nil {
                pid = 0
            }
        }
    }

    var action = args[0]
    if (action == "rt" || action == "restart") {
        if pid == 0 {
            log.Fatal("process is not running")
        }

        if e := syscall.Kill(pid, syscall.SIGKILL); e == nil {
            fmt.Println("old process killed: ", pid)
            run(app)
        } else {
            log.Fatal("can not kill old process", e)
        }

    } else if action == "start" {
        if pid > 0 {
            log.Fatal("app is allready running")
        }

        run(app)
    } else if action == "stop" {
        if pid == 0 {
            log.Fatal("app is allready stopped")
        }

        file,_ := os.OpenFile(app.Pid, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0666)
        defer file.Close()

        if e := syscall.Kill(pid, syscall.SIGKILL); e == nil {
            fmt.Println("old process killed: ", pid)
        } else {
            log.Fatal(e)
        }

    } else {
        log.Fatal("la la la")
    }
}


func run(app *App) {
    p, e := os.StartProcess(app.Cmd, nil, &os.ProcAttr {Dir: app.Dir})
    if e != nil {
        log.Fatal("start app error ", e)
    }

    file, e := os.OpenFile(app.Pid,
        os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0666)
    if e != nil {
        log.Fatal(app.Pid, "file can not open")
    }
    defer file.Close()

    if _,e := file.Write([]byte(fmt.Sprintf("%d", p.Pid))); e != nil {
        log.Fatal(app.Pid, "file can not save")
    }

    fmt.Println("new process running on pid: ", p.Pid)
}


type App struct {
    Name string
    Cmd string `yaml:"cmd"`
    Dir string `yaml:"dir"`
    Pid string `yaml:"pid"`
}

func apps() map[string]*App {
    apps := make(map[string]*App)
    if e := loadYaml(apps, "config.yml"); e != nil {
        log.Fatal(e)
    }

    return apps
}

func loadYaml(v interface{}, filename string) error {
    text, e := loadFile(filename)
    if e != nil {
        return e
    }

    return goyaml.Unmarshal(text, v)
}

func loadFile(filename string) ([]byte, error) {
    file, e := os.Open(filename)
    if e != nil {
        return nil, e
    }
    defer file.Close()

    fileInfo, e := file.Stat()
    if e != nil {
        return nil, e
    }

    text := make([]byte, fileInfo.Size())
    file.Read(text)
    return text, nil
}
