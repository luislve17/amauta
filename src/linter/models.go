package linter

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

type Node struct {
	Info  map[string]interface{}
	Links []*Node
}

type StructureGraph struct {
	Root *Node
}

type RawBlock struct {
	Content   string
	StartLine int
}
