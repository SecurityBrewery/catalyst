package main

type Swagger struct {
	Swagger     string                `yaml:"swagger" json:"swagger"`
	Info        Info                  `yaml:"info" json:"info"`
	Paths       map[string]PathItem   `yaml:"paths" json:"paths"`
	Definitions map[string]Definition `yaml:"definitions" json:"definitions"`
}

type Info struct {
	Version string `yaml:"version" json:"version"`
	Title   string `yaml:"title" json:"title"`
}

type PathItem struct {
	Get  Operation `yaml:"get" json:"get"`
	Post Operation `yaml:"post" json:"post"`
}

type Operation struct {
	Tags        []string            `yaml:"tags" json:"tags"`
	Summary     string              `yaml:"summary" json:"summary"`
	Description string              `yaml:"description" json:"description"`
	OperationID string              `yaml:"operationId" json:"operationId"`
	Parameters  []Parameter         `yaml:"parameters" json:"parameters"`
	Responses   map[string]Response `yaml:"responses" json:"responses"`
}

type Parameter struct {
	// Description string `yaml:"description" json:"description"`
	Name     string `yaml:"name" json:"name"`
	In       string `yaml:"in" json:"in"`
	Required bool   `yaml:"required" json:"required"`
	Schema   Schema `yaml:"schema" json:"schema"`
	Direct   Schema `yaml:"-,inline" json:"-,inline"`
}

type Schema struct {
	Ref                  string      `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Format               string      `yaml:"format,omitempty" json:"format,omitempty"`
	Title                string      `yaml:"title,omitempty" json:"title,omitempty"`
	Description          string      `yaml:"description" json:"description,omitempty"`
	Default              interface{} `yaml:"default,omitempty" json:"default,omitempty"`
	Maximum              interface{} `yaml:"maximum,omitempty" json:"maximum,omitempty"`
	Items                *Schema     `yaml:"items,omitempty" json:"items,omitempty"`
	Type                 string      `yaml:"type" json:"type,omitempty"`
	AdditionalProperties *Schema     `yaml:"additionalProperties,omitempty" json:"additionalProperties,omitempty"`
	Enum                 []string    `yaml:"enum,omitempty" json:"enum,omitempty"`
}

type Response struct {
	Description string `yaml:"description" json:"description"`
	Schema      Schema `yaml:"schema" json:"schema"`
}

type Definition struct {
	Type       string            `yaml:"type" json:"type"`
	Required   []string          `yaml:"required,omitempty" json:"required,omitempty"`
	Embed      string            `yaml:"x-embed" json:"x-embed"`
	Properties map[string]Schema `yaml:"properties" json:"properties"`
}
