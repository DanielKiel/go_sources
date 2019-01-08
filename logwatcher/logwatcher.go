package main

import (
  "fmt"
  "path/filepath"
  "os"
  "io"
  "bytes"
  "github.com/hpcloud/tail"
)

type Logs map[string]Log

type Log struct {
  lines int
}

func main() {
  path := "/home/daniel/Public/"
  //initial files getting
  files := getFiles(path)

  for path,file := range files {
    fmt.Println(path, file.lines)
  }

  //watch them
  for {
    for path,file := range files {

      reader, err := os.Open(path)

      check(err)

      lines, err := lineCounter(reader)

      check(err)

      if (lines > file.lines || lines < file.lines) {

        fmt.Println("___________________________")
        fmt.Println("changes detected")
        fmt.Println(path, file.lines)
        fmt.Println("now",lines)
        fmt.Println("___________________________")

        log := Log{lines: lines}
        files[path] = log
      }


    }
  }
}

func check(e error) {
  if e != nil {
      panic(e)
  }
}

func show(path string) {
  t, err := tail.TailFile(path, tail.Config{Follow: true,ReOpen:true})

  check(err)

  for line := range t.Lines {
      fmt.Println(line.Text)
  }
  fmt.Println(t.Lines,path)
}

func lineCounter(r io.Reader) (int, error) {
  buf := make([]byte, 32*1024)
  count := 0
  lineSep := []byte{'\n'}

  for {
      c, err := r.Read(buf)
      count += bytes.Count(buf[:c], lineSep)

      switch {
      case err == io.EOF:
          return count, nil

      case err != nil:
          return count, err
      }
  }
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
        reader, err := os.Open(path)

        check(err)

        lines, err := lineCounter(reader)

        check(err)

        log := Log{lines: lines}
        logs[ident] = log
      }
      return nil
  })
  if err != nil {
      panic(err)
  }

  return logs

}
