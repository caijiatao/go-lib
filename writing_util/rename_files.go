package writing_util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func removeSpaceFilesInFolder(folderPath string, suffix string) error {
	// 获取文件夹下所有文件
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// 检查文件是否以 .png 结尾并包含空格
		if strings.HasSuffix(file.Name(), suffix) && strings.Contains(file.Name(), " ") {
			// 构建新文件名，去掉空格
			newFileName := strings.ReplaceAll(file.Name(), " ", "")

			// 构建完整的文件路径
			oldFilePath := filepath.Join(folderPath, file.Name())
			newFilePath := filepath.Join(folderPath, newFileName)

			// 重命名文件
			err := os.Rename(oldFilePath, newFilePath)
			if err != nil {
				return err
			}
			fmt.Printf("Renamed: %s to %s\n", oldFilePath, newFilePath)
		}
	}

	return nil
}
