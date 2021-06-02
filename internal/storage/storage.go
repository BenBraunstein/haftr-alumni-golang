package storage

import (
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// storage environment variables
const (
	defaultRegion = "us-east-1"
)

// Config is a representation of pubsub configurations
type Config struct {
	Region string
}

// UploadOption is a representation of an upload uption
type UploadOption func(r *OptionalUploadRequest)

// OptionalUploadRequest is representation of optional parameters for an upload request
type OptionalUploadRequest struct {
	MetaData    map[string]*string
	ContentType string
}

// DefaultConfig returns the default config
func DefaultConfig() Config {
	return Config{Region: defaultRegion}
}

// UploadImageFunc is a function that takes in a reader of an image file and a storage key and uploads it to S3
type UploadImageFunc func(r io.Reader, contentType, key, fileName string) error

type GetImageURLFunc func(key string) (string, error)

// UploadFunc func for uploading data to s3
type UploadFunc func(reader io.Reader, bucket string, key string, opts ...UploadOption) error

// PresignFunc func for presigning s3 object
type PresignFunc func(bucket string, key string) (string, error)

// UploadImage uploads an image file to S3
func UploadImage(upload UploadFunc, bucket string) UploadImageFunc {
	return func(r io.Reader, contentType, key, fileName string) error {
		metaDataOpt := func(r *OptionalUploadRequest) {
			r.MetaData = map[string]*string{
				"clientFilename": &fileName,
			}
		}
		contentTypeOpt := func(r *OptionalUploadRequest) {
			r.ContentType = contentType
		}
		return upload(r, bucket, key, metaDataOpt, contentTypeOpt)
	}
}

func GetImageURL(presignURL PresignFunc, bucket string) GetImageURLFunc {
	return func(key string) (string, error) {
		return presignURL(bucket, key)
	}
}

// UploadToS3 default implementation of s3 uploader
func UploadToS3(c Config) UploadFunc {
	return func(reader io.Reader, bucket string, key string, opts ...UploadOption) error {
		optRequestInput := OptionalUploadRequest{
			MetaData:    map[string]*string{},
			ContentType: "application/octet-stream",
		}
		for _, opt := range opts {
			opt(&optRequestInput)
		}
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(c.Region)},
		)
		uploader := s3manager.NewUploader(sess)

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			ContentType: aws.String(optRequestInput.ContentType),
			Body:        reader,
			Metadata:    optRequestInput.MetaData,
		})
		if err != nil {
			log.Printf("Unable to upload %q to %q, %v", key, bucket, err)
			return err
		}
		log.Printf("Successfully uploaded %q to %q\n", key, bucket)
		return nil
	}
}

// PresignObject default implementation of S3 object URL presigner
func PresignObject(c Config) PresignFunc {
	return func(bucket, key string) (string, error) {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(c.Region)},
		)

		svc := s3.New(sess)
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		urlStr, err := req.Presign(15 * time.Minute)
		if err != nil {
			return "", err
		}

		return urlStr, nil
	}
}
