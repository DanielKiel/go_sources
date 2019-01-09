package main

import (
  "fmt"
  "os"
  fm "logwatcher/filemanager"
  "github.com/subosito/gotenv"
)



func main() {
  gotenv.Load()

  //fmt.Println(os.Getenv("SSH_USERNAME"))

  path := "/home/daniel/Public/"
  //initial files getting
  files := fm.GetFiles(path)

  for path,file := range files {
    fmt.Println(path, file.Lines)
  }

  //watch them
  for {
    for path,file := range files {

      reader, err := os.Open(path)

      check(err)

      lines, err := fm.LineCounter(reader)

      check(err)

      if (lines > file.Lines || lines < file.Lines) {

        fmt.Println("___________________________")
        fmt.Println("changes detected")
        fmt.Println(path, file.Lines)
        fmt.Println("now",lines)
        fmt.Println("___________________________")

        log := fm.Log{Lines: lines}
        files[path] = log

        //show(path)
      }
    }
  }
}

func check(e error) {
  if e != nil {
      panic(e)
  }
}
