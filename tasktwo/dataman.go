package tasktwo

import (
	"encoding/json"
	"github.com/khrees2412/cadana/util"
	"io/ioutil"
	"log"
	"os"
)

func Start() {
	jsonFile, err := os.Open("tasktwo/persons.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()

	byteValue, e := ioutil.ReadAll(jsonFile)
	if e != nil {
		log.Println(err)
	}
	persons, err := UnmarshalPersons(byteValue)
	if err != nil {
		log.Println(err)
	}
	salaryAsc := persons.SortSalary(false)
	util.PrettyPrint(salaryAsc)
}

func UnmarshalPersons(data []byte) (Persons, error) {
	var r Persons
	err := json.Unmarshal(data, &r)
	return r, err
}

type Persons struct {
	Data []Person `json:"data"`
}

type IPerson interface {
	SortSalaryAsc() []Person
	SortSalaryDesc() []Person
	FilterSalaryInUSD() []Person
	GroupCurrency() map[Currency][]Person
}

type Person struct {
	ID         string `json:"id"`
	PersonName string `json:"personName"`
	Salary     Salary `json:"salary"`
}

type Salary struct {
	Value    int      `json:"value"`
	Currency Currency `json:"currency"`
}

type Currency string

const (
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	USD Currency = "USD"
)

func (r *Persons) SortSalary(asc bool) []Person {
	data := r.Data
	if !asc {
		for i := 0; i < len(data); i++ {
			for j := i + 1; j < len(data); j++ {
				if data[i].Salary.Value < data[j].Salary.Value {
					data[i], data[j] = data[j], data[i]
				}
			}
		}
	} else {
		for i := 0; i < len(data); i++ {
			for j := i + 1; j < len(data); j++ {
				if data[i].Salary.Value > data[j].Salary.Value {
					data[i], data[j] = data[j], data[i]
				}
			}
		}
	}

	return data
}

func (r *Persons) FilterSalaryInUSD() []Person {
	var person []Person
	for _, v := range r.Data {
		if v.Salary.Currency == "USD" {
			person = append(person, v)
		}
	}
	return person
}

func (r *Persons) GroupCurrency() map[Currency][]Person {
	m := make(map[Currency][]Person)
	for _, v := range r.Data {
		m[v.Salary.Currency] = append(m[v.Salary.Currency], v)
	}
	return m
}

// Extra methods

func (r *Persons) GetNames() []string {
	var names []string
	for _, v := range r.Data {
		names = append(names, v.PersonName)
	}
	return names
}
func (r *Persons) GetIDs() []string {
	var ids []string
	for _, v := range r.Data {
		ids = append(ids, v.ID)
	}
	return ids
}
func (r *Persons) GetCurrencies() []Currency {
	var currencies []Currency
	for _, v := range r.Data {
		currencies = append(currencies, v.Salary.Currency)
	}
	return currencies
}
func (r *Persons) GetSalariesValue() []int {
	var salaries []int
	for _, v := range r.Data {
		salaries = append(salaries, v.Salary.Value)
	}
	return salaries
}
