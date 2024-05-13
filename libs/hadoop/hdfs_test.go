package hadoop

import (
	"encoding/csv"
	"fmt"
	"github.com/colinmarc/hdfs"
	"io"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	//err := WriteFile()
	//assert.Nil(t, err)
	ReadFile()
}

func TestWriteChunkFile(t *testing.T) {
	// HDFS 连接信息
	client, _ := hdfs.New("192.168.15.58:9000")
	defer client.Close()

	// 读取CSV文件
	csvfile, err := os.Open("C:\\Users\\caijiatao\\PycharmProjects\\data_process\\spark_demo\\article\\item_truncate.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open CSV file: %v\n", err)
		os.Exit(1)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	// 假设CSV文件的第一行是标题，可以根据需要调整
	reader.Read()

	// 分片大小
	chunkSize := 1000 // 假设每个分片包含1000行

	// 分片计数
	chunkCount := 1

	// 循环读取CSV并写入HDFS
	for {
		records := make([][]string, chunkSize)
		i := 0
		for ; i < chunkSize; i++ {
			record, err := reader.Read()
			if err != nil {
				break
			}
			records[i] = record
		}
		if i == 0 {
			break
		}
		records = records[:i]

		// 写入到HDFS
		err := writeChunkToHDFS(client, records, chunkCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write chunk to HDFS: %v\n", err)
			os.Exit(1)
		}
		chunkCount++
	}

	fmt.Println("CSV data uploaded and chunked into multiple files on HDFS.")
}

func TestA(t *testing.T) {

}

func TestUploadCSV(t *testing.T) {
	// 连接到HDFS
	client, _ := hdfs.New("192.168.15.58:9000")
	defer client.Close()

	// 本地CSV文件路径
	localFilePath := "C:\\Users\\caijiatao\\PycharmProjects\\data_process\\spark_demo\\article\\item_truncate.csv"

	// 在HDFS上的目标路径
	hdfsFilePath := "/item_truncate.csv"

	// 打开本地文件
	file, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("Error opening local file:", err)
		return
	}
	defer file.Close()

	// 创建HDFS文件
	hdfsFile, err := client.Create(hdfsFilePath)
	if err != nil {
		fmt.Println("Error creating HDFS file:", err)
		return
	}
	defer hdfsFile.Close()

	// 从本地文件复制到HDFS文件
	_, err = io.Copy(hdfsFile, file)
	if err != nil {
		fmt.Println("Error copying file content:", err)
		return
	}

	fmt.Println("File uploaded successfully to HDFS.")
}
