# Goq

```go
func main() {
	db, err := goq.Open("postgres", conn)
	panicIf(err)
	b := NewBuilder(db.Dialect())

	query := b.Select(
		b.Countries.ID,
		b.Countries.Name,
		b.Cities.All(),
	).From(
		b.Countries,
	).Joins(
		b.InnerJoin(b.Cities).On(b.Cities.CountryID.Eq(b.Countries.ID)),
	).Where(
		b.Countries.Population.Lte(500000),
		b.Cities.Name.Like("New%"),
	).OrderBy(
		b.Countries.Name,
		b.Cities.Name,
	)

	var countries []Country
	var citiesByCountry map[int][]City
	err = db.Query(query).Collect(
		b.Countries.ToSlice(&countries),
		b.Cities.ToSliceMap(&cities).By(b.Countries.ID),
	)
	panicIf(err)

	fmt.Println(countries, cities)
}
```
