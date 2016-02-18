package abs

import (
	"bytes"
	"github.com/Azure/azure-sdk-for-go/storage"
	"hamustro/dialects"
)

// Azure Queue Storage configuration file.
type Config struct {
	Account    string `json:"account"`
	AccessKey  string `json:"access_key"`
	Container  string `json:"container"`
	BlobPath   string `json:"blob_path"`
	FileFormat string `json:"file_format"`
}

// Checks is it valid or not
func (c *Config) IsValid() bool {
	return c.Account != "" && c.AccessKey != "" && c.Container != "" && c.FileFormat != ""
}

// Create a new StorageClient object based on a configuration file.
func (c *Config) NewClient() (dialects.StorageClient, error) {
	serviceClient, err := storage.NewBasicClient(c.Account, c.AccessKey)
	if err != nil {
		return nil, err
	}
	converterFunction, err := dialects.GetBatchConverterFunction(c.FileFormat)
	if err != nil {
		return nil, err
	}
	return &BlobStorage{
		Account:        c.Account,
		AccessKey:      c.AccessKey,
		BlobPath:       c.BlobPath,
		Container:      c.Container,
		FileFormat:     c.FileFormat,
		BatchConverter: converterFunction,
		Client:         serviceClient.GetBlobService()}, nil
}

// Azure Queue Storage dialect.
type BlobStorage struct {
	Account        string
	AccessKey      string
	Container      string
	BlobPath       string
	FileFormat     string
	BatchConverter dialects.BatchConverter
	Client         storage.BlobStorageClient
}

// It is a buffered storage.
func (c *BlobStorage) IsBufferedStorage() bool {
	return true
}

// There is no normal converter for ABS
func (c *BlobStorage) GetConverter() dialects.Converter {
	return nil
}

// Returns the converter function
func (c *BlobStorage) GetBatchConverter() dialects.BatchConverter {
	return c.BatchConverter
}

// Send a single Event into the Azure Queue Storage.
func (c *BlobStorage) Save(workerID int, msg *bytes.Buffer) error {
	buffer, err := dialects.Compress(msg)
	if err != nil {
		return err
	}
	if err := c.Client.CreateBlockBlobFromReader(c.Container, dialects.GetRandomPath(c.BlobPath, c.FileFormat),
		uint64(buffer.Len()), bytes.NewReader(buffer.Bytes()), nil); err != nil {
		return err
	}
	return nil
}
