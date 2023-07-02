package eurostat

import (
	"errors"
	"testing"
)

func TestLatestKey(t *testing.T) {
	type testCase struct {
		input []string
		want  string
		err   error
	}

	cases := []testCase{
		{
			input: []string{"20210112T012012.tsv.gz", "20240101T000000.tsv.gz"},
			want:  "20240101T000000.tsv.gz",
			err:   nil,
		},
		{
			input: []string{},
			want:  "",
			err:   ErrSnapshotsBucketEmpty,
		},
		{
			input: []string{"20210112T012012.tsv.gz"},
			want:  "20210112T012012.tsv.gz",
			err:   nil,
		},
		{
			input: []string{"20210112T012012,.tsv.gz", "20240101T000000,.tsv.gz"},
			want:  "",
			err:   ErrNoParsableObjectsInBucket,
		},
		{
			input: []string{"20210112T012012.tsv.gz", "20200112T012012.tsv.gz", "20220112T012012.tsv.gz", "20250112T012012.tsv.gz"},
			want:  "20250112T012012.tsv.gz",
			err:   nil,
		},
	}

	for _, tc := range cases {
		got, err := latestKey(tc.input)
		if !errors.Is(err, tc.err) {
			t.Fatalf("Expected err %v but got err %v\n", tc.err, err)
		}
		if got != tc.want {
			t.Fatalf("Expected %s but got %s\n", tc.want, got)
		}
	}
}
