package minio

import (
	"context"
	"fmt"
	"io"

	minio "github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/encrypt"
	"github.com/minio/minio-go/v6/pkg/s3utils"
	"github.com/pkg/errors"
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

	store struct {
		bucket   string
		separator string
		path      string

		mc  *minio.Client
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

func New(bucket string, separator string, path string, opt Options) (s *store, err error) {
	s = &store{
		bucket:   bucket,
		separator: separator,
		path:      path,
		mc:        nil,

		originalFn: defOriginalFn,
		previewFn:  defPreviewFn,
	}

	if err = s3utils.CheckValidBucketName(s.bucket); err != nil {
		return nil, err
	}

	if s.mc, err = minio.New(opt.Endpoint, opt.AccessKeyID, opt.SecretAccessKey, opt.Secure); err != nil {
		return nil, err
	}

	if e, err := s.mc.BucketExists(s.bucket); err != nil {
		return nil, err
	} else if !e {
		if opt.Strict {
			return nil, errors.Errorf("bucket %q does not exist", s.bucket)
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
		return errors.Errorf("Invalid name when trying to store object: '%s' (for %s)", name, s.bucket)
	}

	return nil
}

func (s *store) withPrefix(name string) string {
	return s.path + s.separator + name
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
	_, err = s.mc.PutObject(s.bucket, s.withPrefix(name), f, -1, minio.PutObjectOptions{
		ServerSideEncryption: s.sse,
	})

	return err
}

func (s store) Remove(name string) error {
	return s.mc.RemoveObject(s.bucket, s.withPrefix(name))
}

func (s store) Open(name string) (io.ReadSeeker, error) {
	return s.mc.GetObject(s.bucket, s.withPrefix(name), minio.GetObjectOptions{
		ServerSideEncryption: s.sse,
	})
}

func (s *store) Healthcheck(ctx context.Context) error {
	return nil
}
