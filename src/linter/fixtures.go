package linter

var ValidManifest string = `
[[tags]]
public#00FF00: Public API  
internal#AAAAAA: Internal use only  
deprecated#FF6F61: Will be removed soon  
under-dev#FFD966: Still under development  
beta#87CEEB: Beta feature  
admin#FF1493: Admin only

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

[[Items#public]]
description: Endpoints related to user's items

[request@GET:/v1/items]
description: Fetch user items  
contentType: application/json  
header~Authorization@str#internal: Bearer token  
query~limit?int: Max number of items  
query~offset?int: Pagination offset

[response@GET:/v1/items]
status: 200  
contentType: application/json  
body~items@item[]: List of items

[ref@schema:item]
id@str: Item identifier  
type@str: Item type  
created_at@str: ISO timestamp  
status@enum[active,inactive,archived]|null#deprecated: Item status

[[Orders#public,internal]]
description: Handle user orders and transactions

[request@POST:/v1/orders]
description: Submit a new order  
contentType: application/json  

header~Authorization@str#internal: Bearer token  
body~order@order: Order payload  
example: <ref@example:submit_order>

[response@POST:/v1/orders]
status: 201  
contentType: application/json  
body~confirmation@str: Confirmation message

[ref@schema:order]
id@str: Order ID  
items@str[]: Item IDs  
total@float: Total amount  
notes?str#internal: Optional notes

[[Reports#under-dev,beta]]
description: Generate usage and audit reports

[request@GET:/v1/reports/usage]
description: Get usage report  
header~Authorization@str#internal: Bearer token  
query~from@str: Start date  
query~to@str: End date

[response@GET:/v1/reports/usage]
status: 200  
body~report@usage_report: Usage report payload

[ref@schema:usage_report]
users@int: Total users  
active@int: Active users  
api_calls@int: Total API calls

[[Admin#admin]]
description: Admin endpoints

[request@DELETE:/v1/users/:id]
description: Delete a user  
header~Authorization@str#admin: Admin token

[response@DELETE:/v1/users/:id]
status: 204

[response@DELETE:/v1/users/:id]
status: 403  
body~error@str: Forbidden  
body~code@int: Error code

[ref@example:create_user]
profile:  
  name: Jane Doe  
  email: jane@example.com  
  gender: female  
metadata:  
  source: newsletter  
  tags: [beta, test]

[ref@example:submit_order]
order:  
  id: ORD123  
  items: [item-001, item-002]  
  total: 49.99  
  notes: Urgent delivery
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
