package main

import (
  "fmt"
  "os"
  "encoding/csv"
  "encoding/json"
  "bufio"
  "io"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
  "encoding/base64"
  "net/http"
  "time"
  "io/ioutil"
  "runtime"
  "sync/atomic"
)

const CHANELBUFFER = 14
const WORKER = 12
const WRITTEN = 0

type Message struct {
  Name string `json:"name"`
  Foo string `json:"foo"`
}

func main() {
  fmt.Println(runtime.NumCPU())
  runtime.GOMAXPROCS(40)
  now := time.Now()


  readedChan :=read("./Downloads/my_test_encrypted.csv")

  appendedChan := addData(readedChan)

  write("./Creates/my_test_decrypted.csv", appendedChan, now)

  select{}
}

func addData(readedChan <- chan[] string) <- chan []string {
  appendedChan := make(chan []string, CHANELBUFFER)
  var counter uint64

  for worker := 0; worker < WORKER; worker++ {
    go appendLines(readedChan, appendedChan, &counter)
  }

  return appendedChan
}

func appendLines(readedChan <- chan[] string, appendedChan chan <- [] string, counter *uint64) {

  for line := range readedChan {
    atomic.AddUint64(counter, 1)
    myClient := &http.Client{Timeout:90 * time.Second}
    resp, err := myClient.Get("http://192.168.1.4:80")

    if (err != nil) {
      continue
    }
    defer resp.Body.Close()

    //var message Message
    body, err := ioutil.ReadAll(resp.Body)
    check(err)

    var message Message
    js := json.Unmarshal(body,&message)
    check(js)

    line = append(line, message.Name)
    line = append(line, message.Foo)

    i := atomic.LoadUint64(counter)
    fmt.Println("appended", i)

    appendedChan <- line
  }
}

func write(path string, appendedChan <- chan[] string, now time.Time) {
  file, err := os.Create(path)

  check(err)

  //defer file.Close()

  writer := csv.NewWriter(file)
  defer writer.Flush()
  var counter uint64

  for worker := 0; worker < WORKER; worker++ {
    go writeLines(writer, appendedChan,now,&counter)
  }

}

func writeLines(writer *csv.Writer, writeChan <- chan [] string, now time.Time, counter *uint64) {

  for line := range writeChan {
    atomic.AddUint64(counter, 1)
    //encrypted, error := encrypt("new val", "0123456789012345")
    //check(error)
    //line := append(line, encrypted)

    err := writer.Write(line)
    check(err)

    t1 := time.Now()
    i := atomic.LoadUint64(counter)
	  fmt.Printf("________________we have %n items at  %v to run.\n",i, t1.Sub(now))

    fmt.Println(line)
    i++
  }
}

func read(path string) <- chan []string {
  out := make(chan []string, CHANELBUFFER)
  csvFile, err := os.Open(path)
  check(err)

  reader := csv.NewReader(bufio.NewReader(csvFile))
  var counter uint64
  for worker := 0; worker < WORKER; worker++ {
    go readLines(reader, out, &counter)
  }


  return out
}

func readLines(reader *csv.Reader, readedChan chan <-[] string, counter *uint64) {
  for {
    line, err := reader.Read()
    atomic.AddUint64(counter, 1)

    if (err == io.EOF) {
      close(readedChan)
      break
    }

    check(err)

    decrypted, error := decrypt(line[5], "0123456789012345")
    check(error)

    line[5] = decrypted

    i := atomic.LoadUint64(counter)
    fmt.Println("___readed", i)

    readedChan <- line
  }
}

func encrypt(message string, passphrase string) (encmessage string, err error) {
	key := []byte(passphrase)

  plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmessage = base64.URLEncoding.EncodeToString(cipherText)
  return
}

func decrypt(securemess string, passphrase string) (decodedmess string, err error) {
  key := []byte(passphrase)
  cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		//err = errors.New("Ciphertext block size is too short!")
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func check(err error) {
  if (err != nil) {
    panic(err)
  }
}
