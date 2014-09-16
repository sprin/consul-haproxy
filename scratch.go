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
  Attrs   map[string]string
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
		{"0", "app", []string{"HOST=app", "ORDER=1"}, 8000, net.IP{192, 168, 0, 2}, "node0", map[string]string{},},
		{"1", "app", []string{"HOST=app", "ORDER=0"}, 8000, net.IP{192, 168, 0, 3}, "node1", map[string]string{},},
		{"2", "db", []string{"HOST=db", "ORDER=0"}, 8000, net.IP{192, 168, 0, 4}, "node3", map[string]string{},},
		{"3", "app", []string{"HOST=app", "ORDER=2"}, 8000, net.IP{192, 168, 0, 5}, "node4", map[string]string{},},
	}

  for _, server := range servers {
    for _, tag := range server.Tags {
      k, v := kvFromString(tag)
      if k != "" {
        server.Attrs[k] = v;
      }
    }
  }

	failover_order := func(s1, s2 *ServerEntry) bool {
		return s1.Attrs["ORDER"] < s2.Attrs["ORDER"]
	}

	fmt.Println("Unsorted:", servers)
	By(failover_order).Sort(servers)
	fmt.Println("Sorted:", servers)
}
