package etcd

import (
	"errors"
	"math/rand"
)

type RandomBalance struct {
	curIndex int
	rss      []string
}

var (
	Random = "roundrobin"
)

func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

func (r *RandomBalance) Get() (string, error) {
	return r.Next(), nil
}

func (r *RandomBalance) Remove(ip string) {
	for i, v := range r.rss {
		if v == ip {
			r.rss = append(r.rss[:i], r.rss[i+1:]...)
		}
	}
}
