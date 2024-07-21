package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const MysqlDSN = "root:@(127.0.0.1:3306)/janction?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../query",
		OutFile:           "",
		ModelPkgPath:      "../model",
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  false,
		Mode:              gen.WithoutContext | gen.WithDefaultQuery,
	})

	g.UseDB(connectDB(MysqlDSN))

	// generate all table from database
	g.WithModelNameStrategy(func(tableName string) (modelName string) {

		return convertTableToModelName(tableName)
	})
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint":   func(columnType gorm.ColumnType) (dataType string) { return "int" },
		"smallint":  func(columnType gorm.ColumnType) (dataType string) { return "int" },
		"mediumint": func(columnType gorm.ColumnType) (dataType string) { return "int" },
		"bigint":    func(columnType gorm.ColumnType) (dataType string) { return "int" },
		"int":       func(columnType gorm.ColumnType) (dataType string) { return "int" },
	}
	g.WithDataTypeMap(dataMap)
	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}

func connectDB(dsn string) (db *gorm.DB) {
	var err error
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}
	return db
}

func convertTableToModelName(tableName string) (new string) {
	var nameParts = strings.Split(strings.ReplaceAll(tableName, "tb_", ""), "_")
	return joinStringSlice(nameParts, "", cases.Title(language.English).String)
}

func joinStringSlice(elems []string, sep string, convert func(old string) (new string)) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return convert(elems[0])
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteString(convert(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(convert(s))
	}
	return b.String()
}
