package cli

var black = "\033[30m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var magenta = "\033[35m"
var cyan = "\033[36m"
var white = "\033[37m"

var reset = "\033[0m"
var bold = "\033[1m"
var dim = "\033[2m"
var italic = "\033[3m"
var underline = "\033[4m"
var blink = "\033[5m"
var reverse = "\033[7m"
var hidden = "\033[8m"

var buildVersion string = "alpha-0.5"

type styledString string
