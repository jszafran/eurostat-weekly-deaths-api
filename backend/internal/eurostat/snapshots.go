package eurostat

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
	"time"
)

type SnapshotManager struct {
	bucket  string
	session *session.Session
}

func NewSnapshotManager(bucket string) (SnapshotManager, error) {
	var sm SnapshotManager

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			),
		})

	if err != nil {
		return sm, err
	}
	return SnapshotManager{
		bucket:  bucket,
		session: sess,
	}, nil
}

func (sm *SnapshotManager) PersistSnapshot(r io.Reader, timestamp time.Time) error {
	uploader := s3manager.NewUploader(sm.session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(sm.bucket),
		Key:    aws.String(fmt.Sprintf("%s.gzip", timestamp.Format("20060102T150405"))),
		Body:   r,
	})
	if err != nil {
		return err
	}
	return nil
}
