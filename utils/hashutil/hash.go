package hashutil

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func Sha1(str string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}

func Md5File(file string) (string, error) {
	f, err := os.Open(file)
	if nil != err {
		return "", err
	}
	defer f.Close()
	return Md5Reader(f), nil
}

func Md5Reader(reader io.Reader) string {
	m := md5.New()
	io.Copy(m, reader)
	return fmt.Sprintf("%x", m.Sum([]byte("")))
}
func Md5ReaderLen(reader io.Reader) (string, int64) {
	m := md5.New()
	fileSize, _ := io.Copy(m, reader)
	return fmt.Sprintf("%x", m.Sum([]byte(""))), fileSize
}
