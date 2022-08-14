package generator

import (
	"fmt"
	"github.com/derekAHua/goLib/utils"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"io"
	"path"
	"path/filepath"
	"strings"
)

// @Author: Derek
// @Description:
// @Date: 2022/8/14 00:37
// @Version 1.0

func GenerateModel(directory, name string) {
	upperName := strings.ToUpper(name[:1]) + name[1:]
	fileName := name + ".go"
	allFileName := filepath.Join(directory, fileName)

	content := `
		package packageName

		type (
			### interface {}
		
			!!! struct {
				ctx *gin.Context
				model.BaseModel
			}
		)
		
		func New###(ctx *gin.Context) ### {
			return &!!!{ctx, model.NewBaseModel(conf.MysqlClientTest.WithContext(ctx).table("tableName"))}
		}
		
		func New###WithTx(ctx *gin.Context, tx *gorm.DB) ### {
			return &!!!{ctx, model.NewBaseModel(tx.table("tableName"))}
		}
	`
	content = strings.Replace(content, "packageName", path.Dir(allFileName), -1)
	content = strings.Replace(content, "###", upperName, -1)
	content = strings.Replace(content, "!!!", name, -1)
	content = strings.Replace(content, "tableName", "tbl"+generator.CamelCase(name), -1)

	f, err := utils.WriteFile(allFileName)
	if err != nil {
		panic(err)
	}

	defer func() { _ = f.Close() }()
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	}

	err = utils.RunFmt(directory, fileName)
	if err != nil {
		panic(err)
	}

	fmt.Println(" 已生成...")
}
