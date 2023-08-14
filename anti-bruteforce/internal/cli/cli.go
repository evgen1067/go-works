package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

var commandApp *AppCLI

var suggestions = []prompt.Suggest{
	{
		Text:        "rb",
		Description: "Reset Leaky Bucket",
	},
	{
		Text:        "add",
		Description: "Add address to Blacklist/Whitelist (syntax: add blacklist/whitelist 127.0.0.1/25)",
	},
	{
		Text:        "delete",
		Description: "delete address from Blacklist/Whitelist (syntax: delete blacklist/whitelist 127.0.0.1/25)",
	},
	{
		Text:        "exit",
		Description: "Exiting the program",
	},
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

func executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "rb":
		err := commandApp.ResetBucket()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	case "add":
		if len(blocks) < 3 {
			fmt.Println("Insufficient number of arguments")
		} else {
			err := commandApp.AddToList(blocks[1], blocks[2])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case "delete":
		if len(blocks) < 3 {
			fmt.Println("Insufficient number of arguments")
		} else {
			err := commandApp.DeleteFromList(blocks[1], blocks[2])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	default:
		fmt.Println("Unsupported command")
	}
}

func Run(ctx context.Context) {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("> "),
		prompt.OptionTitle("Bruteforce CLI"),
		prompt.OptionInputTextColor(prompt.Yellow),
	)
	commandApp = NewAppCLI(ctx)
	p.Run()
}
