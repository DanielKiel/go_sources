package main

import (
  "fmt"
  "path/filepath"
  "os"
  "github.com/hpcloud/tail"
)

type Logs map[string]Log

type Log struct {
  path string
}

func main() {
  path := "/home/daniel/Public/LogTests/one"
  files := getFiles(path)

  for _,file := range files {
    fmt.Println(file.path)
    notify(file.path)
  }
}

func notify(path string) {
  t, err := tail.TailFile(path, tail.Config{Follow: true,ReOpen:true})

  if (err != nil) {
    fmt.Println("error", err)
  }

  for line := range t.Lines {
      fmt.Println(line.Text)
  }
  fmt.Println(t.Lines,path)
}

func getFiles(root string) Logs {
  logs := make(Logs)

  err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
      //fmt.Println(filepath.Dir(path))
      ident, miss := filepath.Abs(path)

      if (miss != nil) {
        fmt.Println(miss)
      }

      //check if the key already exists
      _,exists := logs[ident]

      if (exists) {
        fmt.Println("we must handle the error here")
      }

      if (filepath.Ext(path) == ".log") {
        log := Log{path: path}
        logs[ident] = log
      }
      return nil
  })
  if err != nil {
      panic(err)
  }

  return logs

}
