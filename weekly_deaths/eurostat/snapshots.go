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
	"sort"
	"strconv"
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

func sortSnapshotKeys(keys []string) ([]string, error) {
	type snapshotKey struct {
		timestamp time.Time
		key       string
	}

	var (
		parsedKeys []snapshotKey
		sortedKeys []string
	)

	for _, k := range keys {
		ts, err := parseTimestamp(k)
		if err != nil {
			log.Printf("unparsable object name %s: %s", k, err)
			continue
		}

		parsedKeys = append(parsedKeys, snapshotKey{
			timestamp: ts,
			key:       k,
		})
	}

	if len(parsedKeys) == 0 {
		return sortedKeys, ErrNoParsableObjectsInBucket
	}

	sort.Slice(parsedKeys, func(i, j int) bool {
		return parsedKeys[i].timestamp.Before(parsedKeys[j].timestamp)
	})

	for _, k := range parsedKeys {
		sortedKeys = append(sortedKeys, k.key)
	}

	return sortedKeys, nil
}

func latestKey(keys []string) (string, error) {
	var latestKey string

	if len(keys) == 0 {
		return latestKey, ErrSnapshotsBucketEmpty
	}

	sk, err := sortSnapshotKeys(keys)
	if err != nil {
		return latestKey, err
	}

	return sk[len(sk)-1], nil
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

func (sm *SnapshotManager) listSnapshotsChronologically() ([]string, error) {
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
		return obj, err
	}

	sorted, err := sortSnapshotKeys(obj)
	if err != nil {
		return obj, err
	}

	return sorted, nil
}

func (sm *SnapshotManager) LatestSnapshotFromS3() (DataSnapshot, error) {
	var ds DataSnapshot

	log.Println("Attempting to fetch latest snapshot from S3.")
	obj, err := sm.listSnapshotsChronologically()
	if err != nil {
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

func (sm *SnapshotManager) CleanupSnapshots() error {
	t := os.Getenv("CLEANUP_KEEP_N_LATEST_SNAPSHOTS")
	threshold, err := strconv.Atoi(t)
	if err != nil {
		return fmt.Errorf("failed to read cleanup threshold env var: %w", err)
	}

	keys, err := sm.listSnapshotsChronologically()
	delta := len(keys) - threshold
	if delta < 0 {
		log.Printf("exiting cleanup operation early as there are only %d snapshots and threshold is %d", len(keys), threshold)
		return nil
	}

	keysToDelete := keys[:delta]

	// TODO: implement actual deletion of objects from S3
	log.Printf("This function would have deleted %+v snapshots!\n", keysToDelete)
	return nil
}
