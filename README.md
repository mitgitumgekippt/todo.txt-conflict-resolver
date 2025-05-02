# todo.txt-conflict-resolver

> [!WARNING]
> This tool is still under development and not yet production-ready. Features and behavior are subject to change. Double-check all results and report any bugs. Because of that, this tool generates a file named ".todo.txt.backup", which contains the content of the original "todo.txt" file and it does not automatically delete the conflicted files yet.

## Project Description

This Go-script provides a basic command-line interface for merging [todo.txt](https://github.com/todotxt/todo.txt) conflict files, if it's used with [syncthing](https://syncthing.net/) (or similar file synchronization tools).
Please note that it currently uses a custom sync algorithm that works for my personal needs and which assumes some assumptions. See the [Limitations](#Limitations) section for more details.

#### Background
I am using [todo.txt](https://github.com/todotxt/todo.txt) to organise my tasks.
To make it practical across devices, I sync my todo files using Syncthing.
However, I often modify these files independently on different devices while offline, which leads to conflicts.
This tool helps merge such conflicting versions back into a single, consistent file.

## Installation
Clone the project and build the project using
```go
go build
```

Then move the binary to the same location as your todo.txt file.

> [!IMPORTANT]
> This tool does *not* support file paths. It must be run in the directory containing your todo.txt.

## Usage
This script is intended to be run in CLI. To get started, run
```bash
./main --init
```

It will create the "[Common base](https://www.geeksforgeeks.org/git-merge/)", which is necessary for merging.

> [!IMPORTANT]
> This tool is not intended to work without this Common base.

> [!IMPORTANT]
> Make sure that no conflicts exist, when the "Common base" is created, as it is crucial for merging. If you generate the file while a conflict is present, the tool will not be able to detect lines/todos altered in both files (aka conflicts).

#### Runtime options:
The following options are available:
```
  --double-check 	 Confirm the changes before writing them. Implies --verbose. Overrides --dry.
  --dry              Perform a dry run without altering any file
  --force            Overwrite conflicts with the most recent entry (aka if every file has a different entry for the same todo)
  --init			 Initializes the common base. If multiple files names are specified, then the first one is taken. Stops after.
  --verbose          Shows the file contents and file changes
  --help             Show the help message and exits`)
  ```

#### Limitations
- This script assumes that no tasks are deleted. If you want to delete tasks [or move them to a done.txt file](https://github.com/mitgitumgekippt/archive_todo.txt), then please make sure that no current conflicts exist.
- Like git-merge, this tool does not automatically resolve conflicting lines. But it will detect them and print them into the console.

> [!NOTE]
> If you try to merge two files but a task was delted, the tool will think each line after the deleted line was altered.
