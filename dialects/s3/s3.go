package s3

import (
	".."
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
)

// Amazon SNS configuration file.
type Config struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Bucket          string `json:"bucket"`
	BlobPath        string `json:"blob_path"`
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	FileFormat      string `json:"file_format"`
}

// Checks is it valid or not
func (c *Config) IsValid() bool {
	return c.AccessKeyID != "" && c.SecretAccessKey != "" && c.Bucket != "" && c.BlobPath != "" && c.Region != "" && c.Endpoint != ""
}

// Create a new StorageClient object based on a configuration file.
func (c *Config) NewClient() (dialects.StorageClient, error) {
	creds := credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, "")
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}
	converterFunction, err := dialects.GetConverterFunction(c.FileFormat)
	if err != nil {
		return nil, err
	}
	config := &aws.Config{
		Region:           &c.Region,
		Credentials:      creds,
		Endpoint:         &c.Endpoint,
		S3ForcePathStyle: aws.Bool(true)}
	return &S3Storage{
		AccessKeyID:     c.AccessKeyID,
		SecretAccessKey: c.SecretAccessKey,
		Bucket:          c.Bucket,
		BlobPath:        c.BlobPath,
		Region:          c.Region,
		Converter:       converterFunction,
		Client:          s3.New(session.New(), config)}, nil
}

// SNS Storage's dialect.
type S3Storage struct {
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	BlobPath        string
	Region          string
	EndPoint        string
	Converter       dialects.Converter
	Client          *s3.S3
}

// It is a buffered storage.
func (c *S3Storage) IsBufferedStorage() bool {
	return true
}

// Returns the converter function
func (c *S3Storage) GetConverter() dialects.Converter {
	return c.Converter
}

// Publish a single Event to SNS topic.
func (c *S3Storage) Save(msg *string) error {
	buffer, err := dialects.Compress(msg)
	if err != nil {
		return err
	}
	fileSize := buffer.Len()
	fileBytes := buffer.Bytes()
	params := &s3.PutObjectInput{
		Bucket:        &c.Bucket,
		Key:           aws.String(dialects.GetRandomPath(c.BlobPath)),
		Body:          bytes.NewReader(fileBytes),
		ContentLength: aws.Int64(int64(fileSize)),
		ContentType:   aws.String(http.DetectContentType(fileBytes)),
		Metadata:      map[string]*string{"Key": aws.String("MetadataValue")}}
	_, err = c.Client.PutObject(params)
	if err != nil {
		return err
	}
	return nil
}
