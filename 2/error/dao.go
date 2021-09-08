package main

import (
    "database/sql"
    "fmt"
    "github.com/pkg/errors"
)

// Fake dao example
func QueryRow(input string) ([]string, error) {
    result, err := fakeDbQueryHandler(input)
    if err != nil {
        if err == sql.ErrNoRows {
            // No rows found for the given condition, handled differently
            fmt.Println("No result found for the input", input)
            return []string{}, nil
        } else {
            // Wrap and throw
            return nil, errors.Wrap(err, "Query failed")
        }
    }
    return result, nil
}

// Fake handler mock
func fakeDbQueryHandler(input string) ([]string, error) {
    switch input {
    case "NR":
        return nil, sql.ErrNoRows
    case "CD":
        return nil, sql.ErrConnDone
    case "TD":
        return nil, sql.ErrTxDone
    default: 
        return []string{"fake-result", input}, nil
    }
}

func checkPrint(input string) {
    fmt.Println("------------")
    fmt.Println("Input:", input)
    r, err := QueryRow(input)
    if err != nil {
        fmt.Printf("FATAL: %+v\n", err)
    } else {
        fmt.Println("Result:", r)
    }
}

func main() {
    checkPrint("NR")
    checkPrint("CD")
    checkPrint("TD")
    checkPrint("AAA")
}
