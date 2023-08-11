package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func main() {
	const MysqlConfig = "root:my-secret-pw@tcp(47.99.140.12:3306)/go_zero_douyin?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := gorm.Open(mysql.Open(MysqlConfig))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}
	// 自定义模型结体字段的标签
	// 将特定字段名的 json 标签加上`string`属性,即 MarshalJSON 时该字段由数字类型转成字符串类型
	// jsonField := gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
	// 	return strings.TrimPrefix(columnName, "u_")
	// })
	// 将非默认字段名的字段定义为自动时间戳和软删除字段;
	// 自动时间戳默认字段名为:`updated_at`、`created_at, 表字段数据类型为: INT 或 DATETIME
	// 软删除默认字段名为:`deleted_at`, 表字段数据类型为: DATETIME
	autoUpdateTimeField := gen.FieldGORMTag("update_time", func(tag field.GormTag) field.GormTag {
		return map[string][]string{"column": []string{"update_time"}, "type": []string{"int", "unsigned"}, "autoUpdateTime": []string{}}
	})
	autoCreateTimeField := gen.FieldGORMTag("create_time", func(tag field.GormTag) field.GormTag {
		return map[string][]string{"column": []string{"create_time"}, "type": []string{"int", "unsigned"}, "autoCreateTime": []string{}}
	})
	softDeleteField := gen.FieldType("delete_time", "gorm.DeletedAt")
	// fieldTrimPrefix := gen.FieldTrimPrefix("u_")
	// 模型自定义选项组
	// fieldOpts := []gen.ModelOpt{jsonField, autoCreateTimeField, autoUpdateTimeField, softDeleteField, fieldTrimPrefix}
	fieldOpts := []gen.ModelOpt{autoCreateTimeField, autoUpdateTimeField, softDeleteField}
	// 自定义字段的数据类型
	// 统一数字类型为int64,兼容protobuf
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint":   func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"smallint":  func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"mediumint": func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"bigint":    func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"int":       func(detailType gorm.ColumnType) (dataType string) { return "int64" },
	}
	// 生成实例
	conf := &gen.Config{
		// 生成query全局查询对象时候，需要区分Model path 和输出path Q会默认使用OutPath最后一个路径当包名称，如果不修改modelPKg会导致生成文件中导入结构体错误
		OutPath: "../cmd/rpc/internal/dao",
		// 生成全局查询文件名称
		OutFile:      "query.go",
		ModelPkgPath: "model",
		// WithDefaultQuery 生成默认查询结构体(作为全局变量使用), 即`Q`结构体和其字段(各表模型)
		// WithoutContext 生成没有context调用限制的代码供查询
		// WithQueryInterface 生成interface形式的查询代码(可导出), 如`Where()`方法返回的就是一个可导出的接口类型
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,

		// 表字段可为 null 值时, 对应结体字段使用指针类型
		FieldNullable: true, // generate pointer when field is nullable

		// 表字段默认值与模型结构体字段零值不一致的字段, 在插入数据时需要赋值该字段值为零值的, 结构体字段须是指针类型才能成功, 即`FieldCoverable:true`配置下生成的结构体字段.
		// 因为在插入时遇到字段为零值的会被GORM赋予默认值. 如字段`age`表默认值为10, 即使你显式设置为0最后也会被GORM设为10提交.
		// 如果该字段没有上面提到的插入时赋零值的特殊需要, 则字段为非指针类型使用起来会比较方便.
		FieldCoverable: true, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values

		// 模型结构体字段的数字类型的符号表示是否与表字段的一致, `false`指示都用有符号类型
		FieldSignable: false, // detect integer field's unsigned type, adjust generated data type
		// 生成 gorm 标签的字段索引属性
		FieldWithIndexTag: false, // generate with gorm index tag
		// 生成 gorm 标签的字段类型属性
		FieldWithTypeTag: false, // generate with gorm column type tag:
		// 生成单元测试
		WithUnitTest: true,
	}
	conf.WithDataTypeMap(dataMap)
	conf.WithOpts(fieldOpts...)
	g := gen.NewGenerator(*conf)
	// 设置目标 db
	g.UseDB(db)

	// 创建模型的结构体
	User := g.GenerateModelAs("user", "User")

	// Generate Type Safe API with Dynamic SQL defined on Queries interface for `model.User` and `model.Company`
	// g.ApplyInterface(func() {}, User, Major, ThirdData, UserAcc)
	g.ApplyBasic(User)

	g.Execute()

}
