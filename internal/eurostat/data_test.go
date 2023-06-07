package eurostat

import (
	"testing"
	"time"
)

func TestTimestampFromFilename(t *testing.T) {
	type TestCase struct {
		fileName   string
		want       time.Time
		shouldPass bool
	}

	cases := []TestCase{
		{fileName: "/foo/bar/20210112T102331.tsv.gz", want: time.Date(2021, 1, 12, 10, 23, 31, 0, time.UTC)},
		{fileName: "20210112T102331.tsv.gz", want: time.Date(2021, 1, 12, 10, 23, 31, 0, time.UTC)},
		{fileName: "/foo/bar/20210113T142331.tsv.gz", want: time.Date(2021, 1, 13, 14, 23, 31, 0, time.UTC)},
	}

	for _, c := range cases {
		got, err := timestampFromFileName(c.fileName)
		if err != nil {
			t.Fatalf("Expected error to be nil but got %s\n", err)
		}
		if got != c.want {
			t.Fatalf("wanted %s but got %s", c.want, got)
		}
	}
}
