package form

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-macaron/binding"
	"github.com/unknwon/com"
)

type Form interface {
	binding.Validator
}

func init() {
	binding.SetNameMapper(com.ToSnakeCase)
}

// Assign 将字段值返回表单
func Assign(form interface{}, data map[string]interface{}) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		} else if len(fieldName) == 0 {
			fieldName = com.ToSnakeCase(field.Name)
		}

		data[fieldName] = val.Field(i).Interface()
	}
}

func validate(errs binding.Errors, data map[string]interface{}, f Form) binding.Errors {
	if errs.Len() == 0 {
		return errs
	}

	data["HasError"] = true
	Assign(f, data)

	typ := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// 忽略的字段
		if fieldName == "-" {
			continue
		}

		// 报错信息字段名
		trName := field.Tag.Get("locale")
		if len(trName) == 0 {
			trName = fieldName
		}

		if errs[0].FieldNames[0] == field.Name {
			switch errs[0].Classification {
			case binding.ERR_REQUIRED:
				data["ErrorMsg"] = trName + "不能为空"
			case binding.ERR_ALPHA_DASH:
				data["ErrorMsg"] = trName + "必须为英文字母、阿拉伯数字或横线（-_）"
			case binding.ERR_ALPHA_DASH_DOT:
				data["ErrorMsg"] = trName + "必须为英文字母、阿拉伯数字、横线（-_）或点"
			case binding.ERR_SIZE:
				data["ErrorMsg"] = trName + fmt.Sprintf("长度必须为 %s", getSize(field))
			case binding.ERR_MIN_SIZE:
				data["ErrorMsg"] = trName + fmt.Sprintf("长度最小为 %s 个字符", getMinSize(field))
			case binding.ERR_MAX_SIZE:
				data["ErrorMsg"] = trName + fmt.Sprintf("长度最大为 %s 个字符", getMaxSize(field))
			case binding.ERR_EMAIL:
				data["ErrorMsg"] = trName + "不是一个有效的邮箱地址"
			case binding.ERR_URL:
				data["ErrorMsg"] = trName + "不是一个有效的 URL"
			default:
				data["ErrorMsg"] = "未知错误" + " " + errs[0].Classification
			}
			return errs
		}
	}
	return errs
}

func getRuleBody(field reflect.StructField, prefix string) string {
	for _, rule := range strings.Split(field.Tag.Get("binding"), ";") {
		if strings.HasPrefix(rule, prefix) {
			return rule[len(prefix) : len(rule)-1]
		}
	}
	return ""
}

func getSize(field reflect.StructField) string {
	return getRuleBody(field, "Size(")
}

func getMinSize(field reflect.StructField) string {
	return getRuleBody(field, "MinSize(")
}

func getMaxSize(field reflect.StructField) string {
	return getRuleBody(field, "MaxSize(")
}
