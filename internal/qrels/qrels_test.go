package qrels

import (
	"strings"
	"testing"
)

const file = `1	1	1
1	2	1
1	3	0
1	4	0
1	5	0
1	6	0
1	7	0
1	8	0
1	9	1
1	10	0
1	11	1
1	12	0
1	13	0
1	14	0
1	15	1
1	16	0
1	17	0
1	18	0
1	19	0
1	20	1
1	21	1
1	22	1
2	1	1
2	2	0
2	3	1
2	4	0
2	5	0
2	6	0
2	7	0
2	8	0
2	9	1
2	10	1
`

func TestParseQrelsFile(t *testing.T) {
	r := strings.NewReader(file)
	qrels, err := ParseQrelsFile(r)
	if err != nil {
		t.Fatalf("got error while parsing: %v", err)
	}
	want := 32
	got := len(qrels)
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	t.Run("Create RelevantDocs", func(t *testing.T) {
		docs := CreateMap(qrels)
		t.Log(docs)
	})
}
