# Note

Note is, yet another, note-taking tool. It's built with the idea of
organization of notes in mind. It specifically caters to people who are in the
habit of creating notes for different projects their working on.

## Features

- Concise command line interface.
- Allows multiple format of notes, normal dump of thoughts, to-do lists,
  bookmarks and github-like issues.
- Provides a TUI form for quick note-taking and better user experience.
- Keeps track of you projects. Allows you to write notes for any project from
  anywhere.
- Is git aware, will always write the notes in the base of the repository.
- Prioritises local notes, and support global notes.

## Installation

1. Clone the repository.
2. Make sure you have `go` and `make` installed.
3. Make sure you have `$GOPATH/bin` in your path.
4. Run `make install` in the root of the repository.

## Usage

- Taking simple notes

```sh
note This is a new note, You do not even need to use quotes. Unless, \
  ofcourse, I am using a special character like '?' or "'" or '"'.
```

- Bookmarking links

```sh
note bookmark https://github.com/Chaitanyabsprip/note

# You can also use short forms of subcommands
note b https://github.com/Chaitanyabsprip/note

# bookmark also has a TUI form option. You can invoke it with the following
# command
note b
```

- Making to-do lists

```sh
note todo I need to get this done

# And of course, with the short forms
note td I need to get this done
note t I need to get this done
```

- Local issues

```sh
note issue --title "The title of the issue needs to be quoted" However the \
  description of the issue does not need to be. This is amazing\!

# And again, the short forms
note i # this will invoke the TUI form.
```

You can invoke the `-h` flag for the main program or any subcommand to know its
CLI.

## Why, yet another, note-taking tool?

I am lazy and did not want to search for a tool and find the one that fulfills
my needs and method of working. I did not want to go through the pain of reading
a bunch of readme-s, usage documentation and trying out tools to find my fit.
And yes I realize I had to, and I am having to do that any way during the
development of this project. I enjoy learning and every project is a new
opportunity for exactly that. Having to go through readme-s of libraries,
tools, documentation and the whole process of planning, designing and
implementing a product is valuable even outside this project.
