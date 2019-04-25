#!/usr/local/bin/lua

local mark = require("markdown")
local toml = require("toml")
local date = require("date")

function get_files(dir)
  local files = {}
  for file in io.popen("ls " .. dir):lines() do
    files[#files + 1] = dir .. file
  end
  return files
end

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

    return {
      meta = get_meta(file),
      html = get_html(file),
      href = path:sub(dir:len()):sub(1, -4) .. "/"
    }
  end

  local posts = {}
  for i, file in pairs(get_files(dir)) do posts[i] = get_post(file) end
  return posts
end

function gen_page(title, active, main)
  local function gen_head(title)
    return [[<head>
      <meta charset="utf-8">
      <meta http-equiv="X-UA-Compatible" content="IE=edge">
      <meta name="viewport" content="width=device-width, initial-scale=1">
      <title>]] .. title .. [[</title>
      <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto+Mono">
      <link rel="stylesheet" href="/index.css">
      <script src="/index.js"></script>
    </head>]]
  end

  local function gen_body(header, main)
    return [[<body>
      <div class="wrap">
        ]] .. header .. main .. [[
        <footer>
          &copy; 2019 Karl McGuire
          <span>Powered by <a href="https://lua.org">Lua</a>.</span>
        </footer>
      </div>
    </body>]]
  end

  local function gen_header(active)
    local function curr(path)
      if active == path then
        return [[ class="active"]]
      else
        return ""
      end
    end

    return [[<header>
      <nav>
        <div>
          <a href="/"]] .. curr("index") ..[[>Karl McGuire</a>
        </div>
        <div>
          <a href="/about/"]] .. curr("about") .. [[>About</a>
          <a href="/contact/"]] .. curr("contact") .. [[>Contact</a>
        </div>
      </nav>
    </header>]]
  end

  local head   = gen_head(title)
  local header = gen_header(active)
  local body   = gen_body(header, main)

  return [[<!doctype html><html lang="en">]] .. head .. body .. [[</html>]]
end

-- gen_list creates the main element for the index page with the list of all
-- posts in descending order (by date)
function gen_list(posts)
  local function gen_map()
    local boxes = ""

    for i = 1,7 do
      boxes = boxes .. [[<div class="map__row">]]
      for a = 1, 52 do
        boxes = boxes ..
          [[<div id="box__]] .. (i .. "__" .. a) .. 
          [[" class="map__box"></div>]]
      end
      boxes = boxes .. [[</div>]]
    end

    return [[
    <div class="map">
    ]] .. boxes .. [[
    </div>
    ]]
  end

  -- gen_post creates the main element for individual post pages
  local function gen_post(data)
    return [[<li class="post">
      <div class="post__title">
        <a href="]] .. data.href .. [[">
          <span>]] .. data.meta.head .. [[</span>
        </a>
      </div>
      <div class="post__date">
      ]] .. data.meta.date .. [[
      </div>
    </li>]]
  end

  -- sort by date descending
  function sort(t)
    -- collect keys from table
    local keys = {}
    for k in pairs(t) do keys[#keys + 1] = k end

    -- sort keys in descending order
    table.sort(keys, function(a, b)
      return date(t[a].meta.date) > date(t[b].meta.date)
    end)

    -- iterator
    local i = 0
    return function()
      i = i + 1
      if keys[i] then return keys[i], t[keys[i]] end
    end
  end

  -- concat each row of the post list, iterating in order of date descending
  local list = ""
  for _, post in sort(posts) do
    list = list .. gen_post(post)
  end

  return gen_map() .. [[<ul class="posts">]] .. list .. [[</ul>]]
end

function gen_post(post)
  return [[<div class="head">
    <h1 class="title"><a href="]] .. post.href .. [[">]]
    .. post.meta.head .. [[</a></h1>
    <p class="meta">]] .. post.meta.date .. [[</p>
  </div>
  <div class="text">
  ]] .. post.html .. [[
  </div>]]
end

for _, post in pairs(get_posts("posts/")) do
  local fold = "." .. post.href
  os.execute("rm -rf " .. fold)
  os.execute("mkdir " .. fold)

  local file = io.open(fold .. "index.html", "w+")
  file:write(gen_page(post.meta.head, "", gen_post(post)))
end

local index = io.open("./index.html", "w+")
index:write(gen_page("Karl McGuire", "index", gen_list(get_posts("posts/"))))
