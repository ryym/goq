🚧  UNDER DEVELOPMENT

# Goq  

[![circleci](https://circleci.com/gh/ryym/goq.svg?style=svg)](https://circleci.com/gh/ryym/goq)

[![appveyor](https://ci.appveyor.com/api/projects/status/5yf0rg3n810cnkih?svg=true)](https://ci.appveyor.com/project/ryym/goq)

SQL-based DB access library for Gophers

## Features

### SQL-based API

Goq provides the low-level API which just wraps SQL clauses as Go methods,
Instead of abstracting a way of query construction as an opinionated API like typical frameworks.
That is, you already know most of the Goq API if you know SQL.

### Flexible Result Mapping

Using Goq, you can collect result rows fetched from DB into various format structures:
a slice of your model, single map, map of slices, combination of them, etc.

### Custom Query Builder Generation

Goq can generate your custom query builder by [`go generate`](https://blog.golang.org/generate)
based on your models mapped to DB tables.
This helps you write a query with more type safety and readability.

## What does it look like?

```go
import (
    "fmt"

    _ "github.com/lib/pq"
    "github.com/ryym/goq"
)

func main() {
    // Connect to DB.
    db, err := goq.Open("postgres", conn)
    panicIf(err)
    defer db.Close()

    // Initialize your builder.
    q := NewBuilder(db.Dialect())

    // Write a query.
    query := q.Select(
        q.Countries.ID,
        q.Countries.Name,
        q.Cities.All(),
    ).From(
        q.Countries,
    ).Joins(
        q.InnerJoin(q.Cities).On(
            q.Cities.CountryID.Eq(q.Countries.ID),
        ),
    ).Where(
        q.Countries.Population.Lte(500000),
        q.Cities.Name.Like("New%"),
    ).OrderBy(
        q.Countries.Name,
        q.Cities.Name,
    )

    var countries []Country
    var citiesByCountry map[int][]City

    // Collect results.
    err = db.Query(query).Collect(
        q.Countries.ToUniqSlice(&countries),
        q.Cities.ToSliceMap(&citiesByCountry).By(q.Countries.ID),
    )
    panicIf(err)

    fmt.Println("Complete!", countries, citiesByCountry)
}
```

## API

<https://godoc.org/github.com/ryym/goq>

## Out of Scope Features

Goq is not a DB management framework so does not support any of these:

- schema migration
- schema generation from Go code
- model generation from schema

## Inspiration

Goq is inspired by [ScalikeJDBC](http://scalikejdbc.org/)
which is a Scala library providing SQL-based API.
