package infrastructure

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	ErrMetadataNotFound = errors.New("metadata not found")
)

type S3Impl struct {
	accessKey string
	secretKey string
	region    string
	bucket    string

	sess *session.Session
	s3   *s3.S3
}

func NewS3Impl(accessKey, secretKey, region, bucket, endpoint string) (*S3Impl, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)

	return &S3Impl{
		accessKey: accessKey,
		secretKey: secretKey,
		region:    region,
		bucket:    bucket,
		sess:      sess,
		s3:        svc,
	}, nil
}

func (s *S3Impl) UploadFile(filePath, destKey string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return err
	}

	_, err = s.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(destKey),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Impl) UploadFileFromMemory(data []byte, destKey string) error {
	uploader := s3manager.NewUploader(s.sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(destKey),
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Impl) DownloadFile(filePath, sourceKey string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(s.sess)

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(sourceKey),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == s3.ErrCodeNoSuchKey {
				return ErrMetadataNotFound
			}
		}
		return err
	}

	return nil
}

func (s *S3Impl) DownloadFileInMemory(sourceKey string) ([]byte, error) {
	downloader := s3manager.NewDownloader(s.sess)

	buff := &aws.WriteAtBuffer{}
	_, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(sourceKey),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == s3.ErrCodeNoSuchKey {
				return nil, ErrMetadataNotFound
			}
		}
		return nil, err
	}

	return buff.Bytes(), nil
}

// ListFiles lists all files names (only the name, not path) in a directory
func (s *S3Impl) ListFiles(prefix string) ([]string, error) {
	svc := s3.New(s.sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == s3.ErrCodeNoSuchKey {
				return nil, ErrBackupNotFound
			}
		}

		return nil, err
	}

	var keys []string
	for _, item := range resp.Contents {
		keys = append(keys, *item.Key)
	}

	return keys, nil
}

// ListFolders lists all folders in a directory and retruns their full path
func (s *S3Impl) ListFolders(prefix string) ([]string, error) {
	prefix = s.formatFolderKey(prefix)
	svc := s3.New(s.sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}

	var folders []string
	for _, item := range resp.CommonPrefixes {
		if *item.Prefix == prefix {
			continue
		}
		folders = append(folders, *item.Prefix)
	}

	return folders, nil
}

func (s *S3Impl) IsFolderExists(prefix string) (bool, error) {
	svc := s3.New(s.sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return false, err
	}

	return len(resp.CommonPrefixes) > 0, nil
}

func (s *S3Impl) DeleteFolder(prefix string) error {
	prefix = s.formatFolderKey(prefix)
	svc := s3.New(s.sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return err
	}

	var keys []string
	for _, item := range resp.Contents {
		keys = append(keys, *item.Key)
	}

	if len(keys) > 0 {
		_, err = svc.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(s.bucket),
			Delete: &s3.Delete{
				Objects: make([]*s3.ObjectIdentifier, len(keys)),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *S3Impl) DeleteFile(key string) error {
	svc := s3.New(s.sess)
	return s.deleteFile(svc, key)
}

func (s *S3Impl) deleteFile(scv *s3.S3, key string) error {
	_, err := scv.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Impl) formatFolderKey(key string) string {
	if key != "" && !strings.HasSuffix(key, "/") {
		key += "/"
	}
	return key
}
