package common

import (
	"fmt"
	"os"
)

func SaveFile(path string, data []byte) (err error) {
	// Remove file if exist
	_, err = os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			return fmt.Errorf("旧文件删除失败: %w", err)
		}
	}

	// save file
	return os.WriteFile(path, data, 0644)
}
