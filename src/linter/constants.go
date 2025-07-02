package linter

var tagSectionRegex string = `^\[\[@tags\]\]`
var tagRegex string = `^([-_\w]+)(#[A-F|\d]{6}):\s*(.*)`
var moduleSectionHeaderRegex string = `^\[\[([A-Z]+[\w| |-|_]*)@api#?([\w|,|-]*)\]\]`
var contentSectionHeaderRegex string = `^\[\[([A-Z]+[\w| |-|_]*)@content#?([\w|,|-]*)\]\]`
