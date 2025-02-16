package analyzer

import (
	"errors"
	"fmt"
	"path/filepath"
	"walrus/compiler/internal/parser"
	"walrus/compiler/internal/typechecker"
	"walrus/compiler/report"
	"walrus/compiler/wio"
)

const HALTED = "compilation halted"

func Analyze(filePath string, displayErrors, debug, save2Json bool) (reports report.Reports, e error) {

	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("%v", r)
			reports = report.GetReports()
		}
		report.ClearReports()
	}()

	//must have .wal file
	if len(filePath) < 5 || filePath[len(filePath)-4:] != ".wal" {
		e = errors.New("error: file must have .wal extension")
		return nil, e
	}

	//get the folder and file name
	folder, fileName := filepath.Split(filePath)

	tree, e := parser.NewParser(filePath, debug).Parse()
	if e != nil {
		return report.GetReports(), e
	}

	if save2Json {
		//write the tree to a file named 'expressions.json' in 'code/ast' folder
		e = wio.Serialize(&tree, folder, fileName)
		if reports != nil {
			err := errors.New(report.TreeFormatString(HALTED, "Error serializing AST", e.Error()))
			return report.GetReports(), err
		}
	}

	typechecker.Analyze(tree, filePath)

	reports = report.GetReports()

	return reports, nil
}
