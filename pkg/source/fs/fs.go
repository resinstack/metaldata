package fs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/resinstack/metaldata/pkg/http"
)

// FS implements the filesystem strategy for retrieving machine
// metadata from files on disk.
type FS struct {
	baseDir string
}

// New sets up the filesystem source and hands it back.
func New(baseDir string) http.InfoSource {
	return &FS{baseDir: baseDir}
}

// GetMachineInfo returns a specific key for a given machine.
func (fs *FS) GetMachineInfo(mach, key string) (string, error) {
	mach = strings.ReplaceAll(mach, ":", "-")

	vals := make(map[string]string)

	bytes, err := ioutil.ReadFile(filepath.Join(fs.baseDir, mach+".json"))
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(bytes, &vals); err != nil {
		return "", err
	}

	v, ok := vals[key]
	if !ok {
		return "", errors.New("machine known, no value available")
	}
	return v, nil
}
