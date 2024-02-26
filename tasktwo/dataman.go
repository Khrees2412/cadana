package tasktwo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func UnmarshalPersons(data []byte) (Persons, error) {
	var r Persons
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Persons) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Persons struct {
	Data []Person `json:"data"`
}

type Person struct {
	ID         string `json:"id"`
	PersonName string `json:"personName"`
	Salary     Salary `json:"salary"`
}

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

type IPerson interface {
	SortSalaryAsc() []Person
	SortSalaryDesc() []Person
	FilterSalaryInUSD() []Person
	GroupCurrency() map[Currency][]Person
}

func Start() {
	jsonFile, err := os.Open("persons.json")
	if err != nil {
		fmt.Println(err)
		//return nil
	}
	defer jsonFile.Close()

	byteValue, e := ioutil.ReadAll(jsonFile)
	if e != nil {
		fmt.Println(e)
		//return nil
	}
	persons, err := UnmarshalPersons(byteValue)
	if err != nil {
		//return nil
	}
	salaryAsc := persons.SortSalaryAsc()
	PrettyPrint(salaryAsc)
	//return &persons

}

func PrettyPrint(p interface{}) {
	pretty, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(pretty))

}

func (r *Persons) GetDataFromJSON() *Persons {
	jsonFile, err := os.Open("persons.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()

	byteValue, e := ioutil.ReadAll(jsonFile)
	if e != nil {
		fmt.Println(e)
		return nil
	}
	persons, err := UnmarshalPersons(byteValue)
	if err != nil {
		return nil
	}
	salaryAsc := r.SortSalaryAsc()
	fmt.Println(salaryAsc)
	return &persons
}

func (r *Persons) SortSalaryAsc() []Person {
	data := r.Data
	for i := 0; i < len(data); i++ {
		for j := i + 1; j < len(data); j++ {
			if data[i].Salary.Value > data[j].Salary.Value {
				data[i], data[j] = data[j], data[i]
			}
		}
	}
	return data
}

func (r *Persons) SortSalaryDesc() []Person {
	data := r.Data
	for i := 0; i < len(data); i++ {
		for j := i + 1; j < len(data); j++ {
			if data[i].Salary.Value < data[j].Salary.Value {
				data[i], data[j] = data[j], data[i]
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
