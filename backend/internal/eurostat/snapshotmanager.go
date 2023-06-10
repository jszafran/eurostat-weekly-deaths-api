package eurostat

type SnapshotManager struct {
	s3Bucket string
}

func New(s3Bucket string) SnapshotManager {
	return SnapshotManager{s3Bucket: s3Bucket}
}

func (sm SnapshotManager) ReadSnapshotFromLocalDisk(path string) DataSnapshot {
	return DataSnapshot{}
}

func (sm SnapshotManager) WriteSnapshotToLocalDisk(snapshot DataSnapshot) error {
	return nil
}

func (sm SnapshotManager) ReadSnapshotFromS3(s3Path string) DataSnapshot {
	return DataSnapshot{}
}

func (sm SnapshotManager) ReadMostRecentSnapshotFromS3() (DataSnapshot, error) {
	return DataSnapshot{}, nil
}

func (sm SnapshotManager) WriteSnapshotToS3(snapshot DataSnapshot) error {
	return nil
}

func (sm SnapshotManager) ReadSnapshotFromEurostat() (DataSnapshot, error) {
	return DataSnapshot{}, nil
}
