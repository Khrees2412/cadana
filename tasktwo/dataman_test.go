package tasktwo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalPersons(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"id": "1",
				"personName": "John Doe",
				"salary": {
					"value": 5000,
					"currency": "USD"
				}
			},
			{
				"id": "2",
				"personName": "Jane Doe",
				"salary": {
					"value": 6000,
					"currency": "EUR"
				}
			}
		]
	}`

	persons, err := UnmarshalPersons([]byte(jsonData))
	assert.Nil(t, err)
	assert.Len(t, persons.Data, 2)
	assert.Equal(t, "John Doe", persons.Data[0].PersonName)
	assert.Equal(t, "Jane Doe", persons.Data[1].PersonName)
	assert.Equal(t, int64(5000), int64(persons.Data[0].Salary.Value))
	assert.Equal(t, int64(6000), int64(persons.Data[1].Salary.Value))
	assert.Equal(t, "USD", string(persons.Data[0].Salary.Currency))
	assert.Equal(t, "EUR", string(persons.Data[1].Salary.Currency))
}

func TestSortSalaryAsc(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"id": "1",
				"personName": "John Doe",
				"salary": {
					"value": 5000,
					"currency": "USD"
				}
			},
			{
				"id": "2",
				"personName": "Jane Doe",
				"salary": {
					"value": 6000,
					"currency": "EUR"
				}
			}
		]
	}`

	persons, err := UnmarshalPersons([]byte(jsonData))
	assert.Nil(t, err)

	sortedPersons := persons.SortSalary(true)
	assert.Equal(t, int64(5000), int64(sortedPersons[0].Salary.Value))
	assert.Equal(t, int64(6000), int64(sortedPersons[1].Salary.Value))
}

func TestSortSalaryDesc(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"id": "1",
				"personName": "John Doe",
				"salary": {
					"value": 5000,
					"currency": "USD"
				}
			},
			{
				"id": "2",
				"personName": "Jane Doe",
				"salary": {
					"value": 6000,
					"currency": "EUR"
				}
			}
		]
	}`

	persons, err := UnmarshalPersons([]byte(jsonData))
	assert.Nil(t, err)

	sortedPersons := persons.SortSalary(false)
	assert.Equal(t, int64(6000), int64(sortedPersons[0].Salary.Value))
	assert.Equal(t, int64(5000), int64(sortedPersons[1].Salary.Value))
}

func TestFilterSalaryInUSD(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"id": "1",
				"personName": "John Doe",
				"salary": {
					"value": 5000,
					"currency": "USD"
				}
			},
			{
				"id": "2",
				"personName": "Jane Doe",
				"salary": {
					"value": 6000,
					"currency": "EUR"
				}
			}
		]
	}`

	persons, err := UnmarshalPersons([]byte(jsonData))
	assert.Nil(t, err)

	filteredPersons := persons.FilterSalaryInUSD()
	assert.Len(t, filteredPersons, 1)
	assert.Equal(t, "John Doe", filteredPersons[0].PersonName)
}

func TestGroupCurrency(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"id": "1",
				"personName": "John Doe",
				"salary": {
					"value": 5000,
					"currency": "USD"
				}
			},
			{
				"id": "2",
				"personName": "Jane Doe",
				"salary": {
					"value": 6000,
					"currency": "EUR"
				}
			}
		]
	}`

	persons, err := UnmarshalPersons([]byte(jsonData))
	assert.Nil(t, err)

	groupedPersons := persons.GroupCurrency()
	assert.Len(t, groupedPersons, 2)
	assert.Len(t, groupedPersons[USD], 1)
	assert.Len(t, groupedPersons[EUR], 1)
	assert.Equal(t, "John Doe", groupedPersons[USD][0].PersonName)
	assert.Equal(t, "Jane Doe", groupedPersons[EUR][0].PersonName)
}
