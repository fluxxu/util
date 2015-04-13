package util

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
	"reflect"
	"strings"
)

type SqlProvider struct {
	*BaseProvider
	table  string
	dbx    *sqlx.DB
	fields []string
}

func (p *SqlProvider) ParseStruct(s interface{}) {
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		dbtag := strings.TrimSpace(f.Tag.Get("db"))
		if dbtag != "" && dbtag != "-" {
			p.fields = append(p.fields, "`"+dbtag+"`")
			//sort, filter
			tag := strings.TrimSpace(f.Tag.Get("provider"))
			for _, attr := range strings.Split(tag, " ") {
				attr = strings.TrimSpace(attr)
				if attr == "sort" {
					p.sortKeys[dbtag] = true
				} else if attr == "filter" {
					p.filterKeys[dbtag] = true
				}
			}
		}
	}
}

type wherePart struct {
	pred interface{}
	args []interface{}
}

func (p *SqlProvider) filterWhereParts() []wherePart {
	var parts []wherePart
	for k, vs := range p.GetFilters() {
		var part wherePart
		var preds []string
		for _, v := range vs {
			preds = append(preds, fmt.Sprintf("(`%s` LIKE ?)", k))
			part.args = append(part.args, "%"+strings.Replace(v, "%", "", -1)+"%")
		}
		part.pred = "(" + strings.Join(preds, " OR ") + ")"
		parts = append(parts, part)
	}
	return parts
}

func (p *SqlProvider) count() (int, error) {
	var count int
	q := squirrel.Select("COUNT(*)").From(p.table)
	for _, part := range p.filterWhereParts() {
		q = q.Where(part.pred, part.args...)
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return 0, err
	}
	//fmt.Println("Count:" + sql)
	if err = p.dbx.Get(&count, sql, args...); err != nil {
		return 0, err
	}
	return count, nil
}

func (p *SqlProvider) Read(dst interface{}) (*ProviderResponse, error) {
	r := &ProviderResponse{}
	q := squirrel.Select(p.fields...).From(p.table).Limit(p.take).Offset(p.skip)
	for _, part := range p.filterWhereParts() {
		q = q.Where(part.pred, part.args...)
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	//fmt.Println("Select:" + sql)
	if err = p.dbx.Select(dst, sql, args...); err != nil {
		return nil, err
	}

	r.Data = dst

	if p.stat {
		count, err := p.count()
		if err != nil {
			return nil, err
		}

		r.Stat = map[string]int{
			"count": count,
			"take":  int(p.take),
			"skip":  int(p.skip),
		}
	}
	return r, nil
}

func NewSqlProvider(dbx *sqlx.DB, table string, s interface{}) *SqlProvider {
	p := &SqlProvider{
		table: "`" + table + "`",
		dbx:   dbx,
	}
	p.BaseProvider = NewBaseProvider()
	p.ParseStruct(s)
	//fmt.Println("db:", p.fields)
	//fmt.Println("filter:", p.filterKeys)
	//fmt.Println("sort:", p.sortKeys)
	return p
}
