## Task 1

### Description
This is a simple API to fetch exchange rate between two currencies.
To start the server, run the following command. A server will start at port 5001
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