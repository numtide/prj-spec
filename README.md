# PRJ Base Directory Specification

The PRJ Base Directory Specification sets a number of conventions for
project-centric tooling. With the ambition that *all* tools adopt it 😇

> [View document](./PRJ_SPEC.md)

**STATUS: unstable**

## Rationale

For the longest time, Unix tools would litter the user's home with dotfiles
that contain a mix of configuration and data files. As a user, it was hard to
know what files should be kept when moving between systems, or even to which
program a file belonged. It was also causing a lot of noise when listing the
directory.

Then one day, some gentle people decided to write the [XDG Base Directory
Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html).
It didn’t do anything on its own except lay down some conventions. Config
files should now live under `~/.config`, cached files would live under
`~/.cache`, and so on.

Nowadays, most programs respect this specification and it’s easier for users
to move their config around by copying the `~/.config` folder around. Or bust
the `~/.cache` folder when running low on disk space, ...

Similarly, code repositories nowadays are filled with top-level configuration
files. Build tools put build results next to config and code. And every tool
re-invents heuristics to find where the project root is.

This document sets a small number of conventions.

## Projects following this spec

In alphabetical order.

PRs to add your own are welcome!

- [direnv](https://direnv.net)
- [treefmt](https://github.com/numtide/treefmt)

## License

This document is licensed under the [CC0 1.0](./LICENSE).
