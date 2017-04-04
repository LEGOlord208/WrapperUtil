# WrapperUtil

WrapperUtil is basically just a program wrapper.  
It can measure time to run, auto-restart, et.c.

But, it also has a handy API for developers to use.  
More of that later.

To see program arguments of WrapperUtil, run `./WrapperUtil --help`.

# API

The API is very small as of currently.  

## The detection API.
Detect WrapperUtil around your program by checking for the `wrapperutil` environment variable.  
It's value will always be `true`.  
(of course you can spoof this, but still)

## The action API
This API needs to be enabled using the `-p` flag.

Before you close the program, print `wrapperutil` and then a JSON:  
| Key     | | Type | Description                                  |
|---------| |------| -------------------------------------------  |
| Restart | | bool | Restarts the program. Less timeout than `-r` |
| Exit    | | bool | Exits the program, overriding the `-r` flag  |

Example:  
`wrapperutil{"Restart":true}`
