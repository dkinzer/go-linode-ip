package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Linode struct {
	Label string
	Ipv4  []string
}

type App struct {
	Linodes         []Linode
	matchingLinodes []Linode
	chosenLinode    int
	exit            bool
}

func LinodesList() []Linode {
	cmd := exec.Command("linode-cli", "--json", "linodes", "list")
	jsonBlob, err := cmd.Output()

	var linodes []Linode

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(jsonBlob, &linodes)

	if err != nil {
		log.Fatal(err)
	}

	return linodes
}

func (app *App) MatchLinodes(name string) {
	app.matchingLinodes = nil

	for _, linode := range app.Linodes {
		if strings.Contains(linode.Label, name) {
			app.matchingLinodes = append(app.matchingLinodes, linode)
		}
	}
}

func (app *App) chooseLinode(choice int) {
	app.chosenLinode = choice
}

func (app App) Ipv4() string {
	index := app.chosenLinode

	if len(app.matchingLinodes) == 0 {
		return ""
	}

	if app.exit {
		return ""
	}

	return app.matchingLinodes[index].Ipv4[0]
}

func (app App) PrintChooseALinode() {
	fmt.Println()
	fmt.Println("Multiple linodes found:")
	fmt.Println("-----------------------")

	for i, linode := range app.matchingLinodes {
		fmt.Printf("(%v) %s\n", i, linode.Label)
	}

	fmt.Println()
	fmt.Println("[u] update matcher")
	fmt.Println("[e] exit")

	fmt.Println("Which linode or command?")
	fmt.Println()
}

func Fetch(name string, linodes ...Linode) string {
	if len(linodes) == 0 {
		linodes = LinodesList()
	}

	app := App{
		Linodes:         linodes,
		matchingLinodes: nil,
		chosenLinode:    0,
	}

	app.MatchLinodes(name)

	for {
		if len(app.matchingLinodes) > 1 {
			reader := bufio.NewReader(os.Stdin)
			app.PrintChooseALinode()

			char, err := reader.ReadString('\n')

			if err != nil {
				log.Fatal(err)
			}

			char = strings.TrimSuffix(char, "\n")

			if char == "e" {
				app.exit = true
				break
			}

			if char == "u" {
				fmt.Println("What matcher?")

				name, err = reader.ReadString('\n')
				name = strings.TrimSuffix(name, "\n")
				app.MatchLinodes(name)

				continue
			}

			i, err := strconv.Atoi(char)

			if err == nil && i >= 0 && i < len(app.matchingLinodes) {
				app.chooseLinode(i)
				break
			}
		} else {
			break
		}
	}

	return app.Ipv4()
}

func main() {
	if len(os.Args) == 1 {
		return
	}

	name := os.Args[1]

	fmt.Printf("%v", Fetch(name))
}
