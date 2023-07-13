[![codecov](https://codecov.io/github/roccolangeweg/rdw-opendata-go/branch/main/graph/badge.svg?token=0G7BZZV6LY)](https://codecov.io/github/roccolangeweg/rdw-opendata-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/roccolangeweg/rdw-opendata-go)](https://goreportcard.com/report/github.com/roccolangeweg/rdw-opendata-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/roccolangeweg/rdw-opendata-go.svg)](https://pkg.go.dev/github.com/roccolangeweg/rdw-opendata-go)
![License](https://img.shields.io/github/license/roccolangeweg/rdw-opendata-go)
# rdw-opendata-go

This package provides a Go client for the RDW Open Data API available at https://opendata.rdw.nl. The client currently supports the following datasets:

- [x] [Registered Vehicles](https://opendata.rdw.nl/Voertuigen/Open-Data-RDW-Gekentekende_voertuigen/m9d7-ebf2)
- [ ] [Recognized Companies](https://opendata.rdw.nl/Voertuigen/Open-Data-RDW-erkende-bedrijven/8ys7-d773)
- [ ] [Recognized Driving Schools](https://opendata.rdw.nl/Voertuigen/Open-Data-RDW-erkende-rijscholen/534e-5vdg)

...and many more

## Installation

```bash
go get github.com/roccolangeweg/rdw-opendata-go
```

To use the client, you need to [register for an App token](https://opendata.rdw.nl/profile/edit/developer_settings). This can be requested free of charge.



## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    rdwopendatago "github.com/roccolangeweg/rdw-opendata-go"
)

func main() {
    client := rdwopendatago.NewClient(nil, "YOUR_APP_TOKEN")
    
    vehicles, err := client.RegisteredVehicles.List(context.Background(), rdwopendatago.RegisteredVehiclesListOptions{})
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(vehicles)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

* [RDW Open Data](https://opendata.rdw.nl)