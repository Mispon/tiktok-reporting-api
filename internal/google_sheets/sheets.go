package gsh

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	credentialsFile = "bin/credentials.json"
	sheetId         = "1033_iix0fnZDvh35GgcIXtZ9HwXRkIZPrc-gmnTbhzA"
)

type GoogleSheet interface {
	Init(ctx context.Context) error
	WriteRow(valuesRange string, values []interface{}) error
}

// New constructor
func New() GoogleSheet {
	return &googleSheet{}
}

type googleSheet struct {
	service *sheets.Service
}

// Init creates sheets service
func (gs *googleSheet) Init(ctx context.Context) error {
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return err
	}

	gs.service = srv
	return nil
}

// WriteRow writes simple row
func (gs *googleSheet) WriteRow(valuesRange string, values []interface{}) error {
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, values)

	res, err := gs.service.Spreadsheets.Values.Append(sheetId, valuesRange, &vr).ValueInputOption("RAW").Do()
	fmt.Println("spreadsheet push ", res)

	return err
}
