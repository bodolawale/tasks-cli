# Tasks CLI (Go)

A simple command-line todo application built with Go.

## Features

- Add tasks
- List tasks
- Mark tasks as complete/incomplete
- Delete tasks
- Persistent storage using CSV

## Usage

```bash
tasks add "Buy milk"
tasks list
tasks list -a
tasks complete 1
tasks uncomplete 1
tasks delete 1
```

## Example Output

```bash
ID  Task        Created
1   Buy milk    a minute ago
```

## Tech Stack

    •	Go
    •	encoding/csv
    •	text/tabwriter
