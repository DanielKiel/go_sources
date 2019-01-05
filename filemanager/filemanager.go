package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "strconv"
)

func main() {
   directories := getDirectories()

   for _, d := range directories {
     fmt.Println(d)
   }
}

func getDirectories() []string {
  files, err := ioutil.ReadDir("/usr/")
  directories := []string{}

  if (err != nil) {
    log.Fatal(err)
  }

  for _, f := range files {
    //fmt.Println(f.Name())
    size := strconv.FormatInt(f.Size(), 10)
    name := f.Name()
    entry := name + ": " + size
    //entry := strings.Join(f.Name(),": ", string(f.Size()))
    //fmt.Println(f.Size())
    directories = append(directories, entry)
  }

  return directories
}
