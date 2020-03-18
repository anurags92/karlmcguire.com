#!/usr/local/bin/lua

local mark = require("markdown")
local toml = require("toml")
local date = require("date")

local VIEW_SERVER = "http://api.karlmcguire.com"

-- get_files returns a table list of filenames inside the directory
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
      return toml.parse(trim(text:sub(5, text:find("---", 4) - 2)))
    end

    -- get_html parses the markdown sans the frontmatter and generates html
    local function get_html(text)
      return trim(mark(text:sub(text:find("---", 4) + 3)))
    end

    local file = io.open(path, "r"):read("*a")

    return {
      meta = get_meta(file),
      html = get_html(file),
      text = trim(file:sub(file:find("---", 4) + 3)),
      href = path:sub(dir:len()):sub(1, -4) .. "/"
    }
  end

  local posts = {}
  for i, file in pairs(get_files(dir)) do 
    local post = get_post(file)
    if post.meta.published then
      posts[i] = post
    end
  end
  return posts
end

function gen_page(title, active, main)
  local function gen_head(title)
    return [[<head>
      <meta charset="utf-8">
      <meta http-equiv="X-UA-Compatible" content="IE=edge">
      <meta name="viewport" content="width=device-width, initial-scale=1">
      <meta name="format-detection" content="date=no,telephone=no">
      <title>]] .. title .. [[</title>
      <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto|Roboto+Mono">
      <link rel="stylesheet" href="/index.css">
      <link rel="icon" type="image/png" href="/favicon.png">
    </head>]]
  end

  local function gen_body(header, main)
    return [[<body>
      <div class="wrap">
        ]] .. header .. main .. [[
        <footer>
          &copy; 2020 Karl McGuire
          <span>Powered by <a href="https://lua.org">Lua</a>.</span>
        </footer>
      </div>
      <script type="text/javascript">
        fetch("]] .. VIEW_SERVER .. [[/views?path=" + window.location.pathname)
          .then((res) => {
            return res.json() 
          })
          .then((data) => {
            document.getElementById("views").innerText = data 
          })
      </script>
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
          <a href="https://github.com/karlmcguire"><svg width="20" height="20" xmlns="http://www.w3.org/2000/svg"><path d="M10 0C4.477 0 0 4.477 0 10c0 4.418 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.342-3.369-1.342-.454-1.155-1.11-1.462-1.11-1.462-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.268 2.75 1.026A9.578 9.578 0 0 1 10 4.836c.85.004 1.705.114 2.504.337 1.909-1.294 2.747-1.026 2.747-1.026.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.579.688.481C17.137 18.163 20 14.418 20 10c0-5.523-4.478-10-10-10" fill="#C63535" fill-rule="evenodd"/></svg></a>
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
        delay = "opacity: " .. math.random() + 0.10 .. ";"
        boxes = boxes ..
          [[<div id="box__]] .. (i .. "__" .. a) ..
          [[" class="map__box" style="]] .. delay .. [["></div>]]
      end
      boxes = boxes .. [[</div>]]
    end

    return [[
    <div class="map">
    ]] .. boxes .. [[
    </div>
    ]]
  end

  -- get_post creates a list item with the post href, title, and date
  local function gen_post(data)
    return [[<li class="post">
      <div class="post__title">
        <a href="]] .. data.href .. [[">
          <span>]] .. data.meta.title .. [[</span>
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

-- gen_post creates the main element for individual post pages
function gen_post(post)
  words = 0
  for word in post.text:gmatch("%w+") do
    words = words + 1
  end

  return [[<div class="head">
    <h1 class="title"><a href="]] .. post.href .. [[">]]
    .. post.meta.title .. [[</a></h1>
    <div class="head__meta">
      <span>]] .. post.meta.date .. [[ &mdash; ]]
      .. string.format("%.0f min read", words / 265) .. [[
      </span>
      <div class="head__meta__views">
        <span id="views" class="head__meta__views__count"></span>
        <svg class="head__meta__views__eye" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 576 512"><path fill="currentColor" d="M572.52 241.4C518.29 135.59 410.93 64 288 64S57.68 135.64 3.48 241.41a32.35 32.35 0 0 0 0 29.19C57.71 376.41 165.07 448 288 448s230.32-71.64 284.52-177.41a32.35 32.35 0 0 0 0-29.19zM288 400a144 144 0 1 1 144-144 143.93 143.93 0 0 1-144 144zm0-240a95.31 95.31 0 0 0-25.31 3.79 47.85 47.85 0 0 1-66.9 66.9A95.78 95.78 0 1 0 288 160z"></path></svg> 
      </div>
    </div>
  </div>
  <div class="text">
  ]] .. post.html .. [[
    <div class="notice">
      If you enjoyed this writing and want to get notified when I post more,
      or see some of my shorter thoughts, check out my <a href="https://t.me/karlsmcguire">Telegram channel.</a>
    </div>
  </div>]]
end

for _, post in pairs(get_posts("./posts/")) do
  os.execute("rm -rf ./docs" .. post.href)
  os.execute("mkdir ./docs" .. post.href)
  local file = io.open("./docs" .. post.href .. "index.html", "w+")
  file:write(gen_page(post.meta.title, "", gen_post(post)))
end

local index = io.open("./docs/index.html", "w+")
index:write(gen_page("Karl McGuire", "index", gen_list(get_posts("./posts/"))))
