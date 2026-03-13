package utils

import (
	"fmt"
	"html/template"
	"strings"
)

// TemplateFuncs 返回模板函数映射
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"url_for": func(name string, args ...interface{}) string {
			switch name {
			case "static":
				if len(args) > 0 {
					// 处理 Flask 风格的关键字参数
					if len(args) == 2 {
						if filename, ok := args[1].(string); ok {
							return "/static/" + filename
						}
					}
					// 处理直接传递路径的情况
					return "/static/" + fmt.Sprint(args[0])
				}
				return "/static/"
			case "login":
				return "/login"
			case "logout":
				return "/logout"
			case "resources":
				return "/resources"
			case "users":
				return "/users"
			case "models":
				return "/models"
			default:
				return "/"
			}
		},
		// 添加其他 Flask 模板函数的等价实现
		"tojson": func(v interface{}) template.JS {
			return template.JS(fmt.Sprintf("%v", v))
		},
		"join": strings.Join,
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		// 添加一个辅助函数来处理静态文件
		"static": func(path string) string {
			return "/static/" + path
		},
	}
}
