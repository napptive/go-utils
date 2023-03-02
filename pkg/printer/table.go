/**
 * Copyright 2020 Napptive
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package printer

import (
	"os"
	"reflect"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/napptive/nerrors/pkg/nerrors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// MinWidth is the minimal cell width including any padding.
	MinWidth = 8
	// TabWidth is the width of tab characters (equivalent number of spaces)
	TabWidth = 4
	// Padding added to a cell before computing its width
	Padding = 4
	// PaddingChar with the ASCII char used for padding
	PaddingChar = ' '
	// TabWriterFlags with the formatting options.
	TabWriterFlags = 0
)

// TablePrinter structure with the implementation required to print in a human-readable table format a given result.
type TablePrinter struct {
	extraTemplateFunctions map[string]any
	templates              map[reflect.Type]string
}

// NewTablePrinter builds a new ResultPrinter whose output is a human-readable table-like representation of the object.
func NewTablePrinter() (ResultPrinter, error) {
	extraFunctions := make(map[string]any, 0)
	printer := &TablePrinter{
		templates: make(map[reflect.Type]string, 0),
	}
	extraFunctions["fromTimestamp"] = printer.fromTimestamp
	extraFunctions["fromTimestampUint"] = printer.fromTimestampUint
	extraFunctions["fromProtoTimestampToDate"] = printer.fromProtoTimestampToDate
	extraFunctions["capitalize"] = printer.CapitalizeWord
	extraFunctions["fromProtoTimestampToUTCTime"] = printer.fromProtoTimestampToUTCTime
	printer.extraTemplateFunctions = extraFunctions
	return printer, nil
}

func (tp *TablePrinter) AddTemplate(templateType reflect.Type, templateContent string) {
	tp.templates[templateType] = templateContent
}

// GetTemplate returns a template to print an arbitrary structure in table format.
func (tp *TablePrinter) GetTemplate(result interface{}) (*string, error) {
	template, exists := tp.templates[reflect.TypeOf(result)]
	if !exists {
		return nil, nerrors.NewUnimplementedError("no template is available to print %s", reflect.TypeOf(result).String())
	}
	return &template, nil
}

// fromTimestamp returns the string representation of a timestamp.
func (tp *TablePrinter) fromTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).String()
}

// fromTimestamp returns the string representation of a timestamp.
func (tp *TablePrinter) fromTimestampUint(timestamp uint64) string {
	return tp.fromTimestamp(int64(timestamp))
}

// fromProtoTimestampToDate transforms a proto timestamp into a date
func (tp *TablePrinter) fromProtoTimestampToDate(timestamp *timestamppb.Timestamp) string {
	return timestamp.AsTime().Format("2006-01-02")
}

// fromProtoTimestampToUTCTime transforms a proto timestamp into a timestamp
func (tp *TablePrinter) fromProtoTimestampToUTCTime(timestamp *timestamppb.Timestamp) string {
	return timestamp.AsTime().UTC().String()
}

// CapitalizeWord sets as upper case the first letters of a word
func (tp *TablePrinter) CapitalizeWord(word string) string {
	return cases.Title(language.Und).String(word)
}

// AddTemplateFunction adds a new function so that it can be called in the golang template.
func (tp *TablePrinter) AddTemplateFunction(name string, f interface{}) {
	tp.extraTemplateFunctions[name] = f
}

// Print the result.
func (tp *TablePrinter) Print(result interface{}) error {
	associatedTemplate, err := tp.GetTemplate(result)
	if err != nil {
		return err
	}
	t := template.New("TablePrinter").Funcs(
		tp.extraTemplateFunctions,
	)
	t, err = t.Parse(*associatedTemplate)
	if err != nil {
		return nerrors.NewInternalErrorFrom(err, "cannot apply template")
	}
	w := tabwriter.NewWriter(os.Stdout, MinWidth, TabWidth, Padding, PaddingChar, TabWriterFlags)
	if err := t.Execute(w, result); err != nil {
		return err
	}
	w.Flush()
	return nil
}

// PrintResultOrError prints the result using a given printer or the error.
func (tp *TablePrinter) PrintResultOrError(result interface{}, err error) error {
	return PrintResultOrError(tp, result, err)
}
