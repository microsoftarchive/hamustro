package file

import (
	"bytes"
	"github.com/wunderlist/hamustro/src/dialects"
	"os"
)

// Local file configuration
type Config struct {
	FilePath   string `json:"file_path"`
	FileFormat string `json:"file_format"`
	Compress   bool   `json:"compress"`
}

// Checks is it valid or not
func (c *Config) IsValid() bool {
	return c.FilePath != "" && c.FileFormat != ""
}

// Create a new StorageClient object based on a configuration file.
func (c *Config) NewClient() (dialects.StorageClient, error) {
	converterFunction, err := dialects.GetBatchConverterFunction(c.FileFormat)
	if err != nil {
		return nil, err
	}

	return &FileStorage{
		FilePath:       c.FilePath,
		FileFormat:     c.FileFormat,
		Compress:       c.Compress,
		BatchConverter: converterFunction}, nil
}

// FileStorage's dialect.
type FileStorage struct {
	FilePath       string
	FileFormat     string
	Compress       bool
	BatchConverter dialects.BatchConverter
}

// It is a buffered storage.
func (c *FileStorage) IsBufferedStorage() bool {
	return true
}

// Returns the converter function
func (c *FileStorage) GetConverter() dialects.Converter {
	return nil
}

// Returns the batch converter function
func (c *FileStorage) GetBatchConverter() dialects.BatchConverter {
	return c.BatchConverter
}

func (c *FileStorage) GetBuffer(msg *bytes.Buffer) (*bytes.Buffer, error) {
	if c.Compress {
		buffer, err := dialects.Compress(msg)
		if err != nil {
			return nil, err
		}
		return buffer, nil
	}
	return msg, nil
}

// Write a single local file with multiple records
func (c *FileStorage) Save(msg *bytes.Buffer) error {
	buffer, err := c.GetBuffer(msg)
	if err != nil {
		return err
	}

	basepath := dialects.ResolvePath(c.FilePath)
	if err := os.MkdirAll(basepath, os.ModePerm); err != nil {
		return err
	}

	path := dialects.GetRandomPath(basepath, c.FileFormat, c.Compress)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data := buffer.Bytes()
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}
