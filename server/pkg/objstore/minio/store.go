package minio

import (
	"context"
	"fmt"
	"io"
	"strings"

	minio "github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/encrypt"
	"github.com/minio/minio-go/v6/pkg/s3utils"
)

type (
	Options struct {
		Endpoint string
		Secure   bool
		Strict   bool

		AccessKeyID     string
		SecretAccessKey string

		ServerSideEncryptKey []byte
	}

	minioClient interface {
		BucketExists(bucketName string) (bool, error)
		MakeBucket(bucketName string, location string) (err error)
		PutObject(bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (n int64, err error)
		RemoveObject(bucketName, objectName string) error
		GetObject(bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	}

	store struct {
		bucket     string
		pathPrefix string
		component  string

		mc  minioClient
		sse encrypt.ServerSide

		originalFn func(id uint64, ext string) string
		previewFn  func(id uint64, ext string) string
	}
)

var (
	defPreviewFn = func(id uint64, ext string) string {
		return fmt.Sprintf("%d_preview.%s", id, ext)
	}

	defOriginalFn = func(id uint64, ext string) string {
		return fmt.Sprintf("%d.%s", id, ext)
	}
)

func New(bucket, pathPrefix, component string, opt Options) (s *store, err error) {
	client, err := minio.New(opt.Endpoint, opt.AccessKeyID, opt.SecretAccessKey, opt.Secure)
	if err != nil {
		return nil, err
	}
	return newWithClient(client, bucket, pathPrefix, component, opt)
}

func newWithClient(mc minioClient, bucket, pathPrefix, component string, opt Options) (s *store, err error) {
	s = &store{
		bucket:     bucket,
		pathPrefix: pathPrefix,
		component:  component,
		mc:         mc,

		originalFn: defOriginalFn,
		previewFn:  defPreviewFn,
	}

	if err = s3utils.CheckValidBucketName(s.bucket); err != nil {
		return nil, err
	}

	if e, err := s.mc.BucketExists(s.bucket); err != nil {
		return nil, err
	} else if !e {
		if opt.Strict {
			return nil, fmt.Errorf("bucket %q does not exist", s.bucket)
		}

		err = s.mc.MakeBucket(s.bucket, "us-east-1")
		if err != nil {
			return nil, err
		}
	}

	if len(opt.ServerSideEncryptKey) > 0 {
		s.sse, err = encrypt.NewSSEC(opt.ServerSideEncryptKey)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (s *store) check(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("Invalid name when trying to store object: '%s' (for %s)", name, s.bucket)
	}

	return nil
}

func (s store) Original(id uint64, ext string) string {
	// @todo presigned URL
	return s.originalFn(id, ext)
}

func (s store) Preview(id uint64, ext string) string {
	// @todo presigned URL
	return s.previewFn(id, ext)

}

func (s store) Save(name string, f io.Reader) (err error) {
	_, err = s.mc.PutObject(s.bucket, s.getObjectName(name), f, -1, minio.PutObjectOptions{
		ServerSideEncryption: s.sse,
	})

	return err
}

func (s store) Remove(name string) error {
	return s.mc.RemoveObject(s.bucket, s.getObjectName(name))
}

func (s store) Open(name string) (io.ReadSeekCloser, error) {
	return s.mc.GetObject(s.bucket, s.getObjectName(name), minio.GetObjectOptions{
		ServerSideEncryption: s.sse,
	})
}

func (s *store) Healthcheck(_ context.Context) error {
	return nil
}

// getObjectName prefix path to object name
func (s *store) getObjectName(name string) (out string) {
	path := strings.Replace(s.pathPrefix, "{component}", s.component, 1)
	return fmt.Sprintf("%s%s", path, name)
}

// GetBucket return bucket name based on storage option bucket, separator or bucketName
func GetBucket(bucket, component string) string {
	return strings.Replace(bucket, "{component}", component, 1)
}
