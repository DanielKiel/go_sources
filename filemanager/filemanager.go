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
   directories := getDirectories()

   for {
     //and now campare them!!
     comparism := getDirectories()

     for key,value := range comparism {
       fmt.Println(value.name)
       fmt.Println(value.size)

       compared := value.size > directories[key].size || value.size < directories[key].size

       fmt.Println("compare is", compared)

       fmt.Println("___________________________________________")
     }
   }


}

func getDirectories()Directories {
  files, err := ioutil.ReadDir("/usr/")
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
