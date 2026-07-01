package common

import (
	"fmt"
	"os"
	"path/filepath"
)

// SaveFile atomically writes data to path: it writes to a temp file in the same
// directory and renames it into place, so a crash or interrupt mid-write can
// never leave a truncated database behind.
func SaveFile(path string, data []byte) (err error) {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".nali-tmp-*")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }() // no-op after a successful rename

	if _, err = tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("写入临时文件失败: %w", err)
	}
	if err = tmp.Close(); err != nil {
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}

	if err = os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}
	return nil
}
