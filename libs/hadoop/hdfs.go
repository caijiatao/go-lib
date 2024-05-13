package hadoop

import (
	"encoding/csv"
	"fmt"
	"github.com/colinmarc/hdfs"
)

func writeChunkToHDFS(client *hdfs.Client, records [][]string, chunkCount int) error {
	// 创建一个新文件
	filename := fmt.Sprintf("/item_chunk_%d.csv", chunkCount)
	file, err := client.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将记录写入到文件
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadFile() {
	client, _ := hdfs.New("192.168.15.58:9000")
	defer client.Close()

	file, _ := client.Open("/item_chunk_1.csv")

	buf := make([]byte, 1024)
	for {
		read, err := file.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(read)
		fmt.Println(string(buf))
		if read == 0 {
			break
		}
	}
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
