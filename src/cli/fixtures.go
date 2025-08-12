package cli

var rootFileContent = `
[[@root]]

[[@tags]]
ref$my-tags

[[@groups]]
ref$my-groups
`

var nonRootFileContent = `
[[my-groups@ref]]
foo-group#fooTag: My foo group
`

var nestedFileContent = `[[my-tags@ref]]
foo-tag#262626: My foo tag
`

var duplicatedRefDeclaration = `
[[my-tags@ref]]
bar-tag#121212: Dup ref declaration tag
`
