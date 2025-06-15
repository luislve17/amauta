package linter

var tagSectionRegex string = `^\[\[tags\]\]`
var tagRegex string = `^([-_\w]+)(#[A-F|\d]{6}):\s*(.*)`
var moduleSectionRegex string = `^\[\[([A-Z]+[\w| |-|_]*)\]\]`
