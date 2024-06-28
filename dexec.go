package main

import (
	"bufio"
	"fmt"
	ac "github.com/PeterHickman/ansi_colours"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type container struct {
	name   string
	id     string
	exited bool
}

func select_one(max int) int {
	for {
		s := fmt.Sprintf("Select [1-%d]: ", max)
		fmt.Print(ac.Bold(s))

		var i int
		fmt.Scanln(&i)

		if i >= 1 && i <= max {
			return i - 1
		}
	}
}

func fetch_containers() []container {
	cmd := exec.Command("docker", "ps", "-a")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return []container{}
	}

	list := []container{}
	sc := bufio.NewScanner(strings.NewReader(string(stdout)))
	for sc.Scan() {
		text := sc.Text()
		fields := strings.Fields(text)

		name := fields[len(fields)-1]

		if name != "NAMES" {
			id := fields[0]
			exited := strings.Contains(text, "Exited")

			list = append(list, container{name, id, exited})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].name < list[j].name
	})

	return list
}

func choose_container(title string, list []container) container {
	fmt.Println(ac.Bold(title))

	for i, e := range list {
		if e.exited {
			s := fmt.Sprintf("    %d %s (down)", i+1, e.name)
			fmt.Println(ac.Red(s))
		} else {
			s := fmt.Sprintf("    %d %s", i+1, e.name)
			fmt.Println(ac.Blue(s))
		}
	}

	i := select_one(len(list))
	return list[i]
}

func choose_command(title string, exited bool) string {
	fmt.Println(ac.Bold(title))

	var o []string
	if exited {
		o = []string{"start"}
	} else {
		o = []string{"stop", "restart", "shell"}
	}

	for i, e := range o {
		s := fmt.Sprintf("    %d %s", i+1, e)
		fmt.Println(ac.Blue(s))
	}

	i := select_one(len(o))
	return o[i]
}

func stop_start_restart(ssr string, id string) {
	cmd := exec.Command("docker", ssr, id)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func shell(id string) {
	cmd := exec.Command("docker", "exec", "-it", id, "bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // add error checking
}

func main() {
	cs := fetch_containers()

	if len(cs) == 0 {
		fmt.Println(ac.Red("There are no containers"))
		os.Exit(1)
	}

	c := choose_container("Available containers", cs)
	cmd := choose_command("Available commands", c.exited)

	switch cmd {
	case "stop", "start", "restart":
		stop_start_restart(cmd, c.id)
	case "shell":
		shell(c.id)
	}
}
