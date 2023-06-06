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
		{fileName: "/foo/bar/20210112T102331.tsv.gz", want: time.Date(2021, 1, 12, 10, 23, 31, 0, time.UTC), shouldPass: true},
		{fileName: "20210112T102331.tsv.gz", want: time.Date(2021, 1, 12, 10, 23, 31, 0, time.UTC), shouldPass: true},
		{fileName: "/foo/bar/20210113T142331.tsv.gz", want: time.Date(2021, 1, 13, 14, 23, 31, 0, time.UTC), shouldPass: true},
	}

	for i, c := range cases {
		ix := i + 1
		got, err := timestampFromFileName(c.fileName)
		if err != nil {
			t.Fatalf("(%d) Expected error to be nil but got %s\n", ix, err)
		}
		if got != c.want {
			t.Fatalf("(%d) wanted %s but got %s", ix, c.want, got)
		}
	}
}
