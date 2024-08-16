package main

import (
	"bytes"
	"compress/gzip"
	"github.com/golang/snappy"
	"io/ioutil"
	"log"
)

// 压缩数据
func snappyCompress(data []byte) []byte {
	return snappy.Encode(nil, data)
}

// 解压缩数据
func snappyDecompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}

// 压缩数据
func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	defer writer.Close()

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	writer.Close()
	return buf.Bytes(), nil
}

// 解压缩数据
func decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func main() {
	originalData := []byte("this is the original data")
	log.Printf("Original size: %d\n", len(originalData))

	compressedData, err := compress(originalData)
	if err != nil {
		log.Fatalf("Error compressing data: %v", err)
	}
	// compressedData的大小
	log.Printf("Compressed size: %d\n", len(compressedData))

	decompressedData, err := decompress(compressedData)
	if err != nil {
		log.Fatalf("Error decompressing data: %v", err)
	}

	if string(decompressedData) != string(originalData) {
		log.Fatalf("Decompressed data is different from original data")
	}
}
