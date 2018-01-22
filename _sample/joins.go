package main

import "github.com/ryym/goq/gql"

func (p *Prefs) Cities(c *Cities) *gql.Join {
	return &gql.Join{p, p.ID.Eq(c.PrefectureID)}
}
