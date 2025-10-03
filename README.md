# Gotem
Gotem is a cli tool to create, store, and send HTTP requests.
This tool is meant to be a simple way to handle HTTP requests using a custom *.gohttp file.

# Table of contents
- [Installation](#installation)
- [GOHTTP](#gohttp)
- [FAQ](#faq)
- [Features](#features)

# Installation
Currently this project is still pre-realease all version before 1.0.0 are going to be breaking changes. This also means we dont have a proper realease thus you will need go to run the CLI.\
Run this command to install:
```bash
go get github.com/wvan1901/Gotem
go install github.com/wvan1901/Gotem
```

# GOHTTP
An HTTP request is message, this message needs to follow a certain format for it to be a valid request. It starts with a request line followed by optinal headers and a body.\
GoHttp build off this format, for a request it will need a custom request line, RFC valid optional headers, RFC valid body, and new lines to indicate end of request.\
All sections will be refering to GoHttp.\
`NOTE:` There is a lot more to a http message than what I said, gohttp is based on HTTP RFC 9112.

## Labels
Labels is gohttp custom symantics.
A label is a string key value pair. This is the symatics of a label:
```
A LABEL begins on a new line with char '@' followed by optional multiple spaces (OMP), label name, OMP, equals, OMP and label value.
  @ OMP label-name OMP EQUAL label-value NEW-LINE
```
`NOTE:` Label names are case insensitive they also must be alphanumeric or dash or underscore (ANDU). Label values also must be ANDU.\
Labels are used to give custom metadata. This allows us to over ride strings in the http request. An example of this being useful is overriding base URL.

### Requred Labels
The only requirement needed other than a HTTP request is a name for the request. This is called a request name.
To give a request a request name we add a label with the label name as 'name' before the request. Ex:
```
@name=request-name
```
Remember label names are case insensitive so the label name for request name could also be `@NAME`.

### Optional labels
There is only one custom label so far and that is `description`. This is used by gotem to provide info when using the `-h` flag.

### Custom labels
Remeber a label is just a string key value pair.
So to create a custom label just follow the symantics. Here is a example of a custom label:
```
@user-name=gopher4life
# The Label name is: user-name
# The Label value is: gopher4life
```
Gotem has reserved labels so the following label names should not be used as a custom label: `name, description`.

### Label order
Labels must follow this order: First label must be the name label. That it we can put any custom labels between name label and HTTP request.

## Request Line
A request line consist of a request method & target. The request method & target is `REQUIRED`.
```
A request-line begins with a request method token, followed by a single space (SP), the request url, and ends with a new line.
  request-line   = method SP request-url new-line

Example:
GET http://localhost:42069/health

So Method=GET, Url=http://localhost:42069 ,Path=/health
```
The allowed Request methods are the following:
* GET
* POST
* DELETE
* PUT
* PATCH

`NOTE:` The request methods MUST be in all uppercase.\
The request url will be any characters between the SP after method and all chars until the new line
## Headers
Each header is optional leading white space (OLWS), a case-insensitive field name, a colon (`:`), OLWS, field value, optional trailing white space (OTWS)
```
OLWS field-name: OLWS field-body OTWS

Example:
Content-Type: text/html

So FieldName=content-type, FieldValue=text/html
```

## Request Body
After the headers a body is optional but if one is included an empty line is required between last header and body.
```
GET http://localhost:42069/health
Host: localhost:42069

hello world!
```

## Comments
Comments are simple they can be placed anywhere in the file with the sole exception of a HTTP Request.
Comments must also begin on a new line.
Here are the semantics of a comment:
```
A COMMENT begins on a new line with the char '#' followed by any string and ends with a NEW-LINE
  # ANY-STRING
```
`REMEMBER:` a comment cannot be inside label value of a REQUST so the following is invalid:
```
@Name=request1
@Description=health check
GET http://localhost:42069/health
# This is an INVALID line, comments cannot exist within a request.
# Gotem will attempt to parse the comment as part of the http request which will fail!
Host: localhost:42069

# Even worse it will take this line as part of the body and send it.
```
## Entire request
A valid http request must start with a name label and end with a HTTP request.
Here is an example of valid file:
```
# Some comment!
@Name=request2
@Description=submit
POST http://localhost:42069/submit
Host: localhost:42069
Content-Length: 13

hello world!

# Comment 2
@Name=health
@Description=health-check
# Comment before request
GET http://localhost:42069/health
Host: localhost:42069
```
`NOTE:` A empty new line is required between headers and body.

# FAQ
*How is this better than a, b, ... z tool?*\
Its not, this project is meant to suit my needs. Those tools are great and offers solutions to issues many people have.
I haven't found a tool that feels good to use with my development set up. This is why I made this.\
\
*Why do we need a config file?*\
Memorizing 100 requests is impossible, so I need to store then them in some format. This tool currently isn't meant for
creating a http request without a file. Other great tools already exsist that handle this pretty well (curl, httpie, ... etc).\
\
*Why is the input file a custom http file?*\
There is no standard for .gohttp file and different tools have different formats for them.
Http request follow a specific format, this allows the file to reflect closely to what is being sent to a server.\
\
*Why is the output JSON?*\
Its a preference since I like to use jq for some of my tasks. Having json as the output helps me quickly get the infomation I need and allows
me to use the output with other tools. This also allows me to chain HTTP request which I won't be implementing into this tool or in the near future.

# Features
Gotem is very young and has limited funtionallity. All checked items are feature already added, uncheck means I will start it someday.
### CLI
- [ ] Send & View Request
  - [X] Rest
  - [ ] GraphQL
- [ ] Variable support
- [ ] Cookie support
- [ ] Auth support
  - [ ] OAuth2
  - [ ] Basic
- [ ] Code generation
- [ ] Web socket support

### GOHTTP
- [ ] Syntax highlighting
- [ ] Http oriented scripting language (Probably never, lol)
