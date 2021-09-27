package service

import (
	"errors"
	"fmt"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
	"log"
	"strconv"
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

func SendContact(srv *sheets.Service, id string, contact Result) error {
	_,y, err := searchRows(srv, id, "")
	if err != nil {
		return err
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	lenSheet := strconv.Itoa(y[len(y)-1]+1)
	inValues := make([]interface{}, 5)
	inValues[0] = contact.Name
	inValues[1] = contact.Job
	inValues[2] = contact.Cell
	inValues[3] = contact.Tg
	inValues[4] = contact.Other
	outValue := make([][]interface{}, 1)
	outValue[0] = inValues
	valRen := sheets.ValueRange{
		MajorDimension:  "ROWS",
		Range:           "",
		Values:          outValue,
		ServerResponse:  googleapi.ServerResponse{},
		ForceSendFields: nil,
		NullFields:      nil,
	}
	rangeAE := "Sheet1!A"+lenSheet+":E"+lenSheet
	resp, err := srv.Spreadsheets.Values.Update(id, rangeAE, &valRen).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
		return err
	}
	fmt.Println(resp)
	return nil
}

func findContacts(srv *sheets.Service, id, name string) ([]Result, error) {
	resp, y, err := searchRows(srv, id, name)
	out := make([]Result, 0)
	if err != nil {
		return []Result{}, err
	}
	for i := range y {
		a := fmt.Sprintf("%v", y[i])
		e := fmt.Sprintf("%v", y[i])
		res := "Sheet1!A" + a + ":" + "E" + e
		rsp, err := srv.Spreadsheets.Values.Get(id, res).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet: %v", err)
		}
		if len(resp.Values) == 0 {
			fmt.Println("No data found.")
		} else {
			for _, row := range rsp.Values {
				var newrow [5]string
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

func searchRows(srv *sheets.Service, id, name string) (*sheets.ValueRange, []int, error) {
	readRange := "Sheet1!A:B"
	resp, err := srv.Spreadsheets.Values.Get(id, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	var y []int
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return &sheets.ValueRange{
			MajorDimension:  "",
			Range:           "",
			Values:          nil,
			ServerResponse:  googleapi.ServerResponse{},
			ForceSendFields: nil,
			NullFields:      nil,
		}, []int{}, errors.New("no rows")
	} else {
		for i, row := range resp.Values {
			for _, cell := range row {
				c1 := strings.ToLower(fmt.Sprintf("%v", cell))
				c2 := strings.ToLower(fmt.Sprintf("%v", cell))
				if strings.Contains(c1, strings.ToLower(name)) || strings.Contains(c2, strings.ToLower(name)) {
					y = append(y, i+1)
				}
			}
		}
	}
	return resp, y, nil
}
