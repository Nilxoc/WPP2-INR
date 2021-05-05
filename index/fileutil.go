package index

import (
	"encoding/gob"
	"os"
)

//LoadIndex loads the index dump, you need to cast it to an Index object afterwards..
//Moved it here to avoid circular dependency with util class...
//Not sure if working properly. Need to test later
func loadIndex(path string) (*Index, error) {
	data := Index{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	err2 := dec.Decode(&data)
	if err2 != nil {
		return nil, err2
	}
	return &data, nil
}
