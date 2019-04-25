#!/usr/local/bin/lua

local mark = require("markdown")
local toml = require("toml")

-- trim removes the surrounding whitespace from a string
function trim(s)
  return (s:gsub("^%s*(.-)%s*$", "%1"))
end

-- post loads a post from path and returns a table representation of it, with
-- the corresponding meta information and parsed markdown as html
function post(path)
  local meta = function(text)
    return toml.parse(trim(text:sub(5, text:find("+++", 4) - 2)))
  end

  local html = function(text)
    return trim(mark(text:sub(text:find("+++", 4) + 3)))
  end

  local file = io.open(path, "r"):read("*a")

  return {meta = meta(file), html = html(file)}
end

-- list returns a table of all post relative paths
function list(path)
  local temp = {}

  for line in io.popen("ls " .. path):lines() do 
    temp[#temp + 1] = path .. line
  end

  return temp
end

-- posts returns a table of all posts
function posts(dir)
  local files = list(dir)
  local temp = {}

  for i, file in pairs(files) do
    temp[i] = post(file)
  end

  return temp
end

print(posts("posts/")[2].meta.head)
