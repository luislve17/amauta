<div align="center">
  
![logo](https://github.com/luislve17/amauta/raw/refs/heads/main/assets/amauta-banner.svg)

<h1>
  Docs for people
</h1>

![Static Badge](https://img.shields.io/badge/version-alpha0.1-2b7573)
![Static Badge](https://img.shields.io/badge/engine-go-00ADD8?logo=go)
![Static Badge](https://img.shields.io/badge/ui-html-F06529?logo=html5)


</div>

# Links

1. Live doc: [link](https://luislve17.github.io/amauta/doc) (👷‍♂️⚠️Under development)

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
You may find the fully detailed protocol on our [our documentation](https://luislve17.github.io/amauta/doc) (👷‍♂️⚠️Under development)

## Linter/Renderer
Documenting expects to export a final page that is both professional and easy to navigate. Being honest, there are plenty platforms that do this quiet well, [after](https://www.dreamfactory.com/) you [click](https://redocly.com/) the "[talk to sales team](https://scalar.com/)" button [of course](https://stoplight.io/).

From a FOSS point of view, most alternatives to render a doc manifest into a page lack a lean, professional interface. And yes, the user could invest time and effort modifying the style to get an specific theme around the baseline that these tools offer, but with **_Amauta_** is expected to have available a responsive result that has a simplified but elegant look.

<details>
<summary>😕 What other free alternatives offer out of the box</summary>

<img width="1882" height="1092" alt="image" src="https://github.com/user-attachments/assets/c620e9af-1521-4b42-95c1-284395b264d6" />
<img width="1851" height="980" alt="image" src="https://github.com/user-attachments/assets/bb3036c2-4afc-4e9b-b5d4-b27ba1f9f320" />
<img width="1851" height="980" alt="image" src="https://github.com/user-attachments/assets/95c7012f-ebe4-4b0f-a4e1-b85f526126c8" />
  
</details>


<details>
<summary>:moyai: What Amauta offers (alpha0.1 version) </summary>
  
`amauta --render -i ./dist/manifest.amauta -theme default`

<img width="1689" height="1068" alt="image" src="https://github.com/user-attachments/assets/af95e9a0-2615-409c-a7f3-bd3ec3217929" />
<img width="649" height="1043" alt="image" src="https://github.com/user-attachments/assets/49e6bbbb-9152-4c54-b751-9e96f8189051" />


`amauta --render -i ./dist/manifest.amauta -theme dark`

<img width="1689" height="1068" alt="image" src="https://github.com/user-attachments/assets/d842c3cb-1844-4e02-a98e-3eb0e296c5b2" />
<img width="649" height="1043" alt="image" src="https://github.com/user-attachments/assets/f55fddb9-a675-45e1-9bed-9a8a7671f95e" />
  
</details>

Base themes support, responsive and elegant, out of the box, from a single cli run, into a single fully-embedded HTML file (with minor dependencies for fonts, for example) for you to self-host, toy around, customize with themes and include in your own pipeline. No free-trial, no credit-card reader, no "talk with sales" call to action, no hosting fee's.

## CLI

As mentioned, we expect you to just write the docs, and get your page back. A single, all-in-it `.html` for you to use as you wish. Some minor dependencies are expected for loading icons, fonts and code highlight. (👷‍♂️⚠️Under development)

```
amauta -v
Amauta: version alpha-0.2

amauta -h
Usage of Amauta CLI (alpha-0.2):
-v	Build version
-i	Input path
-o	Output HTML file path (defaults to './dist/doc.html')
-lint	Lint doc manifest
-render	Render HTML from doc manifest
-theme	Name of the selected theme
```

~Enjoy! 🎶

# FAQ

> 0. Why did you create this?
1. I dislike long yaml files
3. I dislike deep indentation
4. I dislike clunky websites that look from the earlier 2000's, [with some exceptions of course](https://www.spacejam.com/1996/jam.htm)
5. I haven't being hired yet

> 1. Who asks this questions? You don't have any real users yet...

The voices in my head.

> 2. Is this going to be free forever?

I expect this to be maintained only by myself, on my own time, while I try to figure out life in general. While I enjoy it, yes, it will be free.

> 3. How often do you plan to make releases?

Please refer to (2)

> 4. Is there a way I can contribute

You may create issues requesting features or reporting bugs, but since the development is in a really early stage for now, those are probably going to be addressed anyway and the issues section will get (optimistically speaking) floaded. PRs are also just for me for the moment.
Ideally and for now: using the tool, starring the project, leaving your toughts in an issue using the "Input" template and sharing the repo with friends and collegues is enough.
Thank you! 🙂

> 5. I want to give you my moneh tho

Here you go:

<a href='https://ko-fi.com/Q5Q0P976H' target='_blank'><img height='36' style='border:0px;height:36px;' src='https://storage.ko-fi.com/cdn/kofi6.png?v=6' border='0' alt='Buy Me a Coffee at ko-fi.com' /></a>

And thank you. Specially now, support is much appreaciated at home.
