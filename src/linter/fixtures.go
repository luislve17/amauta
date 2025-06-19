package linter

var ValidManifest string = `
[[tags]]
public#00FF00: Public API
internal#AAAAAA: Internal use only
deprecated#FF6F61: Will be removed soon
under-dev#FFD966: Still under development

[[Users#public,internal]]
description: Endpoints related to user operations

[request@POST:/v1/users]
description: Create a new user
contentType: application/json

header~Authorization@str#internal: Bearer token
body~profile@user_profile: Main user information
body~metadata@object: Tracking info
body~metadata.source?str#internal: Origin of signup
body~metadata.tags?str[]#internal: Internal labels
example: <ref@example:create_user>

[response@POST:/v1/users]
status: 201
contentType: application/json

body~user@user_profile: Created user data

[response@POST:/v1/users]
status: 400
contentType: application/json
body~error@str: Error message
body~code@int: Error code

[ref@schema:user_profile]
name@str: Full name
email@str: Valid email
gender?enum[male, female, other]|null#deprecated: Gender selection
age?int|null#under-dev: Optional age

[ref@example:create_user]
profile:
  name: Jane Doe
  email: jane@example.com
  gender: female
metadata:
  source: newsletter
  tags: [beta, test]
`

var manifestWithValidTags string = `
[[tags]]
public#00FF00: Public API
internal#AAAAAA: Internal use only
deprecated#FF6F61: Will be removed soon
under-dev#FFD966: Still under development
`

var manifestWithInvalidTags string = `



[[tags]]
public#00FF00: Valid tag format



public2#00FF00: Valid tag format 2
internal@AAAAAA: Invalid tag format
#AAAAAA: Missing name data
`

var manifestWithEmptyTags string = `
[[tags]]

`

var manifestWithValidModule string = `
[[Users]]
description: Endpoints related to user operations

[[Items]]
description: Endpoints related to items owned by users

[[invalidModule]]
description: This should never be loaded due syntax error
`

var manifestWithValidTaggedModules string = `
[[tags]]
public#00FF00: Public API
internal#AAAAAA: Internal use only
deprecated#FF6F61: Will be removed soon
under-dev#FFD966: Still under development

[[Users#public,under-dev]]
description: Endpoints related to user operations

[[Items#internal]]
description: Endpoints related to items owned by users
`

var manifestWithUnexistentTaggedModules string = `
[[tags]]
public#00FF00: Public API
internal#AAAAAA: Internal use only

[[Users#public,under-dev]]
description: Endpoints related to user operations
`
