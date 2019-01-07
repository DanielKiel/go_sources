package main

import (
  "fmt"
  "io/ioutil"
  "log"
)

type Directory struct {
  name string
  size int64
}

type Directories map[string]Directory

func main() {
  path := "/usr/"
  directories := getDirectories(path)

  for {
   //and now campare them!!
   comparism := getDirectories(path)

   for key,value := range comparism {
     fmt.Println(value.name)
     fmt.Println(value.size)

     compared := value.size > directories[key].size || value.size < directories[key].size

     fmt.Println("compare is", compared)

     fmt.Println("___________________________________________")
   }
  }
}

func getDirectories(path string)Directories {
  files, err := ioutil.ReadDir(path)
  directories := make(Directories)
  if (err != nil) {
    log.Fatal(err)
  }

  for _, f := range files {
    directoryObj := Directory {name: f.Name(),size: f.Size()}
    directories[f.Name()] = directoryObj
  }

  return directories
}
