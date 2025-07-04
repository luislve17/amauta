package linter

var ValidManifest string = `
-->
This is a multiline comment block
Starts with --> and ends with <--
but both must be on a new line with nothing
else but that characters sequence.
<--

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
[[Users@api#public,internal]]
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

-- Multiline object declaration
body: This is the object root description
-- Here the type of the field is custom 'user_profile', and must
-- import from other declaration of type 'ref@schema'.
.profile@user_profile: Main user information  
.metadata@object: Tracking info  
.metadata
..source?str#internal: Origin of signup  
..labels?str[]#internal: Internal labels  

-- To import an example, expect to be declared as 'ref@example'
-- 'example' field will not be validated to follow the specified
-- schema. Works as wildcard
example: <ref@example:create_user>

[ref@schema:user_profile]
name@str: User's name
email@str: User's email
-- For special types like 'datetime' declared here
-- the parser will completely ignore its validation, since
-- is not supported by the expected basby the expected basic type
-- Works as a wildcard for types
timezone@@datetime: User's timezone

[ref@example:create_user]
profile:  
.name: Jane Doe  
.email: jane@example.com  
.gender: female  
metadata:  
.source: newsletter  
.tags: [beta, test]
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
[[Users@api]]
description: Endpoints related to user operations

[[Items@api]]
description: Endpoints related to items owned by users

[[invalidModule@api]]
description: This should never be loaded due syntax error
`

var manifestWithValidTaggedModules string = `
[[@tags]]
public#00FF00: Public API
internal#AAAAAA: Internal use only
deprecated#FF6F61: Will be removed soon
under-dev#FFD966: Still under development

[[Users@api#public,under-dev]]
description: Endpoints related to user operations

[[Items@api#internal]]
description: Endpoints related to items owned by users
`

var manifestWithUnexistentTaggedModules string = `
[[@tags]]
public#00FF00: Public API
internal#AAAAAA: Internal use only

[[Users@api#public,under-dev]]
description: Endpoints related to user operations
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
