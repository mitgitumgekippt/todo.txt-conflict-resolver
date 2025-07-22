<!-- Icon -->
<p align="center">
    <img src="icon.png" height="173"/></a>
</p>

<!-- Name of Project -->
<h1 align="center">todo.txt-conflict-resolver</h1>

<!-- Oneline, which explains purpose -->
<h4 align="center">
    Tool for automatically resolving sync conflicts of todo.txt files
</h4>


<div align="center">
    <a href="https://moritz-mander.de">My Home Page</a> |
    <a href="https://moritz-mander.de/">My Blog</a> |
    (Documentation)
</div>
<br>

 ![Maintained](https://img.shields.io/badge/Maintained-yes-yellow)
 ![Documentation](https://img.shields.io/badge/Documentation-yes-blue)

> [!WARNING]
> This tool is still under development and is far from being production-ready. Features and behavior are subject to change. Double-check all results and report any bugs. To recover in case of a fatal bug: A backup called ".todo.txt.backup" is created, which contains the content of the original "todo.txt" file. Additionally it does not automatically delete the conflicted files yet.

## Project Description

This Go-script provides a basic command-line interface for merging [todo.txt](https://github.com/todotxt/todo.txt) conflict files, if it's used with [syncthing](https://syncthing.net/) (or similar file synchronization tools).


> [!WARNING]
> Please note that it currently uses a custom sync algorithm that works for my personal needs and which assumes some assumptions. See the [Limitations](#Limitations) section for more details. I currently work on a implementation based on the 3way-merge-algorithm of git, which resolves most of the Limitations.

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


## Author
Moritz Mander

- [My Website](https://moritz-mander.de/)
- [My Blog](https://moritz-mander.de/blog/)
- [Github](https://github.com/mitgitumgekippt)

## License
This project is licensed by the Do What The Fuck You Want To Public License (WTFPL). For more details, check the LICENSE file or [wtfpl.net](https://www.wtfpl.net/about/)
