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
)

type Message struct {
  Name string `json:"name"`
  Foo string `json:"foo"`
}

func main() {
  //reader := read("./Downloads/my_test_2.csv")
  write("./Creates/my_test_decrypted.csv", read("./Downloads/my_test_encrypted.csv"))
}

func write(path string, reader <- chan[] string) {
  file, err := os.Create(path)

  check(err)

  defer file.Close()

  writer := csv.NewWriter(file)
  defer writer.Flush()

  for line := range reader {
    //encrypted, error := encrypt("new val", "0123456789012345")
    //check(error)
    //line := append(line, encrypted)

    err := writer.Write(line)
    check(err)
  }
}

func read(path string) <- chan []string {

  out := make(chan []string)
  csvFile, err := os.Open(path)
  check(err)

  go func() {
    reader := csv.NewReader(bufio.NewReader(csvFile))

    for {
      line, err := reader.Read()

      if (err == io.EOF) {
        close(out)
        break
      }

      check(err)

      decrypted, error := decrypt(line[5], "0123456789012345")
      check(error)

      line[5] = decrypted

      myClient := &http.Client{Timeout:60 * time.Second}
      resp, err := myClient.Get("http://192.168.1.4:80")

      check(err)
      defer resp.Body.Close()

      //var message Message
      body, err := ioutil.ReadAll(resp.Body)
      check(err)

      var message Message
      js := json.Unmarshal(body,&message)
      check(js)
fmt.Println(message.Name,message.Foo)
      //line[6] = message.Name
      //line[7] = message.Foo

      out <- line
    }
  }()

  return out
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
