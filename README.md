# Router Milter

This is a simple postfix milter that allows you to execute a javascript script
on the processing of any email going through postfix. This script can then allow
or reject the email based on it's return value.

Look at `config.example.yaml` and `script.example.js` for examples of how this
works, and `config/config.go` for a full list of config fields and
flags.

## Known Limitations

The JavaScript Engine ([goja](https://github.com/dop251/goja)) fully supports
ECMAScript 5.1, but does not have full EC6 support.

Additionally, a lot of browser-based or NodeJS based functions or libraries are
not available (like `setTimeout()` for example), as this is not running inside
NodeJS or a browser. Look at the goja README.md file for more details on this.
