#!/usr/local/bin/lua -i

local toml = require("toml")
local temp = toml.parse('head = "Hello, world!"')

print(temp["head"])
