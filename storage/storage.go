package storage

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/SecurityBrewery/catalyst/generated/pointer"
)

type Storage struct {
	session *session.Session
}

type Config struct {
	Host     string
	User     string
	Password string
}

func New(config *Config) (*Storage, error) {
	s, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.User, config.Password, ""),
		Endpoint:         aws.String(config.Host),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})

	return &Storage{s}, err
}

func (s *Storage) S3() *s3.S3 {
	return s3.New(s.session)
}

func (s *Storage) Downloader() *s3manager.Downloader {
	d := s3manager.NewDownloader(s.session)
	d.Concurrency = 1

	return d
}

func (s *Storage) Uploader() *s3manager.Uploader {
	d := s3manager.NewUploader(s.session)
	d.Concurrency = 1

	return d
}

func (s *Storage) DeleteBucket(name string) error {
	_, err := s.S3().DeleteBucket(&s3.DeleteBucketInput{Bucket: pointer.String("catalyst-" + name)})

	return err
}

func CreateBucket(client *s3.S3, ticketID string) error {
	_, err := client.CreateBucket(&s3.CreateBucketInput{Bucket: pointer.String("catalyst-" + ticketID)})
	if err == nil {
		err = client.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: pointer.String("catalyst-" + ticketID)})
		if err != nil {
			return err
		}
	} else {
		var awsError awserr.Error
		if errors.As(err, &awsError) && (awsError.Code() == s3.ErrCodeBucketAlreadyExists || awsError.Code() == s3.ErrCodeBucketAlreadyOwnedByYou) {
			return nil
		}

		return err
	}

	return err
}
