package mcsv

import "errors"

type CsvColumnReader struct {
	header       []string
	headerLen    int
	rowValues    []string
	rowValuesLen int

	// 多语言列, 预定义语言列表
	langArr []string
}

func NewCsvColumnValueReader(header, langArr []string) *CsvColumnReader {
	return &CsvColumnReader{header: header, headerLen: len(header), langArr: langArr}
}

func (cr *CsvColumnReader) SetRawValue(rowValue []string) error {
	cr.rowValues = rowValue
	cr.rowValuesLen = len(rowValue)
	if cr.rowValuesLen != cr.headerLen {
		return errors.New("got wrong row value number")
	}
	return nil
}

func (cr *CsvColumnReader) GetColumnValue(columnName string) string {
	for idx, name := range cr.header {
		if name == columnName && idx < cr.rowValuesLen {
			return cr.rowValues[idx]
		}
	}

	return ""
}

func (cr *CsvColumnReader) GetColumnMultiLanguageValue(columnName string) map[string]string {

	ret := make(map[string]string)
	for _, l := range cr.langArr {
		for idx, name := range cr.header {
			if name == columnName+"@"+l && idx < cr.rowValuesLen {
				ret[l] = cr.rowValues[idx]
			}
		}
	}

	return ret
}
