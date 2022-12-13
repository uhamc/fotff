package rec

import (
	"fmt"
	"fotff/tester"
	"fotff/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/sirupsen/logrus"
	"reflect"
	"sort"
)

const css = `<head>
<style type="text/css">
table{
    border:1px solid;
    border-spacing: 0;
}
th{
    font-size: 11;
	border:1px solid;
    padding: 10px;
	background-color: rgb(137,190,178);
}
td{
    font-size: 11;
    border:1px solid;
    padding: 10px;
	background-color: rgb(160,191,124);
}
.bg-red{
	background-color: rgb(220,87,18);
}
.bg-yellow{
	background-color: rgb(244,208,0);
}
</style>
</head>
`

func Report(curPkg string, taskName string) {
	subject := fmt.Sprintf("[%s] %s test report", curPkg, taskName)
	rt := reflect.TypeOf(Record{})
	tb := table.NewWriter()
	tb.SetIndexColumn(rt.NumField() + 1)
	var row = table.Row{"test case"}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.IsExported() {
			row = append(row, f.Tag.Get("col"))
		}
	}
	tb.AppendHeader(row)
	tb.SetRowPainter(func(row table.Row) text.Colors {
		for _, col := range row {
			if str, ok := col.(string); ok {
				if str == tester.ResultFail {
					return text.Colors{text.BgRed}
				} else if str == tester.ResultOccasionalFail {
					return text.Colors{text.BgYellow}
				}
			}
		}
		return nil
	})
	var rows []table.Row
	for k, rec := range Records {
		var row = table.Row{k}
		rv := reflect.ValueOf(rec)
		for i := 0; i < rv.NumField(); i++ {
			if rv.Field(i).CanInterface() {
				row = append(row, rv.Field(i).Interface())
			}
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0].(string) < rows[j][0].(string)
	})
	tb.AppendRows(rows)
	if err := utils.SendMail(subject, css+tb.RenderHTML()); err != nil {
		logrus.Errorf("failed to send report mail: %v", err)
		return
	}
	logrus.Infof("send mail successfully")
}
