package main


import (
    "os"
    "log"
    "flag"
    "launchpad.net/goyaml"
)


func init() {
    idx = flag.String("idx", "", "Platform id for test, will push into amqp")
    action = flag.Bool("", false, "show current version")
}

func main() {
    apps := make(map[string]*App)
    if e := LoadYaml(apps, "foods.yml"); e != nil {
        log.Fatal(e)
    }


    pa := &os.ProcAttr {
        Dir: root,
        Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
    }

    p, e := os.StartProcess("notes", nil, pa)

    log.Println(p, e)
}




type App struct {
    Name string
    Bin string `yaml:"bin"`
    Pid string `yaml:"pid"`
}

func LoadYaml(v interface{}, filename string) error {
    text, e := LoadFile(filename)
    if e != nil {
        return e
    }

    return goyaml.Unmarshal(text, v)
}

func LoadFile(filename string) ([]byte, error) {
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
