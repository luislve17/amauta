package linter

var ValidManifest string = `
[[@root]]
LogoUrl: https://raw.githubusercontent.com/luislve17/amauta/refs/heads/main/assets/amauta-banner.svg
GithubUrl: https://github.com/luislve17/amauta
<--
This is a multiline comment block
Starts with --> and ends with <--
but both must be on a new line with nothing
else but that characters sequence.
-->

<--
Basic types to identify are:
.int
.float
.str
.object
.list
-->

-- Inline comment block.

-- tags: Special section  
-- items follow: '<name>#<color>: <description>' syntax
[[@tags]]
public#00FF00: Public API  
internal#AAAAAA: Internal use only  
deprecated#FF6F61: Will be removed soon  
under-dev#FFD966: Still under development  
beta#87CEEB: Beta feature  
admin#FF1493: Admin only

[[@groups]]
getting-started#public: Getting started
finishing#public: Almost done
api#public: API


<--
module: Special section  
header follows: '[[<name>#<tag1,tag2,...>]]'  
fields follow: '<name>(@|?)<type>#<tag1,tag2,...>: <description>'
-->
[[About amauta@content#public]]
group: getting-started
summary: <md>
# Amauta
Welcome to Amauta. The right tool to write documentation. This line is excesively long just to verify how line breaks render once that they load on templates

## What is Amauta?
Amauta is a non-standardized documentation tool. Improves readability and maintanibility for small and big developer teams.

## List of things

* Item 1
* Item 2
* Item 3


### Sub topic

> "This is a quote" ~ The quoter

Here is a table:

| Item              | In Stock | Price |
| :---------------- | :------: | ----: |
| Python Hat        |   True   | 23.99 |
| SQL Hat           |   True   | 23.99 |
| Codecademy Tee    |  False   | 19.99 |
| Codecademy Hoodie |  False   | 42.99 |

</md>

[[Users@api#public,internal]]
group: api
summary: Endpoints related to user operations

-- section  
-- header follows: '[<section>@<type>:<description>]'  
-- fields follow: '<name>(@|?)<type>#<tag1,tag2,...>: <description>'
[request@POST:/v1/users]
summary: Create a new user  
contentType: application/json  

-- Inline object declaration  
-- Only allowed for 1-depth object
header@object: This is the object root description  
header.Authorization@str#internal: Bearer token. This is only the field's description

-- Explicit field paths for nested objects
body@object: This is the object root description  
body.profile@user_profile: Main user information  
body.metadata@object: Tracking info  
body.metadata.source?str#internal: Origin of signup  
body.metadata.labels?str[]#internal: Internal labels  
body.metadata.status@enum[active,inactive,archived]|null#deprecated: User status

-- To import an example, expect to be declared as 'ref@example'  
-- 'example' field will not be validated to follow the specified  
-- schema. Works as wildcard
example: <ref@example:create_user>

[ref@schema:user_profile]
name@str: User's name  
email@str: User's email  
-- For special types like 'datetime' declared here  
-- the parser will completely ignore its validation, since  
-- it is not supported by the expected basic type  
-- Works as a wildcard for types
timezone@custom:datetime: User's timezone

[ref@example:create_user]
profile.name: Jane Doe  
profile.email: jane@example.com  
profile.gender: female  
metadata.source: newsletter  
metadata.tags: [beta, test]

[[Final section@content#public]]
group: finishing
summary: <md>
# Final section I guess
This is a separate content section

## Foo
Foo indeed

## List of things

* Item 1
* Item 2
* Item 3
`

var manifestWithValidTags string = `
[[@tags]]
public#00FF00: Public API
internal#AAAAAA: Internal use only
deprecated#FF6F61: Will be removed soon
under-dev#FFD966: Still under development
`

var manifestWithInvalidTags string = `



[[@tags]]
public#00FF00: Valid tag format



public2#00FF00: Valid tag format 2
internal@AAAAAA: Invalid tag format
#AAAAAA: Missing name data
`

var manifestWithEmptyTags string = `
[[@tags]]

`

var manifestWithValidModule string = `
[[@groups]]
api: API

[[Users@api]]
group: api
summary: Endpoints related to user operations

[[Items@api]]
group: api
summary: Endpoints related to items owned by users

[[invalidModule@api]]
group: api
summary: This should never be loaded due syntax error
`

var manifestWithValidTaggedModules string = `
[[@groups]]
api: API

[[@tags]]
public#00FF00: Public API
internal#AAAAAA: Internal use only
deprecated#FF6F61: Will be removed soon
under-dev#FFD966: Still under development

[[Users@api#public,under-dev]]
group: api
summary: Endpoints related to user operations

[[Items@api#internal]]
group: api
summary: Endpoints related to items owned by users
`

var manifestWithUnexistentTaggedModules string = `
[[@tags]]
public#00FF00: Public API
internal#AAAAAA: Internal use only

[[@groups]]
api: API

[[Users@api#public,under-dev]]
group: api
summary: Endpoints related to user operations
`

var ValidManifestWithInlineComments string = `
-- tags: Special section
-- items follows: '<name>#<color>: <description>' syntax
[[@tags]]
public#00FF00: Public API  
internal#AAAAAA: Internal use only  
deprecated#FF6F61: Will be removed soon  
under-dev#FFD966: Still under development  
beta#87CEEB: Beta feature  
admin#FF1493: Admin only

-- module: Special section
-- header follows: '[[<name>#<tag1,tag2,...>]]'
-- fields follow: '<name>(@|?)<type>#<tag1,tag2,...>: <description>'
`

var ValidManifestWithMultilineComments string = `
<--
tags: Special section
items follows: '<name>#<color>: <description>' syntax
-->
[[@tags]]
public#00FF00: Public API  
internal#AAAAAA: Internal use only  
deprecated#FF6F61: Will be removed soon  
under-dev#FFD966: Still under development  
beta#87CEEB: Beta feature  
admin#FF1493: Admin only

<--
module: Special section
header follows: '[[<name>#<tag1,tag2,...>]]'
fields follow: '<name>(@|?)<type>#<tag1,tag2,...>: <description>'
-->
`

var ValidManifestWithContentSection string = `
[[@groups]]
getting-started#public: getting started

[[About amauta@content#public]]
group: getting-started
summary: <md>
# Amauta
Welcome to Amauta. The right tool to write documentation. This line is excesively long just to verify how line breaks render once that they load on templates

## What is Amauta?
Amauta is a non-standardized documentation tool.
</md>
`

var manifestWithValidGroup string = `
[[@groups]]
getting-started#public: getting started
api#public: client api
`

var manifestWithRootDetails string = `
[[@root]]
LogoUrl: https://raw.githubusercontent.com/luislve17/amauta/refs/heads/main/assets/amauta-logo.svg
GithubUrl: https://github.com/luislve17/amauta
`

var manifestWithComplexMdContent string = `
[[@groups]]
custom-section#public: Custom test section

[[Content@content]]
group: custom-section
summary: <md>
` + "# Content" + "\n```toml" + `
[[<name>@content#<tag_ids>]]
group: <group-id>
summary: <summary>
` + "\n```\n" + "This is an inline code `<md>` and `</md>`\n" + `

### Example
` + "\n```toml" + `
[[Tax management@content#internal]]
group: finances
summary: \<md\>
# Support for markdown!
## Subtitle
List:
1. One
2. Two
3. Three

> This is a quote

| tables | also | work |
|--------|------|------|
| tables | also | work |
\</md\>
` + "\n```" + `
</md>
`
