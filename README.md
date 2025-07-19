<div align="center">
  
![logo](https://github.com/luislve17/amauta/raw/refs/heads/main/assets/amauta-banner.svg)

<h1>
  Docs for people
</h1>

![Static Badge](https://img.shields.io/badge/version-alpha0.1-2b7573)
![Static Badge](https://img.shields.io/badge/engine-go-00ADD8?logo=go)
![Static Badge](https://img.shields.io/badge/ui-html-F06529?logo=html5)


</div>

# Introduction

**_Amauta_** focuses on fulfilling the need to create maintainable, redable and intuitive documentation. It defines a new protocol to declare documentation for RESTful APIs, marked-up content, SDKs and virtually anything that could be used in a collaborative team, open source tools or personal projects. 

# Content

## Protocol

**_Amauta_** defines a completely new way to declare the different pieces that exist in a documentation, covering syntax features that are expected from a regular documenation generator, but with simplicity and ease of use.

<details>
  <summary>🔴 What <strike>OpenAPI</strike> other protocols expect you to declare</summary>

  ```yaml
openapi: 3.1.0
info:
  title: Users API
  version: 1.0.0
  description: Endpoints related to user operations
tags:
  - name: api
    description: public, internal

paths:
  /v1/users:
    post:
      summary: Create a new user
      tags: [api]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                profile:
                  $ref: '#/components/schemas/UserProfile'
                metadata:
                  type: object
                  description: Tracking info
                  properties:
                    source:
                      type: string
                      description: Origin of signup
                    labels:
                      type: array
                      items:
                        type: string
                      description: Internal labels
                    status:
                      type: string
                      nullable: true
                      enum: [active, inactive, archived]
                      deprecated: true
                      description: User status
              required: [profile]
            example:
              profile:
                name: Jane Doe
                email: jane@example.com
                gender: female
              metadata:
                source: newsletter
                tags: [beta, test]
      parameters:
        - name: Authorization
          in: header
          required: true
          schema:
            type: string
          description: Bearer token. This is only the field's description
      responses:
        '200':
          description: OK

components:
  schemas:
    UserProfile:
      type: object
      properties:
        name:
          type: string
          description: User's name
        email:
          type: string
          description: User's email
        timezone:
          type: string
          format: datetime
          description: User's timezone
  ```
  
</details>

<details>
  <summary>🟢 What Amauta wants you to write</summary>

  ```
  [[Users@api#public,internal]]
group: api
summary: Endpoints related to user operations

[request@POST:/v1/users]
summary: Create a new user  
contentType: application/json  

header@object: This is the root @object description  
header.Authorization@str#internal: Bearer token. This is only the field's description

body@object: This is the @object root description  
body.profile@user_profile: Main user information. This field references a non-native type "@user_profile", to be searched within the linting scope.  
body.metadata@object: Tracking info. Below, it defines each of their fields
body.metadata.source?str#internal: Origin of signup.
body.metadata.labels?str[]#internal: Internal labels  
body.metadata.status@enum[active,inactive,archived]|null#deprecated: User status

example: <ref@example:create_user>

[ref@schema:user_profile]
name@str: User's name  
email@str: User's email  
timezone@custom:datetime: User's timezone. Defines a custom type, linted normally as wildcard to give the writer flexibility.

[ref@example:create_user]
profile.name: Jane Doe  
profile.email: jane@example.com  
profile.gender: female  
metadata.source: newsletter  
metadata.tags: [beta, test]

  ```
</details>

Focusing on simplicity and readability, **_Amauta_** avoids a nested/indented syntax. Some repetition to access inner fields might be encountered, while encouraging separate modular definitions on more complex entities.
You may find the fully detailed protocol on our documentation (👷‍♂️⚠️Under development)

## Linter/Renderer
Documenting expects to export a final page that is both professional and easy to navigate. Being honest, there are plenty platforms that do this quiet well, [after](https://www.dreamfactory.com/) you [click](https://redocly.com/) the "[talk to sales team](https://scalar.com/)" button [of course](https://stoplight.io/).
From a FOSS point of view, most alternatives to render a doc manifest into a web page to present, most solutions lack a lean, professional interface. **_Amauta_** is also improving this by providing a responsive product that uses a simplified but elegant look.

<details>
<summary>😕 What other free alternatives offer out of the box</summary>

<img width="1882" height="1092" alt="image" src="https://github.com/user-attachments/assets/c620e9af-1521-4b42-95c1-284395b264d6" />
<img width="1851" height="980" alt="image" src="https://github.com/user-attachments/assets/bb3036c2-4afc-4e9b-b5d4-b27ba1f9f320" />
  
</details>


<details>
<summary>:moyai: What Amauta offers </summary>
</details>
