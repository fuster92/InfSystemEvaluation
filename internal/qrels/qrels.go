package qrels

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Doc struct {
	DocID string
	Rel   string
}

type Qrel struct {
	InfNeed string
	D       Doc
}

type RelevantDocs map[string]map[string]Doc

func (r RelevantDocs) IsRelevant(docID, infNeed string) bool {
	doc, ok := r[infNeed]
	if !ok {
		return false
	}
	rel, ok := doc[docID]
	if !ok {
		return false
	}

	return rel.Rel == "1"
}

func CreateMap(qrels []Qrel) RelevantDocs {
	m := make(map[string]map[string]Doc)
	for _, qrel := range qrels {
		if _, ok := m[qrel.InfNeed]; !ok {
			m[qrel.InfNeed] = make(map[string]Doc)
		}
		m[qrel.InfNeed][qrel.D.DocID] = qrel.D
	}
	return m
}

func ParseQrelsFile(r io.Reader) (qrels []Qrel, err error) {
	line := 0
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line++
		qrel, err := parseLine(sc.Text())
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error in line %d formating qrel: %v", line, err))
		}
		qrels = append(qrels, qrel)
	}
	return qrels, nil
}

func parseLine(line string) (Qrel, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return Qrel{}, errors.New("qrel has not the correct amount of fields")
	}
	return Qrel{
		InfNeed: fields[0],
		D: Doc{
			DocID: fields[1],
			Rel:   fields[2],
		},
	}, nil
}

func (q Qrel) String() string {
	return fmt.Sprintf("%s\t%s\t%s", q.InfNeed, q.D.DocID, q.D.Rel)
}
