# Gotem
This is a cli tool that makes HTTP request.

# Table of contents
[How to use](#how-to-use)
[FAQ](#faq)

# How to use
This tool works by receiving a json config file that conatins multiple http requests.
User selects a request and Gotem makes that http request.
It returns the http response as a JSON format.
## Requirements
You need a config file with the following json format:
```json
[
	{
		"name": "request1",
		"description": "health check",
		"method": "GET",
		"url": "https://someurl/health",
		"headers": [
			{
				"key": "header-one",
				"values": ["1","2"]
			}
		]
	},
	{
		"name": "request2",
		"description": "...",
		"method": "POST",
		"url": "https://someurl"
	}
]
```
This file is an array of Request that Gotem uses to make requests.\
Gotem by default looks at the current directory for the config file `gotem.config.json`, if youd like to
provide a different file just use the `-f` flag to provide a filepath to another file.

## Make a request
In order to make a request you will need to select a request to use from the config file. You can select a
request by passing the `-req-name` flag. This selects the first request that matches that name.
`NOTE:` If there is only one request in the config then there is no need to pass a name, Gotem will use that request by default.

## Getting a response
Gotem will return the response in a custom JSON format.

## Listing all requests
Sometimes you'd like to see all requests inside your config. To achieve this use the `-ls` flag.
This flag will print a summary of all the requests in a friendly format.
`NOTE:` This flag will not make a http request. This interrupts the regular program flow.

# FAQ
*How is this better than a, b, ... z tool?*\
Its not, this project is meant to suit my needs. Those tools are great and offers solutions to issues many people have.
I haven't found a tool that feels good to use with my development set up. This is why I made this.\
*Why do we need a config file?*\
Memorizing 100 requests is impossible for me, so I need to store then them in some format. This tool currently isn't meant for
creating a http request without a config file. Other great tools already exsist that handle this pretty well (curl, httpie, ... etc).\
*Why is the config and output JSON?*\
Its a preference since I like to use jq for some of my tasks. Having json as the output helps me quickly get the infomation I need and allows
me to use the output with other tools.
