package linter

import (
	"html/template"
)

// From Linter
type ManifestContent string
type LintStatus string

const (
	LintStatusOK    LintStatus = "ok"
	LintStatusError LintStatus = "error"
	LintStatusWarn  LintStatus = "warn"
)

type LintResult struct {
	Status    LintStatus
	Msg       string
	Structure *StructureGraph
}

type NodeInfo interface {
	isNodeInfo()
}

type Node struct {
	Info  NodeInfo
	Links []*Node
}

type StructureGraph struct {
	Root *Node
}

type RawBlock struct {
	Content   string
	StartLine int
}

// From Structure Graph
type Identifiable struct {
	Id string
}

func (i Identifiable) GetId() string {
	return i.Id
}

type LinkFields struct {
	_tagIds   []string
	_groupIds []string
}

func (l LinkFields) GetLinkIds(key string) []string {
	switch key {
	case "_tagIds":
		return l._tagIds
	case "_groupIds":
		return l._groupIds
	default:
		return nil
	}
}

type Linkable interface {
	GetId() string
	GetLinkIds(key string) []string
}

func (Root) isNodeInfo() {}

type Root struct {
	Identifiable
	ThemeStyle template.CSS
	BlockType  string
	LogoUrl    string
	GithubUrl  string
}

func (Content) isNodeInfo() {}

type Content struct {
	Identifiable
	BlockType string
	Summary   template.HTML
	LinkFields
}

func (Group) isNodeInfo() {}

type Group struct {
	Identifiable
	BlockType   string
	Description string
	LinkFields
}

func (Tag) isNodeInfo() {}

type Tag struct {
	Identifiable
	BlockType   string
	color       string
	Description string
}

func (Module) isNodeInfo() {}

type Module struct {
	Identifiable
	BlockType string
	Summary   template.HTML
	LinkFields
}
