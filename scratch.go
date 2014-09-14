package main

import (
	"fmt"
	"net"
	"regexp"
	"sort"
)

type ServerEntry struct {
	ID      string
	Service string
	Tags    []string
	Port    int
	IP      []byte
	Node    string
}

type By func(p1, p2 *ServerEntry) bool

func (by By) Sort(servers []ServerEntry) {
	bs := &serverSorter{
		servers: servers,
		by:      by,
	}
	sort.Sort(bs)
}

type serverSorter struct {
	servers []ServerEntry
	by      func(s1, s2 *ServerEntry) bool
}

func (s *serverSorter) Len() int {
	return len(s.servers)
}

func (s *serverSorter) Swap(i, j int) {
	s.servers[i], s.servers[j] = s.servers[j], s.servers[i]
}

func (s *serverSorter) Less(i, j int) bool {
	return s.by(&s.servers[i], &s.servers[j])
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

func main() {
	servers := []ServerEntry{
		{"0", "app", []string{"backup"}, 8000, net.IP{192, 168, 0, 2}, "node0"},
		{"1", "app", []string{}, 8000, net.IP{192, 168, 0, 3}, "node1"},
		{"2", "app", []string{"backup"}, 8000, net.IP{192, 168, 0, 4}, "node3"},
	}
	backups_last := func(s1, s2 *ServerEntry) bool {
		return !stringInSlice("backup", s1.Tags)
	}
	fmt.Println("Unsorted:", servers)
	By(backups_last).Sort(servers)
	fmt.Println("Sorted, backups last:", servers)

	for _, tag := range []string{"FOO=BAR", "FOO", "FOO=", "=BAZ", "123=456"} {
		k, v := kvFromString(tag)
		if k != "" {
			fmt.Println(k, v)
		}
	}
}
