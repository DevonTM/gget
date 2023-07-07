package progress

import (
	"compress/zlib"
	"encoding/json"
	"os"
)

type Data struct {
	Chunk     []bool `json:"chunk"`
	ChunkSize int64  `json:"chunk_size"`
	FileSize  int64  `json:"-"`
}

// load progress from file or create new progress file if not exist
func (pd *Data) LoadProgress(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = pd.CreateProgressFile(path)
		if err != nil {
			return err
		}
	} else {
		err = pd.ReadProgressFile(path)
		if err != nil {
			return err
		}
	}

	return nil
}

// calculate chunk slice and create progress file
func (pd *Data) CreateProgressFile(path string) error {
	chunkCount := int(pd.FileSize / pd.ChunkSize)
	if pd.FileSize%pd.ChunkSize != 0 {
		chunkCount++
	}

	pd.Chunk = make([]bool, chunkCount)

	err := pd.SaveProgressFile(path)
	if err != nil {
		return err
	}

	return nil
}

// save progress data to file with zlib compression
func (pd *Data) SaveProgressFile(path string) error {
	jsonData, err := json.Marshal(pd)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	zlibWriter, err := zlib.NewWriterLevel(file, zlib.BestSpeed)
	if err != nil {
		return err
	}
	defer zlibWriter.Close()

	_, err = zlibWriter.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

// read compressed progress file and decode to progress data
func (pd *Data) ReadProgressFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader, err := zlib.NewReader(file)
	if err != nil {
		return err
	}
	defer reader.Close()

	err = json.NewDecoder(reader).Decode(pd)
	if err != nil {
		return err
	}

	return nil
}

// check if all chunk already finished
func (pd *Data) IsFinished() bool {
	for _, chunk := range pd.Chunk {
		if !chunk {
			return false
		}
	}

	return true
}
