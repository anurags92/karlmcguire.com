#!/usr/bin/env python3.7

import toml
from collections import namedtuple
from markdown2 import markdown as md

Post = namedtuple("Post", "meta html")
Meta = namedtuple("Meta", "head desc date tags")

def load(data):
    def meta(data):
        temp = toml.loads(data[4:data[4:].find("+++")+3])
        return Meta(head=temp["head"], desc=temp["desc"],
                    date=temp["date"], tags=temp["tags"])
    def html(data):
        return md(data[data[4:].find("+++")+8:]).strip()

    return Post(meta=meta(data), html=html(data))


print(load(open("post.md").read()).meta)
