package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
)

func main() {
    ch1, err := read("file1.csv")
    if err != nil {
        panic(fmt.Errorf("Could not read file1 %v", err))
    }
    br1 := breakup("1", ch1)
    br2 := breakup("2", ch1)
    br3 := breakup("3", ch1)

    for {
        fmt.Println("Hello")
        if br1 == nil && br2 == nil && br3 == nil {
            break
        }

        select {
        case v, ok := <-br1:
            fmt.Println("channel1", v)
            if !ok {
                br1 = nil
            }
        case v, ok := <-br2:
            fmt.Println("channel2", v)
            if !ok {
                br2 = nil
            }
        case v, ok := <-br3:
            fmt.Println("channel3", v)

            if !ok {
                br3 = nil
            }
        }
    }

    fmt.Println("All completed, exiting")
}

func breakup(worker string, ch <-chan []string) chan []string {
    chE := make(chan []string)

    go func() {
        for v := range ch {
            chE <- v
        }
        close(chE)
    }()
    return chE
}
func read(file string) (<-chan []string, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, fmt.Errorf("opening file %v", err)
    }

    ch := make(chan []string)

    cr := csv.NewReader(f)

    go func() {
        for {
            record, err := cr.Read()
            if err == io.EOF {
                close(ch)

                return
            }

            ch <- record
        }
    }()

    return ch, nil
}
