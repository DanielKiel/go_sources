package main

import (
  "fmt"
  "os"
  fm "logwatcher/filemanager"
  "github.com/subosito/gotenv"
  "golang.org/x/crypto/ssh"
)



func main() {
  //fmt.Println(os.Getenv("SSH_USERNAME"))

  path := "/home/daniel/Public/"
  stream := fm.Stream(path)
  //initial files getting
  //files := <- stream//fm.GetFiles(path)
  //fmt.Println(file.Path, file.Lines)
  for {
      file := <- stream

      reader, err := os.Open(file.Path)

      check(err)

      lines, err := fm.LineCounter(reader)

      check(err)

      if (lines > file.Lines || lines < file.Lines) {

        fmt.Println("___________________________")
        fmt.Println("changes detected")
        fmt.Println(file.Path, file.Lines)
        fmt.Println("now",lines)
        fmt.Println("___________________________")
      }
  }
}

func check(e error) {
  if e != nil {
      panic(e)
  }
}

func connectToHost(host string) (*ssh.Client, *ssh.Session, error) {
  gotenv.Load()
	user := os.Getenv("SSH_USERNAME")
  pass := os.Getenv("SSH_PASSWORD")

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}
