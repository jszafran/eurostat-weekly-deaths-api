package eurostat

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
	"os"
	"time"
)

var (
	ErrSnapshotsBucketEmpty      = errors.New("snapshots bucket is empty")
	ErrNoParsableObjectsInBucket = errors.New("no objects with parsable names found in S3")
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
		Key:    aws.String(fmt.Sprintf("%s%s", timestamp.Format(timestampLayout), dataFileExtension)),
		Body:   r,
	})
	if err != nil {
		return err
	}
	return nil
}

func (sm *SnapshotManager) GetSnapshot(key string) (DataSnapshot, error) {
	var (
		ds   DataSnapshot
		buff aws.WriteAtBuffer
	)

	s3Client := s3.New(sm.session)
	downloader := s3manager.NewDownloaderWithClient(s3Client)
	_, err := downloader.Download(&buff, &s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(sm.bucket),
	})

	if err != nil {
		return ds, err
	}

	ts, err := parseTimestamp(key)
	if err != nil {
		return ds, err
	}

	r, err := gzip.NewReader(bytes.NewReader(buff.Bytes()))
	if err != nil {
		return ds, err
	}

	data, err := ParseData(r)
	ds.Data = data
	ds.Timestamp = ts

	return ds, nil
}

func latestKey(keys []string) (string, error) {
	var (
		latestKey string
		maxIx     int
		maxTs     time.Time
	)

	switch n := len(keys); n {
	case 0:
		return latestKey, ErrSnapshotsBucketEmpty

	case 1:
		_, err := parseTimestamp(keys[0])
		if err != nil {
			//return latestKey, fmt.Errorf("the only object in s3 has unparsable name: %s", keys[0])
			return latestKey, fmt.Errorf("unparsable object name %s: %w", keys[0], err)
		}

		return keys[0], nil

	default:
		for i, k := range keys {
			ts, err := parseTimestamp(k)
			if err != nil {
				log.Printf("found object with unparsable name in s3 bucket: %s", k)
				continue
			}

			if ts.After(maxTs) {
				maxIx = i
				maxTs = ts
			}
		}
	}

	if maxTs.IsZero() {
		return latestKey, ErrNoParsableObjectsInBucket
	}

	return keys[maxIx], nil
}

func (sm *SnapshotManager) LatestSnapshotFromS3() (DataSnapshot, error) {
	var ds DataSnapshot

	log.Println("Attempting to fetch latest snapshot from S3.")
	obj := make([]string, 0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	s3Client := s3.New(sm.session)

	callbackFn := func(o *s3.ListObjectsOutput, b bool) bool {
		for _, o := range o.Contents {
			obj = append(obj, *o.Key)
		}
		return true
	}

	if err := s3Client.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(sm.bucket),
		Prefix: aws.String(""),
	}, callbackFn); err != nil {
		return ds, err
	}

	lk, err := latestKey(obj)
	if err != nil {
		return ds, err
	}

	ds, err = sm.GetSnapshot(lk)
	if err != nil {
		return ds, err
	}
	log.Printf("Successfully fetched S3 snapshot (timestamp: %s)\n", ds.Timestamp)
	return ds, nil
}
