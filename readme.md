# Draft

Draft is, yet another, note-taking tool. It's built with the idea of
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

1. Clone the repository
2. Make sure you have `go` and `make` installed

## Usage

- Taking simple notes

```sh
draft This is a new note, You do not even need to use quotes. Unless, \
  ofcourse, I am using a special character like '?' or "'" or '"'.
```

- Making link bookmarks

```sh
draft bookmark https://github.com/Chaitanyabsprip/draft

# You can also use short forms of subcommands
draft b https://github.com/Chaitanyabsprip/draft

# bookmark also has a TUI form option. You can invoke it with the following
# command
draft b
```

- Making to-do lists

```sh
draft todo I need to get this done

# And ofcourse, with the short forms
draft td I need to get this done
draft t I need to get this done
```

- Local issues

```sh
draft issue --title "The title of the issue needs to be quoted" However the \
  description of the issue does not need to be. This is amazing\!

# And again, the short forms
draft i # this will invoke the TUI form.
```

You can invoke the `-h` flag for the main program or any subcommand to know its
CLI.
