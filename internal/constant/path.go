package constant

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// ResolveDBPath returns an absolute path for a database file. If filename is
// already absolute it is returned unchanged; otherwise it is joined with
// DataDirPath (which must already be initialised via InitPaths).
func ResolveDBPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(DataDirPath, filename)
}

var (
	ConfigDirPath string
	DataDirPath   string
)

// InitPaths resolves the config/data directories from the environment (or XDG
// defaults) and ensures they exist. It must be called once at startup before
// any code reads ConfigDirPath/DataDirPath.
func InitPaths() error {
	if naliHome := os.Getenv("NALI_HOME"); len(naliHome) != 0 {
		ConfigDirPath = naliHome
		DataDirPath = naliHome
	} else {
		ConfigDirPath = os.Getenv("NALI_CONFIG_HOME")
		if len(ConfigDirPath) == 0 {
			ConfigDirPath = filepath.Join(xdg.ConfigHome, "nali")
		}

		DataDirPath = os.Getenv("NALI_DB_HOME")
		if len(DataDirPath) == 0 {
			DataDirPath = filepath.Join(xdg.DataHome, "nali")
		}
	}

	if err := prepareDir(ConfigDirPath); err != nil {
		return err
	}
	return prepareDir(DataDirPath)
}

func prepareDir(dir string) error {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("can not create dir %q: %w", dir, err)
		}
		return nil
	} else if err != nil {
		return err
	} else if !stat.IsDir() {
		return fmt.Errorf("path already exists, but is not a dir: %s", dir)
	}
	return nil
}
