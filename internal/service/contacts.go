package service

import (
	"fmt"
	"google.golang.org/api/sheets/v4"
	"log"
	"strings"
)

type Result struct {
	Name  string
	Job   string
	Cell  string
	Tg    string
	Other string
}

func GetContact(srv *sheets.Service, id, name string) ([]Result, error) {
	res, err := findContacts(srv, id, name)
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	return res, nil
}

func findContacts(srv *sheets.Service, id, name string) ([]Result, error) {
	readRange := "Sheet1!A:B"
	resp, err := srv.Spreadsheets.Values.Get(id, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	var y []int
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return []Result{}, nil
	} else {
		for i, row := range resp.Values {
			for _, cell := range row {
				c1 := fmt.Sprintf("%v", cell)
				c2 := fmt.Sprintf("%v", cell)
				if strings.Contains(c1, strings.ToLower(name))||strings.Contains(c2, strings.ToLower(name)) {
					y = append(y, i+1)
					println(y)
				}
			}
		}
		out := make([]Result, 0)

		for i := range y {
			a := fmt.Sprintf("%v", y[i])
			e := fmt.Sprintf("%v", y[i])
			res := "Sheet1!A" + a + ":" + "E" + e
			fmt.Println(res)
			rsp, err := srv.Spreadsheets.Values.Get(id, res).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve data from sheet: %v", err)
			}
			if len(resp.Values) == 0 {
				fmt.Println("No data found.")
			} else {
				for _, row := range rsp.Values {
					var newrow [5]string
					fmt.Println(len(row))
					for i := 0; i < len(row); i++ {
						newrow[i] = fmt.Sprintf("%v", row[i])
					}
					for i := 4; i < 5-len(row); i++ {
						newrow[i] = fmt.Sprintf("%v", "")
					}
					tmp := Result{
						Name:  fmt.Sprintf("%v", newrow[0]),
						Job:   fmt.Sprintf("%v", newrow[1]),
						Cell:  fmt.Sprintf("%v", newrow[2]),
						Tg:    fmt.Sprintf("%v", newrow[3]),
						Other: fmt.Sprintf("%v", newrow[4]),
					}
					out = append(out, tmp)
				}
			}
		}
		return out, nil
	}
	return []Result{}, err
}
