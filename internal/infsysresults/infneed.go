package infsysresults

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Result struct {
	Need  string
	DocID string
	Pos   int
}

type ResultMap map[string][]Result

func CreateMap(res []Result) ResultMap {
	m := make(map[string][]Result)
	for _, r := range res {
		m[r.Need] = append(m[r.Need], r)
	}
	return m
}

func ParseResults(r io.Reader) (needs []Result, err error) {
	line := 0
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line++
		infNeed, err := parseLine(sc.Text())
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error in line %d formating needs: %v", line, err))
		}
		infNeed.Pos = line
		needs = append(needs, infNeed)
	}
	return needs, nil
}

func parseLine(line string) (Result, error) {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return Result{}, errors.New("result hasn't the correct amount of fields")
	}
	return Result{
		Need:  fields[0],
		DocID: fields[1],
	}, nil
}
