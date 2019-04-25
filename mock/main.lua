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

    return {
      meta = get_meta(file),
      html = get_html(file),
      path = path:sub(dir:len())
    }
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

  local function gen_post(data)
    local function gen_href(path)
      return path:sub(1, path:len()-3) .. "/"
    end

    return [[<li class="post">
      <div class="post__title">
        <a href="]] .. gen_href(data.path) .. [[">
          <span>]] .. data.meta.head .. [[</span>
        </a>
      </div>
      <div class="post__date">
      ]] .. data.meta.date .. [[
      </div>
    </li>]]
  end

  local list = ""

  for _, post in pairs(posts) do
    list = list .. gen_post(post)
  end

  return gen_map() .. [[<ul class="posts">]] .. list .. [[</ul>]]
end

print(gen_page("Karl McGuire", "index", gen_list(get_posts("posts/"))))
