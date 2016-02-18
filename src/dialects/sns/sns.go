package sns

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sub-ninja/hamustro/src/dialects"
)

// Amazon SNS configuration file.
type Config struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	TopicArn        string `json:"topic_arn"`
	Region          string `json:"region"`
}

// Checks is it valid or not
func (c *Config) IsValid() bool {
	return c.AccessKeyID != "" && c.SecretAccessKey != "" && c.TopicArn != "" && c.Region != ""
}

// Create a new StorageClient object based on a configuration file.
func (c *Config) NewClient() (dialects.StorageClient, error) {
	creds := credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, "")
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}
	return &SNSStorage{
		AccessKeyID:     c.AccessKeyID,
		SecretAccessKey: c.SecretAccessKey,
		TopicArn:        c.TopicArn,
		Client:          sns.New(session.New(), &aws.Config{Region: &c.Region, Credentials: creds})}, nil
}

// SNS Storage's dialect.
type SNSStorage struct {
	AccessKeyID     string
	SecretAccessKey string
	TopicArn        string
	Region          string
	Client          *sns.SNS
}

// It is a buffered storage.
func (c *SNSStorage) IsBufferedStorage() bool {
	return false
}

// Returns the converter function
func (c *SNSStorage) GetConverter() dialects.Converter {
	return dialects.ConvertJSON
}

// Returns the batch converter function
func (c *SNSStorage) GetBatchConverter() dialects.BatchConverter {
	return nil
}

// Publish a single Event to SNS topic.
func (c *SNSStorage) Save(msg *bytes.Buffer) error {
	params := &sns.PublishInput{
		Message:  aws.String(msg.String()),
		TopicArn: &c.TopicArn}
	_, err := c.Client.Publish(params)
	if err != nil {
		return err
	}
	return nil
}
