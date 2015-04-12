package util

import (
	"net/http"
	"regexp"
	"strconv"
)

const MaxNumOfFilterValue int = 5
const MaxTake uint64 = 30

type Filters map[string][]string

type Pager interface {
	SetTake(v uint)
	SetSkip(v uint)
}

type Filter interface {
	LimitFilterKeys(keys ...string)
	AddFilter(key, value string)
	GetFilters() Filters
}

type Sorter interface {
	LimitSortKeys(keys ...string)
	SetSort(key string)
	GetSort() string
}

type ProviderResponse struct {
	Stat map[string]int `json:"stat"`
	Data interface{}    `json:"data"`
}

type Provider interface {
	Pager
	Filter
	Sorter
	ParseRequest(r *http.Request)

	Read(dst interface{}) (*ProviderResponse, error)
}

type BaseProvider struct {
	stat       bool
	take       uint64
	skip       uint64
	filterKeys map[string]bool
	sortKeys   map[string]bool

	filters Filters
	sort    string
}

func (p *BaseProvider) SetTake(v uint64) {
	if v > MaxTake {
		v = MaxTake
	}
	p.take = v
}

func (p *BaseProvider) SetSkip(v uint64) {
	p.skip = v
}

func (p *BaseProvider) AddFilter(k, v string) {
	if _, ok := p.filterKeys[k]; ok {
		s := p.filters[k]
		if len(s) < MaxNumOfFilterValue {
			for _, item := range s {
				if item == v {
					return
				}
			}
			p.filters[k] = append(s, v)
		}
	}
}

func (p *BaseProvider) GetFilters() Filters {
	return p.filters
}

func (p *BaseProvider) SetStat(v bool) {
	p.stat = v
}

func (p *BaseProvider) SetSort(k string) {
	if _, ok := p.sortKeys[k]; ok {
		p.sort = k
	}
}

func (p *BaseProvider) GetSort() string {
	return p.sort
}

func (p *BaseProvider) LimitFilterKeys(keys ...string) {
	for _, v := range keys {
		p.filterKeys[v] = true
	}
}

func (p *BaseProvider) LimitSortKeys(keys ...string) {
	for _, v := range keys {
		p.sortKeys[v] = true
	}
}

var filterR = regexp.MustCompile(`^filter_([\w]+)`)

func (p *BaseProvider) ParseRequest(r *http.Request) {
	params := r.URL.Query()
	for k, values := range params {
		if len(values) == 0 {
			continue
		}

		lv := values[len(values)-1]

		if k == "sort" {
			p.SetSort(lv)
		} else if k == "stat" {
			p.SetStat(true)
		} else if k == "take" {
			if iv, err := strconv.Atoi(lv); err == nil {
				p.SetTake(uint64(iv))
			}
		} else if k == "skip" {
			if iv, err := strconv.Atoi(lv); err == nil {
				p.SetSkip(uint64(iv))
			}
		} else if m := filterR.FindSubmatch([]byte(k)); m != nil {
			fk := string(m[1])
			for _, v := range values {
				p.AddFilter(fk, v)
			}
		}
	}
}

func NewBaseProvider() *BaseProvider {
	p := &BaseProvider{
		filterKeys: make(map[string]bool),
		sortKeys:   make(map[string]bool),
		filters:    make(Filters),
		take:       MaxTake,
	}
	return p
}
