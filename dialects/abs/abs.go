package abs

import (
	".."
	"bytes"
	"github.com/Azure/azure-sdk-for-go/storage"
)

// Azure Queue Storage configuration file.
type Config struct {
	Account   string `json:"account"`
	AccessKey string `json:"access_key"`
	Container string `json:"container"`
	BlobPath  string `json:"blob_path"`
}

// Checks is it valid or not
func (c *Config) IsValid() bool {
	return c.Account != "" && c.AccessKey != "" && c.BlobPath != "" && c.Container != ""
}

// Create a new StorageClient object based on a configuration file.
func (c *Config) NewClient() (dialects.StorageClient, error) {
	serviceClient, err := storage.NewBasicClient(c.Account, c.AccessKey)
	if err != nil {
		return nil, err
	}
	return &BlobStorage{
		Account:   c.Account,
		AccessKey: c.AccessKey,
		BlobPath:  c.BlobPath,
		Container: c.Container,
		Client:    serviceClient.GetBlobService()}, nil
}

// Azure Queue Storage dialect.
type BlobStorage struct {
	Account   string
	AccessKey string
	Container string
	BlobPath  string
	Client    storage.BlobStorageClient
}

// It is a buffered storage.
func (c *BlobStorage) IsBufferedStorage() bool {
	return true
}

// Send a single Event into the Azure Queue Storage.
func (c *BlobStorage) Save(msg *string) error {
	buffer, err := dialects.Compress(msg)
	if err != nil {
		return err
	}
	if err := c.Client.CreateBlockBlobFromReader(c.Container, dialects.GetRandomPath(c.BlobPath),
		uint64(buffer.Len()), bytes.NewReader(buffer.Bytes()), nil); err != nil {
		return err
	}
	return nil
}
