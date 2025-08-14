<div align="center">
  
![logo](https://github.com/luislve17/amauta/raw/refs/heads/main/assets/amauta-banner.svg)

<h1>
  Beautiful docs made simple
</h1>

![Static Badge](https://img.shields.io/badge/version-alpha--0.6-2b7573)
![Static Badge](https://img.shields.io/badge/engine-go-00ADD8?logo=go)
![Static Badge](https://img.shields.io/badge/ui-html-F06529?logo=html5)

</div>

# Links

1. Live doc: [link](file:///home/luis/Documents/Programming/amauta/src/dist/doc.html#section-About%20amauta) (🚧 Under development)

# Introduction

**Amauta** is a documentation generator that prioritizes simplicity and readability. It introduces a clean, intuitive syntax for documenting RESTful APIs, SDKs, and any collaborative project without the complexity of traditional documentation formats.

## Key Features

✨ **Zero hosting required** - Generate a single, self-contained HTML file  
📝 **Human-readable syntax** - Write docs that are as easy to read as they are to write  
🚀 **No build tools needed** - No npm, node, or complex setup required  
🎨 **Beautiful out of the box** - Professional themes included by default  
📱 **Fully responsive** - Looks great on any device  

# Installation

## Option 1: Quick Install (Linux)

```bash
wget -qO - https://raw.githubusercontent.com/luislve17/amauta/release/install.sh | bash
```

You'll be suggested to move the binary to `/usr/local/bin` or any other `$PATH`-avaiable location.

## Option 2: Manual Download

Download the precompiled binary from our [releases page](https://github.com/luislve17/amauta/releases) and add it to your PATH.

Verify installation:
```zsh
$ amauta -v
Amauta: version alpha-1.0
```

# Why Amauta?

## Syntax Comparison

Traditional documentation formats can be verbose and difficult to maintain. Here's how Amauta compares:

<details>
  <summary>📋 Traditional YAML approach</summary>

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
  <summary>✨ Amauta's approach</summary>

  ```
[[Users@api#public,internal]]
group: api
summary: Endpoints related to user operations

[request@POST:/v1/users]
summary: Create a new user  
contentType: application/json  

header@object: Authentication and request headers
header.Authorization@str#internal: Bearer token

body@object: Request payload
body.profile@user_profile: Main user information
body.metadata@object: Tracking info
body.metadata.source?str#internal: Origin of signup
body.metadata.labels?str[]#internal: Internal labels  
body.metadata.status@enum[active,inactive,archived]|null#deprecated: User status

example: <ref@example:create_user>

[ref@schema:user_profile]
name@str: User's name  
email@str: User's email  
timezone@custom:datetime: User's timezone

[ref@example:create_user]
profile.name: Jane Doe  
profile.email: jane@example.com  
profile.gender: female  
metadata.source: newsletter  
metadata.tags: [beta, test]
  ```
</details>

Amauta's syntax focuses on readability and maintainability, reducing nesting while keeping all the power you need for comprehensive documentation.

## Beautiful Output

Generate professional documentation that looks great without additional styling:

<details>
<summary>🎨 Default Theme</summary>
  
`amauta --render -i ./dist/manifest.amauta -theme default`

<img width="1689" height="1068" alt="Default theme desktop view" src="https://github.com/user-attachments/assets/af95e9a0-2615-409c-a7f3-bd3ec3217929" />
<img width="649" height="1043" alt="Default theme mobile view" src="https://github.com/user-attachments/assets/49e6bbbb-9152-4c54-b751-9e96f8189051" />

</details>

<details>
<summary>🌙 Dark Theme</summary>

`amauta --render -i ./dist/manifest.amauta -theme dark`

<img width="1689" height="1068" alt="Dark theme desktop view" src="https://github.com/user-attachments/assets/d842c3cb-1844-4e02-a98e-3eb0e296c5b2" />
<img width="649" height="1043" alt="Dark theme mobile view" src="https://github.com/user-attachments/assets/f55fddb9-a675-45e1-9bed-9a8a7671f95e" />
  
</details>

And more themes to come soon!

# Usage

The CLI is designed to be simple and straightforward:

```bash
$ amauta -h
Usage of Amauta CLI (alpha-1.0):
-h      Show this help
-v      Build version
-i      Input path
-o      Output HTML file path (defaults to './dist/doc.html')
-lint   Lint doc manifest
-render Render HTML from doc manifest
-theme  Name of the selected theme (available: 'default', 'dark')
```

Generate your documentation:
```bash
amauta --render -i ./docs/api.amauta -theme default -o ./dist/documentation.html
```

That's it! You now have a beautiful, self-contained HTML file ready to host anywhere.

# Documentation

For the complete protocol specification and advanced features, visit our [documentation](https://luislve17.github.io/amauta/doc) (currently under development).

# Contributing

Amauta is in active development. Here's how you can help:

- ⭐ **Star the repository** to show your support
- 🐛 **Report bugs** by creating an issue
- 💡 **Request features** using our issue templates  
- 📢 **Share with colleagues** who might find it useful

## Support Development

If you find Amauta useful, consider supporting its development:

<a href='https://ko-fi.com/Q5Q0P976H' target='_blank'><img height='36' style='border:0px;height:36px;' src='https://storage.ko-fi.com/cdn/kofi6.png?v=6' border='0' alt='Buy Me a Coffee at ko-fi.com' /></a>

Your support helps maintain and improve Amauta!

# FAQ

**Q: What makes Amauta different from other documentation tools?**  
A: Amauta prioritizes human readability in both the source format and the generated output. You write in a clean, minimal syntax and get beautiful documentation without any additional configuration.

**Q: Is Amauta free?**  
A: Yes! Amauta is and will remain free and open source.

**Q: How stable is the alpha version?**  
A: While Amauta is in alpha, it's functional for creating documentation. The core features work well, though some advanced features are still being developed.

**Q: Can I customize the themes?**  
A: Theme customization is planned for future releases. Currently, we provide default and dark themes that work well out of the box.

**Q: How can I contribute?**  
A: The best way to contribute right now is by using Amauta, reporting any issues you find, and sharing feedback about your experience.

---

Made with ❤️ by a developer who believes documentation should be simple and beautiful.
