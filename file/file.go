package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func AbsPath(p string) (string, error) {
	if filepath.IsAbs(p) {
		return p, nil
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, p), nil
}

func readAsByte(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func ReadAsString(path string) (string, error) {
	bytes, err := readAsByte(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

//GetFileReader returns a os.File Object implementing the io.Reader interface. It is required to call file.Close() when finished using this file object!
func GetFileReader(path string) (*os.File, error) {
	return os.Open(path)
}

func SaveIndex(data interface{}, path string) error {
	/*file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	return enc.Encode(data)*/
	return fmt.Errorf("not supported with vektor")
}

func ListFiles(pathD string) ([]string, error) {
	entr, err := os.ReadDir(pathD)
	res := make([]string, 0, len(entr))
	if err != nil {
		return nil, err
	}
	for _, f := range entr {
		if !f.IsDir() {
			res = append(res, filepath.Join(pathD, f.Name()))
		}
	}
	return res, nil
}
