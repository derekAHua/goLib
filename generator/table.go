package generator

import (
	"fmt"
	"github.com/derekAHua/goLib/utils"
	"io"
	"log"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/golang/protobuf/protoc-gen-go/generator"
	"gorm.io/gorm"
)

type (
	table struct {
		Name    string `gorm:"column:Name"`
		Comment string `gorm:"column:Comment"`
	}

	field struct {
		Field      string `gorm:"column:field"`
		Type       string `gorm:"column:Type"`
		Null       string `gorm:"column:Null"`
		Key        string `gorm:"column:Key"`
		Default    string `gorm:"column:Default"`
		Extra      string `gorm:"column:Extra"`
		Privileges string `gorm:"column:Privileges"`
		Comment    string `gorm:"column:Comment"`
	}
)

func GenerateTable(directory string, client *gorm.DB, dbName string, tableNames ...string) {
	tableNamesStr := ""
	for _, name := range tableNames {
		if tableNamesStr != "" {
			tableNamesStr += ","
		}
		tableNamesStr += "'" + name + "'"
	}
	tables := getTables(client, dbName, tableNamesStr) // 生成所有表信息
	for index := range tables {
		fields := getFields(client, tables[index].Name)
		generateModel(directory, tables[index], fields)
	}
}

// 获取表信息
func getTables(client *gorm.DB, dbName string, tableNames string) (ret []table) {
	if tableNames != "" {
		client.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + dbName + "' and TABLE_NAME in (" + tableNames + ");").Find(&ret)
		return
	}

	client.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + dbName + "';").Find(&ret)
	return
}

// 获取所有字段信息
func getFields(client *gorm.DB, tableName string) (ret []field) {
	client.Raw("show FULL COLUMNS from " + tableName + ";").Find(&ret)
	return
}

// 生成Model
func generateModel(directory string, table table, fields []field) {
	daoFileName := strings.ToLower(strings.Split(table.Name, "tbl")[1])
	structName := generator.CamelCase(daoFileName)
	fileName := filepath.Join(directory, daoFileName+".go")

	content := "package " + path.Dir(fileName) + "\n\n"        // package
	content += "// " + structName + " " + table.Comment + "\n" // 表注释
	content += "type " + structName + " struct {\n"            // 生成struct
	// 生成字段
	content += "model.Model\n"
	for _, field := range fields {
		fieldName := generator.CamelCase(field.Field)
		if fieldName == "Id" || fieldName == "CreateTime" || fieldName == "LastModifyTime" || fieldName == "Deleted" {
			continue
		}
		fieldType := getFiledType(field)
		fieldGorm := getFieldGorm(field)
		fieldJson := getFieldJson(field)
		fieldComment := getFieldComment(field)
		content += "	" + fieldName + " " + fieldType + " `" + fieldGorm + " " + fieldJson + "` " + fieldComment + "\n"
	}
	content += "}\n"

	content += "func (" + "*" + structName + ") TableName() string {\n"
	content += "	" + `return "` + table.Name + `"`
	content += "\n}\n"

	f, err := utils.WriteFile(fileName)
	defer func() { _ = f.Close() }()

	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	}
	fmt.Println(generator.CamelCase(table.Name) + " 已生成...")

	err = utils.RunFmt(directory, daoFileName+".go")
	if err != nil {
		log.Fatal(err)
		return
	}
}

// 获取字段类型
func getFiledType(field field) string {
	var tyArr []string
	if strings.Contains(field.Type, " ") {
		tyArr = strings.Split(field.Type, " ")
	} else {
		tyArr = strings.Split(field.Type, "(")
	}

	switch tyArr[0] {
	case "int":
		if strings.Contains(field.Type, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	case "integer":
		if strings.Contains(field.Type, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	case "mediumint":
		if strings.Contains(field.Type, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	case "bit":
		if strings.Contains(field.Type, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	case "year":
		if strings.Contains(field.Type, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	case "smallint":
		if strings.Contains(field.Type, "unsigned") {
			return "uint16"
		} else {
			return "int16"
		}
	case "tinyint":
		if strings.Contains(field.Type, "unsigned") {
			return "uint8"
		} else {
			return "int8"
		}
	case "bigint":
		if strings.Contains(field.Type, "unsigned") {
			return "uint64"
		} else {
			return "int64"
		}
	case "decimal":
		return "float64"
	case "double", "float", "real", "numeric":
		return "float32"
	case "datetime":
		return "jsontime.JsonTime"
	case "timestamp", "time", "date":
		return "time.Time"
	default:
		return "string"
	}
}

// 获取字段json描述
func getFieldJson(field field) string {
	return `json:"` + firstLet(generator.CamelCase(field.Field)) + `"`
}

// 首字母小写
func firstLet(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// 获取字段gorm描述
func getFieldGorm(field field) string {
	fieldContext := `gorm:"column:` + field.Field
	if field.Field == "create_time" {
		fieldContext = fieldContext + `;autoCreateTime`
	} else if field.Field == "update_time" {
		fieldContext = fieldContext + `;autoUpdateTime`
	}
	if field.Key == "PRI" {
		fieldContext = fieldContext + `;primaryKey`
	}
	if field.Key == "UNI" {
		fieldContext = fieldContext + `;unique`
	}
	if field.Default != "" {
		fieldContext = fieldContext + `;default:` + field.Default
	}
	if field.Extra == "auto_increment" {
		fieldContext = fieldContext + `;autoIncrement`
	}
	if field.Null == "NO" {
		fieldContext = fieldContext + `;not null`
	}
	return fieldContext + `"`
}

// 获取字段说明
func getFieldComment(field field) string {
	if len(field.Comment) > 0 {
		return "// " + strings.Replace(strings.Replace(field.Comment, "\r", "\\r", -1), "\n", "\\n", -1)
	}
	return ""
}
