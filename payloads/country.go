package payloads

type QuarterlyCountryUpdate struct {
	Name              string
	Quarter           int
	TotalPopulation   int
	WorkingPopulation int
}
