package data

import (
	"strconv"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/shopspring/decimal"

	"github.com/ardizariat/generate-pdf/internal/locale"
)

// moneyRange describes a deterministic fake amount. Values are generated in
// minor units and divided by div (100 for USD cents, 1 for whole rupiah).
type moneyRange struct {
	min int
	max int
	div int64
}

func (r moneyRange) fake(faker *gofakeit.Faker) decimal.Decimal {
	return decimal.NewFromInt(int64(faker.IntRange(r.min, r.max))).Div(decimal.NewFromInt(r.div))
}

// dataset is the language-specific reference data used by the generator.
type dataset struct {
	firstNames    []string
	lastNames     []string
	cities        []string
	streets       []string
	companies     []string
	goods         []string
	units         []string
	positions     []string
	purposes      []string
	deliveryNotes []string
	prose         []string

	intro              string
	fictional          string
	phonePattern       string
	addressNumberFirst bool
	price              moneyRange
	receipt            moneyRange
	report             moneyRange
}

func dataFor(lang locale.Language) dataset {
	if lang == locale.Indonesian {
		return indonesianData
	}
	return englishData
}

func (d dataset) address(faker *gofakeit.Faker) string {
	number := strconv.Itoa(faker.IntRange(1, 200))
	street := faker.RandomString(d.streets)
	if d.addressNumberFirst {
		return number + " " + street
	}
	return street + " No. " + number
}
