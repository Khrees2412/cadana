## Task 1

### Description
This is a simple API to fetch exchange rate between two currencies.
To start the server, run the following command. A server will start at port 5001

Set main.go
```go
func main() {
	taskone.Start()
}
```
Then
```go
    go run main.go
```
### Run the test
```go
    go test
```

### API Endpoints
- Fetch exchange rate
    - GET /exchange-rate?currency_pair=USD-EUR



*Response*
```json
    {
        "USD-EUR": 0.85
    }
```

## Additional Task

Set main.go
```go
func main() {
	tasktwo.Start()
}
```

Modify dataman.go
```go
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
	// to sort salary in ascending order
	salaryAsc := persons.SortSalaryAsc()
	// to sort salary in descending order
    salaryDesc := persons.SortSalaryDesc()
    util.PrettyPrint(salaryAsc)
    util.PrettyPrint(salaryDesc)
}
```

Then
```go
    go run main.go
```