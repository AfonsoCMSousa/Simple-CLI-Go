package main

import (
	comand "CLI-Go/cmd"
	promp "CLI-Go/pkg"
	"flag"
	"fmt"
	"os"
	"strings"
)

var commands = map[string]comand.Properties{}
var commandOrder []string

func registerCommand(cmd comand.Properties) {
	commandName := cmd.(*comand.Command).Name
	commands[commandName] = cmd
	commandOrder = append(commandOrder, commandName)
}

var history []string
var historyIndex int

var files []*os.File
var filesIndex int

func main() {
	promp := promp.NewPromp()
	command, bufferSize, err := promp.AgrFlag()
	if err != nil {
		fmt.Println("Error parsing command line arguments")
		return
	}

	history = make([]string, bufferSize)

	if command == "" {
		for {
			command, args := promp.CommandPrompt(history[:], files, &historyIndex, bufferSize)
			if command == "exit" {
				for i := 0; i < filesIndex; i++ {
					files[i].Close()
					fmt.Print("Closed file: ", files[i].Name(), "\n")
				}
				fmt.Printf("Exiting...\n")
				break
			}

			// verify if user is building a pipeline command
			if cmd, ok := commands[command]; ok {
				var runSuccess bool
				content := ""
				for i := 0; i < len(args); i++ {
					if args[i] == ">" {
						// run the first command with all arg untill the ">"
						content = cmd.Run(args[:i])
						// use the next argument for the 2nd argument of next command
						if i+1 < len(args) {
							nextCommand := args[i+1]
							if nextCmd, ok := commands[nextCommand]; ok {
								nextArgs := append([]string{args[i+2]}, content)
								content = nextCmd.Run(nextArgs)
								fmt.Print(content)
							} else {
								fmt.Printf("Command not found (%s)\n", nextCommand)
							}
						} else {
							fmt.Print(content)
							runSuccess = true
						}
						break
					}
				}

				if !runSuccess {
					content = cmd.Run(args)
					fmt.Print(content)
				}
			} else {
				fmt.Printf("Command not found (%s)\n", command)
			}
		}
	} else {
		if cmd, ok := commands[command]; ok {
			content := cmd.Run(flag.Args())
			fmt.Print(content)
		} else {
			fmt.Printf("Command not found (%s)\n", command)
		}
	}
}

func init() {
	registerCommand(&comand.Command{
		Name:        "echo",
		Description: "Responds back with a test message",
		UsageText:   "echo <message>",
		Execute: func(_args []string) string {
			if len(_args) > 0 {
				return strings.Join(_args, " ") + "\n"
			}
			return "Echoing Back!\n"
		},
	})

	registerCommand(&comand.Command{
		Name:        "load",
		Description: "Loads a file",
		UsageText:   "load \t<filename>",
		Execute: func(args []string) string {
			if len(args) == 0 {
				return "No file provided\n"
			}

			for i := 0; i < filesIndex; i++ {
				if files[i].Name() == args[0] {
					return "File already loaded\n"
				}
			}

			// file, err := os.Open(args[0])
			file, err := os.OpenFile(args[0], os.O_RDWR, 0755)
			if err != nil {
				return "File doesnt exist or something else happend\n"
			}

			files = append(files, file)
			filesIndex++

			return "File loaded\n"
		},
	})

	registerCommand(&comand.Command{
		Name:        "close",
		Description: "Closes a file (can handle `*`)",
		UsageText:   "close <filename>",
		Execute: func(args []string) string {
			if len(args) == 0 {
				return "No file provided\n"
			}

			if args[0] == "*" {
				for i := 0; i < filesIndex; i++ {
					files[i].Close()
					fmt.Printf("Closed file: %s\n", files[i].Name())
				}
				files = []*os.File{}
				filesIndex = 0
				return "Closed all files\n"
			}

			for i := 0; i < filesIndex; i++ {
				if files[i].Name() == args[0] {
					files[i].Close()
					files = append(files[:i], files[i+1:]...)
					filesIndex--
					return "Closed file\n"
				}
			}
			return "File not found\n"
		},
	})

	registerCommand(&comand.Command{
		Name:        "write",
		Description: "Write to a file (can handle `*`)",
		UsageText:   "write <filename> <content>",
		Execute: func(args []string) string {
			if len(args) < 2 {
				return "No file or content provided\n"
			}

			fileName := args[0]
			content := args[1]
			for j := 2; j < len(args); j++ {
				content += " " + args[j]
			}

			if fileName == "*" {
				for i := 0; i < filesIndex; i++ {
					fileInfo, err := files[i].Stat()
					if err != nil {
						return "Error getting file info\n"
					}

					if fileInfo.Size() > 0 {
						fmt.Printf("File %s already has content. Overwrite? (y/n) ", files[i].Name())
						asw, _ := promp.NewPromp().CommandPrompt(history[:], files, &historyIndex, 64)
						if asw == "y" || asw == "Y" || asw == "yes" || asw == "Yes" {
							files[i].Truncate(0)
							files[i].Seek(0, 0)
						} else {
							files[i].Seek(0, 2)
							content = "\n" + content
						}
					}

					_, err = files[i].WriteString(content)
					if err != nil {
						return "Error writing to file\n"
					}
				}
				return "Wrote to all files successfully\n"
			}

			for i := 0; i < filesIndex; i++ {
				if files[i].Name() == fileName {
					fileInfo, err := files[i].Stat()
					if err != nil {
						return "Error getting file info\n"
					}

					if fileInfo.Size() > 0 {
						fmt.Printf("File already has content. Overwrite? (y/n) ")
						asw, _ := promp.NewPromp().CommandPrompt(history[:], files, &historyIndex, 64)
						if asw == "y" || asw == "Y" || asw == "yes" || asw == "Yes" {
							files[i].Truncate(0)
							files[i].Seek(0, 0)
						} else {
							files[i].Seek(0, 2)
						}
					}

					_, err = files[i].WriteString(content)
					if err != nil {
						return "Error writing to file\n"
					}

					return "Wrote to file successfully\n"
				}
			}

			return "File not loaded or does not exist\n"
		},
	})

	registerCommand(&comand.Command{
		Name:        "writeLn",
		Description: "Write to a file with a new line (can handle `*`)",
		UsageText:   "writeLn <filename> <content>",
		Execute: func(args []string) string {
			if len(args) < 2 {
				return "No file or content provided\n"
			}

			fileName := args[0]
			content := args[1]
			for j := 2; j < len(args); j++ {
				content += " " + args[j]
			}
			content += "\n"

			if fileName == "*" {
				for i := 0; i < filesIndex; i++ {
					_, err := files[i].WriteString(content)
					if err != nil {
						return "Error writing to file\n"
					}
				}
				return "Wrote to all files successfully\n"
			}

			for i := 0; i < filesIndex; i++ {
				if files[i].Name() == fileName {
					fileInfo, err := files[i].Stat()
					if err != nil {
						return "Error getting file info\n"
					}

					if fileInfo.Size() > 0 {
						fmt.Printf("File already has content. Overwrite? (y/n) ")
						asw, _ := promp.NewPromp().CommandPrompt(history[:], files, &historyIndex, 64)
						if asw == "y" || asw == "Y" || asw == "yes" || asw == "Yes" {
							files[i].Truncate(0)
							files[i].Seek(0, 0)
						} else {
							files[i].Seek(0, 2)
						}
					}

					_, err = files[i].WriteString(content)
					if err != nil {
						return "Error writing to file\n"
					}

					// For some reason, go doesnt concatenate the string until the file is closed and reopened
					files[i].Close()

					file, err := os.OpenFile(args[0], os.O_RDWR|os.O_CREATE, 0755)
					if err != nil {
						return "Error saving file\n"
					}

					files[i] = file

					return "Wrote to file successfully\n"
				}
			}

			return "File not loaded or does not exist\n"
		},
	})

	registerCommand(&comand.Command{
		Name:        "read",
		Description: "Reads content from a file (can handle `*`)",
		UsageText:   "read <filename>",
		Execute: func(args []string) string {
			if len(args) == 0 {
				return "No file provided\n"
			}

			if args[0] == "*" {
				content := ""
				for i := 0; i < filesIndex; i++ {
					fileContent, err := os.ReadFile(files[i].Name())
					if err != nil {
						return "Error reading file\n"
					}
					content += fmt.Sprintf("\nContent of %s:\n%s", files[i].Name(), string(fileContent))
				}
				return content
			}

			for i := 0; i < filesIndex; i++ {
				if files[i].Name() == args[0] {
					content, err := os.ReadFile(files[i].Name())
					if err != nil {
						return "Error reading file\n"
					}
					return fmt.Sprint("\nContent of " + files[i].Name() + ":\n" + string(content))
				}
			}
			return "File not loaded or does not exist\n"
		},
	})

	registerCommand(&comand.Command{
		Name:        "find",
		Description: "Find a string in a file",
		UsageText:   "find <filename> <string>",
		Execute: func(args []string) string {
			if len(args) < 2 {
				return "No file or string provided\n"
			}

			fileName := args[0]
			content := ""

			if fileName == "*" {
				for i := 0; i < filesIndex; i++ {
					fileContent, err := os.ReadFile(files[i].Name())
					if err != nil {
						return "Error reading file\n"
					}
					fileContentStr := string(fileContent)
					finalString := fmt.Sprintf("\nFound in %s:\n", files[i].Name())

					//cut content into lines
					lines := strings.Split(fileContentStr, "\n")
					var wasFound bool = false

					for j := 0; j < len(lines); j++ {
						if strings.Contains(lines[j], args[1]) {
							finalString += "\t>>> [" + fmt.Sprint(j+1) + "] -> " + lines[j] + "\n"
							wasFound = true
						}
					}

					if wasFound {
						content += finalString
					}
				}

				if content == "" {
					return "String not found in any file\n"
				}

				return content
			}

			for i := 0; i < filesIndex; i++ {
				if files[i].Name() == fileName {
					fileInfo, err := files[i].Stat()
					if err != nil {
						return "Error getting file info\n"
					}

					if fileInfo.Size() == 0 {
						return "File is empty\n"
					}

					fileContent, err := os.ReadFile(files[i].Name())
					if err != nil {
						return "Error reading file\n"
					}

					content = string(fileContent)
					finalString := fmt.Sprintf("\nFound in %s:\n", files[i].Name())

					//cut content into lines
					lines := strings.Split(content, "\n")
					var wasFound bool = false

					for j := 0; j < len(lines); j++ {
						if strings.Contains(lines[j], args[1]) {
							finalString += "\t>>> [" + fmt.Sprint(j+1) + "] -> " + lines[j] + "\n"
							wasFound = true
						}
					}

					if !wasFound {
						return "String not found\n"
					}

					return finalString
				}
			}

			return "File not loaded or does not exist\n"
		},
	})

	registerCommand(&comand.Command{
		Name:        "history",
		Description: "Shows the history of commands",
		UsageText:   "history",
		Execute: func(args []string) string {
			content := "\nHistory:\n"
			for i, cmd := range history {
				if cmd == "" {
					break
				}
				content += fmt.Sprintf("\t%d: %s\n", i, cmd)
			}
			return content
		},
	})

	registerCommand(&comand.Command{
		Name:        "clear",
		Description: "Clears a file (can handle *)",
		UsageText:   "clear <filename>",
		Execute: func(args []string) string {
			if len(args) == 0 {
				return "No file provided\n"
			}

			if args[0] == "*" {
				for i := 0; i < filesIndex; i++ {
					files[i].Truncate(0)
					files[i].Seek(0, 0)
				}
				return "Cleared all files\n"
			}

			for i := 0; i < filesIndex; i++ {
				if files[i].Name() == args[0] {
					files[i].Truncate(0)
					files[i].Seek(0, 0)
					return "Cleared file\n"
				}
			}
			return "File not loaded or does not exist\n"
		},
	})

	registerCommand(&comand.Command{
		Name:        "help",
		Description: "Shows this message",
		UsageText:   "help <command>",
		Execute: func(args []string) string {
			if len(args) > 0 {
				if cmd, ok := commands[args[0]]; ok {
					return fmt.Sprintf("%s:\t%s\nUsage: %s", args[0], cmd.Help(), cmd.Usage())
				}
				return fmt.Sprintf("Command not found (%s)\n", args[0])
			}

			content := "\n"
			for _, cmdName := range commandOrder {
				cmd := commands[cmdName]
				content += fmt.Sprintf("%10s:\t%s\n", cmdName, cmd.Help())
			}
			return content
		},
	})
}
