package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	. "catface/AutoMigrateMySQL/config"

	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 仿照 AutoMigrate 原本的效果
func convertToSnakeCase(name string) string {
	// 使用正则表达式找到大写字符并在前面加上下划线，然后转换为小写
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(name, "${1}_${2}")
	return strings.ToLower(snake)
}

// 从 AST 中提取结构体字段类型
func getFieldType(expr ast.Expr) reflect.Type {
	switch t := expr.(type) {
	case *ast.Ident:
		// fmt.Println("t.Name:", t.Name)
		switch t.Name {
		case "string":
			return reflect.TypeOf("")
		case "int":
			return reflect.TypeOf(0)
		case "bool":
			return reflect.TypeOf(true)
		case "uint8":
			return reflect.TypeOf(uint8(0))
		case "uint16":
			return reflect.TypeOf(uint16(0))
		case "uint32":
			return reflect.TypeOf(uint32(0))
		case "uint64":
			return reflect.TypeOf(uint64(0))
		case "float64":
			return reflect.TypeOf(float64(0))
		}
	case *ast.ArrayType:
		elemType := getFieldType(t.Elt)
		if elemType != nil {
			return reflect.SliceOf(elemType)
		}
	case *ast.SelectorExpr: // info time.Time 的特化识别
		if pkgIdent, ok := t.X.(*ast.Ident); ok {
			if pkgIdent.Name == "time" && t.Sel.Name == "Time" {
				return reflect.TypeOf(time.Time{})
			}
			if pkgIdent.Name == "datatypes" && t.Sel.Name == "JSON" {
				return reflect.TypeOf(datatypes.JSON{})
			}
		}
	case *ast.StarExpr: // question 暂时好像不影响。
		// Handle pointer to a type
		return reflect.PtrTo(getFieldType(t.X))
	}
	return nil
}

// 用于保存结构体信息的map
var structs = make(map[string]reflect.Type)

// 遍历文件中的所有声明
func getStruct(fDecls []ast.Decl) {
	for _, decl := range fDecls {
		// 检查声明是否为类型声明（type T struct {...})
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		// 遍历类型声明中的所有规格（可能有多个类型在一个声明中，例如：type (A struct{}; B struct{})）
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// 检查类型是否为结构体
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// 过滤空表
			if len(structType.Fields.List) == 0 {
				continue
			}

			// 构建反射类型
			fields := make([]reflect.StructField, 0)
			// fmt.Println(typeSpec.Name.Name, len(structType.Fields.List))

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					// 处理嵌入结构体
					ident, ok := field.Type.(*ast.Ident)
					if !ok {
						log.Printf("Unsupported embedded type for field %v\n", field.Type)
						continue
					}
					embedType, ok := structs[ident.Name]
					if !ok {
						log.Printf("Embedded type %s not found\n", ident.Name)
						continue
					}
					// 获取嵌入结构体的所有字段
					for i := 0; i < embedType.NumField(); i++ {
						fields = append(fields, embedType.Field(i))
					}
				} else {
					for _, fieldName := range field.Names {
						fieldType := getFieldType(field.Type)
						if fieldType == nil {
							continue
						}

						// 处理标签
						tag := ""
						if field.Tag != nil {
							tag = field.Tag.Value
						}

						fields = append(fields, reflect.StructField{
							Name: fieldName.Name,
							Type: fieldType,
							Tag:  reflect.StructTag(tag),
						})
						// fmt.Println(fieldName.Name, field.Type, fieldType, tag)
					}
				}
			}

			// 创建结构体类型
			structName := typeSpec.Name.Name
			structReflectType := reflect.StructOf(fields)
			structs[structName] = structReflectType
			fmt.Println("get struct: ", structName)
		}
	}
}

func autoMigrate() {
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalln("Error loading config: %v", err)
	}
	// info 初始化数据库
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MySQL.Username,
		config.MySQL.Password,
		config.MySQL.Host,
		config.MySQL.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // 打开 DB 连接

	if err != nil {
		log.Fatal(err)
	}

	// 通过反射创建结构体实例并迁移数据库
	for name, typ := range structs {
		if name == "General" {
			continue
		}
		instance := reflect.New(typ).Interface()

		// 手动设定表名
		tableName := convertToSnakeCase(name)
		db = db.Table(tableName)

		fmt.Printf("Created instance of %s: %+v\n", name, instance)
		if err := db.AutoMigrate(instance); err != nil {
			log.Fatalf("Failed to migrate %s: %v", name, err)
		}
	}
}

func main() {
	const dirPath = "./table_defs"       // 指定目录路径
	const rootFileName = "table_defs.go" // info 根结构体所在的文件
	fset := token.NewFileSet()           // 创建文件集，用于记录位置

	// mark stage-1
	// 列出指定目录下的所有文件和子目录
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// 前置根文件
	var targetIndex int
	var found bool
	for i, entry := range entries {
		if entry.Name() == rootFileName {
			targetIndex = i
			found = true
			break
		}
	}

	if found {
		targetEntry := entries[targetIndex]
		entries = append(entries[:targetIndex], entries[targetIndex+1:]...)
		entries = append([]os.DirEntry{targetEntry}, entries...)
	} else {
		log.Fatalf("File %s not found in directory %s", rootFileName, dirPath)
	}

	// 正常遍历
	for _, entry := range entries {
		// 构建完整路径
		path := filepath.Join(dirPath, entry.Name())

		// 检查文件后缀是否为 .go
		if !entry.IsDir() && filepath.Ext(path) == ".go" {
			// 解析文件，得到 *ast.File 结构
			f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				log.Printf("Error parsing file %s: %s", path, err)
				continue
			}

			// mark stage-2
			getStruct(f.Decls)
			log.Printf("Parsed file: %s", path)
		}
	}

	// mark stage-3
	autoMigrate()
}
