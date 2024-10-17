#!/bin/env python3

import os
import pathlib
import re
import sys
from collections import defaultdict

FUNC_PATTERN = re.compile(r"^\s*func\s+(\w+)")


def read_file(filename):
    with open(filename) as source:
        for index, line in enumerate(source):
            yield index + 1, line


def find_files(where, pattern):
    return tuple(
        sorted(x.relative_to(where) for x in pathlib.Path(where).rglob(pattern))
    )


def find_pattern(pattern, fileset):
    for filename in fileset:
        for number, line in read_file(filename):
            for item in pattern.finditer(line):
                yield f"{filename}:{number}", item.group(1)


def process(limit):
    functions = defaultdict(set)
    files = find_files(os.getcwd(), "*.go")
    for filename, function in find_pattern(FUNC_PATTERN, files):
        functions[function].add(filename)
    keys = "|".join(sorted(functions.keys()))
    pattern = re.compile(f"({keys})")
    counters = defaultdict(int)
    linerefs = defaultdict(set)
    width = 0
    for fileref, value in find_pattern(pattern, files):
        counters[value] += 1
        linerefs[value].add(fileref)
        width = max(width, len(fileref))
    for key, value in sorted(counters.items()):
        if key.startswith("Test"):
            continue
        definitions = len(functions[key]) - 1
        if value != limit + definitions:
            continue
        for link in sorted(linerefs[key]):
            fill = " " * (width - len(link))
            print(f"{link}{fill} {key}")


if __name__ == "__main__":
    process(int(sys.argv[1]) if len(sys.argv) > 1 else 1)
