package main

import (
	"regexp"
	"sort"
)

type By func(p1, p2 *ServerEntry) bool

func (by By) Sort(servers []*ServerEntry) {
	bs := &serverSorter{
		servers: servers,
		by:      by,
	}
	sort.Sort(bs)
}

type serverSorter struct {
	servers []*ServerEntry
	by      func(s1, s2 *ServerEntry) bool
}

func (s *serverSorter) Len() int {
	return len(s.servers)
}

func (s *serverSorter) Swap(i, j int) {
	s.servers[i], s.servers[j] = s.servers[j], s.servers[i]
}

func (s *serverSorter) Less(i, j int) bool {
	return s.by(s.servers[i], s.servers[j])
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func kvFromString(a string) (string, string) {
	re := regexp.MustCompile(`^(\w+)=(\w+)$`)
	if m := re.FindAllStringSubmatch(a, -1); len(m) != 0 {
		return m[0][1], m[0][2]
	}
	return "", ""
}

func failoverOrder(s1, s2 *ServerEntry) bool {
	return s1.Attrs["ORDER"] < s2.Attrs["ORDER"]
}
