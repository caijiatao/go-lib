package hadoop

import (
	"fmt"
	"github.com/colinmarc/hdfs"
)

func ReadFile() {
	client, _ := hdfs.New("192.168.15.58:9000")
	defer client.Close()

	file, _ := client.Open("/test.txt")

	buf := make([]byte, 1024)
	read, err := file.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(read)
	fmt.Println(string(buf))
}

func WriteFile() error {
	client, _ := hdfs.New("192.168.15.58:9000")
	defer client.Close()
	err := client.Remove("/test.txt")
	if err != nil {
		return err
	}
	file, err := client.CreateFile("/test.txt", 1, 10485760, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte("hello, world"))
	if err != nil {
		return err
	}

	//err = file.Flush()
	//if err != nil {
	//	return err
	//}

	return nil
}