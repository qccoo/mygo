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
            // 没有对应查询条件的行，相应处理逻辑, 不同于其他error
            fmt.Println("No result found for the input", input)
            return []string{}, nil
        } else {
            // 其他错误wrap抛给上层(或处理时可考虑log.Fatal)
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
        return
    }
    fmt.Println("Result:", r)
}

func main() {
    checkPrint("NR")
    checkPrint("CD")
    checkPrint("TD")
    checkPrint("AAA")
}
