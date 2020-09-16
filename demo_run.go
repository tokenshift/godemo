package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/kballard/go-shellquote"
)

type DemoRun struct {
	Step      int
	Steps     DemoStepList
	Variables DemoVariables
}

type DemoAction func(DemoStepList, DemoVariables) (DemoAction, bool, error)

var (
	CommandFormat = color.New(color.FgGreen, color.Bold)
	CommentFormat = color.New(color.FgCyan)
	ErrorFormat   = color.New(color.FgHiRed, color.Bold)
)

func DisplayStepN(stepIndex int) DemoAction {
	return func(steps DemoStepList, vars DemoVariables) (DemoAction, bool, error) {
		if stepIndex < 0 {
			stepIndex = 0
		}

		if stepIndex >= len(steps) {
			return nil, false, nil
		}

		step := steps[stepIndex]

		// Display the comment if there is one
		if step.Comment != "" {
			CommentFormat.Printf("(%d) # %s\n", stepIndex+1, strings.TrimSpace(step.Comment))
		}

		// Echo the command to be run (sanitized) if there is one
		if step.Cmd != "" {
			cmd := os.Expand(strings.TrimSpace(step.Cmd), vars.sanitize)
			CommandFormat.Printf("(%d) $ %s\n", stepIndex+1, cmd)
		}

		// If there was neither, move on to the next step
		if step.Comment == "" && step.Cmd == "" {
			return DisplayStepN(stepIndex + 1), true, nil
		}

		// Otherwise, prompt for the next action to take.
		// Default next action is to execute the command that was just echoed.
		action, err := PromptForAction(stepIndex, ExecuteStepN(stepIndex))
		return action, err == nil, err
	}
}

func ExecuteStepN(stepIndex int) DemoAction {
	return func(steps DemoStepList, vars DemoVariables) (DemoAction, bool, error) {
		if stepIndex < 0 {
			stepIndex = 0
		}

		if stepIndex >= len(steps) {
			return nil, false, nil
		}

		step := steps[stepIndex]

		// If there's no command, just move on to display the next step.
		if step.Cmd == "" {
			return DisplayStepN(stepIndex + 1), true, nil
		}

		// Otherwise, execute the command and display its output, then prompt for
		// the next action (default is to display the next demo step).

		cmdString := os.Expand(strings.TrimSpace(step.Cmd), vars.mapping)

		args, err := shellquote.Split(cmdString)
		if err != nil {
			return nil, false, err
		}

		cmd := exec.Command(args[0], args[1:]...)
		if step.Echo {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		err = cmd.Run()
		if err != nil {
			ErrorFormat.Println("ERROR:", err)
		}

		action, err := PromptForAction(stepIndex, DisplayStepN(stepIndex+1))
		return action, err == nil, err
	}
}

func ExitDemo(DemoStepList, DemoVariables) (DemoAction, bool, error) {
	return nil, false, nil
}

func DisplayValidNextSteps() {
	fmt.Println(`Valid options:
(n)ext   - Proceed to next step
(p)prev  - Go back to the last step
(r)eplay - Repeat the same step
(q)uit   - Exit the demo
{number} - Goto a specific step
`)
}

func PromptForAction(currentIndex int, defaultAction DemoAction) (DemoAction, error) {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return ExitDemo, err
	}

	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "":
		return defaultAction, nil
	case "next", "n", "forward", "f":
		return DisplayStepN(currentIndex + 1), nil
	case "previous", "prev", "p", "back", "b":
		return DisplayStepN(currentIndex - 1), nil
	case "replay", "r", "same", "s":
		return DisplayStepN(currentIndex), nil
	case "quit", "q", "exit", "x":
		return ExitDemo, nil
	default:
		match, _ := regexp.MatchString(`^\d+`, input)
		if match {
			// Enter a number to go to the numbered step
			nextStep, err := strconv.ParseInt(input, 10, 0)

			if err != nil {
				fmt.Println(input, "is not a valid option.")
				DisplayValidNextSteps()
				return PromptForAction(currentIndex, defaultAction)
			}

			return DisplayStepN(int(nextStep) - 1), nil
		} else {
			fmt.Println(input, "is not a valid option.")
			DisplayValidNextSteps()
			return PromptForAction(currentIndex, defaultAction)
		}
	}
}
