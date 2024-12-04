# Simple-CLI-Go  

**Simple-CLI-Go** is a lightweight and user-friendly Command Line Interface (CLI) tool written in Go, designed for basic file operations and interaction. This tool simplifies text file management, making it a handy utility for everyday tasks.  

## Features  
- **`echo`**: Send a test message and get an immediate response.  
- **`load` & `close`**: Open and close files with ease (supports batch operations using `*`).  
- **`write` & `writeLn`**: Add content to files, either inline or with a newline, with support for multiple files.  
- **`read`**: View the contents of files directly in the terminal.  
- **`find`**: Search for specific strings within a file.  
- **`history`**: Keep track of executed commands.  
- **`clear`**: Erase the contents of a file (supports batch operations with `*`).  
- **`help`**: Display a concise guide to all available commands.  

## Getting Started  

Clone the repository and build the CLI to start managing your files with ease.  

```bash
# Clone the repo
git clone https://github.com/yourusername/Simple-CLI-Go.git

# Navigate to the project folder
cd Simple-CLI-Go

# Build and run
go build
./Simple-CLI-Go
```

## Usage
- Run the CLI and use the following commands:

```
  echo:     Responds back with a test message
  load:     Loads a file
  close:    Closes a file (can handle `*`)
  write:    Write to a file (can handle `*`)
  writeLn:  Write to a file with a new line (can handle `*`)
  read:     Reads content from a file (can handle `*`)
  find:     Find a string in a file
  history:  Shows the history of commands
  clear:    Clears a file (can handle `*`)
  help:     Shows this message
```
