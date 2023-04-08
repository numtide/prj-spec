# PRJ Base Directory Specification

Authors:

- Jonas Chevalier [zimbatm@zimbatm.com](mailto:zimbatm@zimbatm.com)

Version 0.1

# Introduction

This document establishes a set of conventions for project-centric tools to follow. By adopting a common posture, tools can avoid re-inventing the same logic over and over again.

A project-centric tool is any software that operates primarily inside of a source tree. For example code formatters, code linters, build systems, ... This is in contrast with user-centric tools that are meant to be executed by the user, like applications. Or system-centric tools that are meant to be executed by the system, like system services.

The intended audience is developers who write such tools or users who might want to ask developers to adopt this as a standard.

The approach adopted by this document is similar in spirit to [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) which establishes conventions for user-centric tools for where to put their configuration and cache in the users’s homes.

# Conventions

The PRJ Base Directory Specification is based on the following conventions:

## Project root

The project root is a single base directory that points to the root of the project. Typically that directory would be the root of the source tree checked out by the SCM (git, svn, ...).

Most tools need to find the project root in order to find their configuration file or to traverse it (for example to format all of the files).

If a tool needs a project root, it should follow the following heuristic:

### Specification

The `PRJ_ROOT` MUST be an absolute path that points to the project root.

If the environment variable `$PRJ_ROOT` is set, tools MUST use it over any built-in heuristic to find the project root.

Otherwise, the tool is free to define its own heuristic.

### Example

A typical bash implementation, that looks for the `.config` folder as a fallback heuristic would be:

```bash
# Look for the .config folder in the current directory and up
find_prj_config() (
  local old_pwd
  while [[ $old_pwd != $PWD ]]; do
    if [[ -d .config ]]; then
      echo "$PWD"
      return 0
    fi
    old_pwd=$PWD
    cd ..
    if [[ $old_pwd = "$PWD" ]]; then
       # We're at the top and didn't find anything
       echo "ERROR: could not find project root" >&2
       return 1
    if
  done
)

: "${PRJ_ROOT:=$(find_prj_root)"
```

If the tool doesn’t have any configuration and can assume Git as the SCM, it might look like this:

```bash
: "${PRJ_ROOT:=$(git rev-parse --show-toplevel)"
```

## Config home

The config home represents the folder that stores the tool configuration file.

Historically that home would be the project root in most cases. We aim to change this to a `.config` sub-folder relative to the project root.

### Specification

If the environment variable `$PRJ_CONFIG_HOME` is set, tools MUST use it to look for their configuration files.

Otherwise, if the environment variable `$PRJ_ROOT` is set, the tool MUST use it to look for configuration files in `$PRJ_ROOT/.config`.

Otherwise, the tool is free to pick its own logic. Which might include deciding to abort.

### Example

A typical bash implementation would be:

```bash
: "${PRJ_CONFIG_HOME:=${PRJ_ROOT}/.config}"
my_config_file=${PRJ_CONFIG_HOME}/my_tool_name.ext
```

### Recommendation

In most cases, it’s best if the tool doesn’t load user-level configuration files to avoid differences between machines.

## Project ID

The project ID is an optional unique identifier for the project. It’s mostly useful combined with other things.

### Specification

The PRJ_ID value MUST pass the following regular expression: `^[a-zA-Z0-9_-]{,32}$`. It can be a UUIDv4 or some other random identifier.

If the environment variable `$PRJ_ID` is set, tools MUST use it.

Otherwise, if the `PRJ_CONFIG_HOME` is set and a `prj_id` file exists, it would load it after stripping any trailing white spaces.

Otherwise, the tool is free to pick its own logic.

### Example

Bash implementation:

```bash
if [[ -z "${PRJ_ID:-}" && -f "${PRJ_CONFIG_HOME/prj_id}" ]]; then
  PRJ_ID=$(< "${PRJ_CONFIG_HOME}/prj_id")
fi
```

### TODO

Consider adding some sort of fallback based on `PRJ_ROOT | sha256sum`?

## Cache home

The cache home represents the folder that stores intermediate results that can be re-created from other sources. It SHOULD be safe for the user to delete that folder at any time and not lose any information.

### Specification

If the environment variable `$PRJ_CACHE_HOME` is set, tools MUST use it. The value MUST be an absolute path.

Otherwise, if the `$PRJ_ID` is set, tools MUST set it to `${XDG_CACHE_HOME}/prj/${PRJ_ID}`.

Otherwise, the tool MUST set it to `${PRJ_ROOT}/.cache`.

### Example

Bash implementation

```bash
: "${XDG_CACHE_HOME:=${HOME}/.cache}"

if [[ -z "${PRJ_CACHE_HOME:-}" ]]; then
  if [[ -n "${PRJ_ID:-}" ]]; then
    PRJ_CACHE_HOME="${XDG_CACHE_HOME}/prj/${PRJ_ID}"
  else
    PRJ_CACHE_HOME="${PRJ_ROOT}/.cache"
  fi
fi
```

### Rationale

Cache directories often contain a deep file structure that can hit some filesystem limits on some system. Given that the user’s home is usually higher in the file hierarchy, it’s best to place it there. Another factor is that users might have multiple checkouts of a project and that allows sharing the build caches.
