package model

import "github.com/HuckOps/forge/server/repository/generic"

type Language string

const (
	Golang     Language = "golang"
	Python     Language = "python"
	Javascript Language = "javascript"
	Bin        Language = "bin"
)

type Exporter struct {
	BaseModel
	Name     string                 `json:"name" bson:"name"`
	Args     map[string]interface{} `json:"args" bson:"args"`
	Language Language               `json:"language" bson:"language"`

	// 插件下载连接
	// Python： wheel文件
	// Go, Bin： 二进制文件
	// JS：git仓库
	SourceURL string `json:"source_url" bson:"source_url"`
}

func (model *Exporter) Repository() *generic.Repository[Exporter] {
	return generic.NewRepository[Exporter]("exporter")
}

//func (model *Exporter) Indexes() {}
