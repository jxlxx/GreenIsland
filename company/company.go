package company

import "github.com/jxlxx/GreenIsland/config"

func New(fileName string) *Company {
	company := Company{}
	config.ReadConfig(fileName, &company)
	return &company
}

func (c Company) QuarterlyReport() error {
	return nil
}

func (c Company) DailyJitter() error {
	return nil
}
