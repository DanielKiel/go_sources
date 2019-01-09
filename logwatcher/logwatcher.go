package main

import (
  "fmt"
  "path/filepath"
  "os"
  "io"
  "bytes"
  "bufio"
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

        show(path)
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
  file, err := os.Open(path)

  check(err)

  defer file.Close()

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    fmt.Println(scanner.Text())
  }
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
