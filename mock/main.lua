#!/usr/local/bin/lua

local mark = require("markdown")
local toml = require("toml")

-- get_posts takes a dir path param and returns a table of post objects for
-- every post contained in the directory
function get_posts(dir)
  -- get_post takes a path and returns a post object
  local function get_post(path)
    -- trim removes leading and trailing whitespace (like python's strip())
    local function trim(text)
      return (text:gsub("^%s*(.-)%s*$", "%1"))
    end

    -- get_meta parses the toml frontmatter of the raw post markdown and returns
    -- a meta object for use later in sorting and organization
    local function get_meta(text)
      return toml.parse(trim(text:sub(5, text:find("+++", 4) - 2)))
    end

    -- get_html parses the markdown sans the frontmatter and generates html
    local function get_html(text)
      return trim(mark(text:sub(text:find("+++", 4) + 3)))
    end

    local file = io.open(path, "r"):read("*a")
    return {meta = get_meta(file), html = get_html(file)}
  end

  -- get_list takes a path and returns a list of post paths
  local function get_list(path)
    local list = {}
    for line in io.popen("ls " .. path):lines() do
      list[#list + 1] = path .. line
    end
    return list
  end

  local posts = {}
  for i, file in pairs(get_list(dir)) do posts[i] = get_post(file) end
  return posts
end
