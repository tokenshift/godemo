package main

import (
	"github.com/alecthomas/kong"
)

type DemoCommand struct {
	Filename string `kong:"required,help='Location of the demo definition file.',short='f',type='path'"`
}

type SetupCmd DemoCommand
type RunCmd DemoCommand
type TeardownCmd DemoCommand

func (r *SetupCmd) Run(ctx *kong.Context) error {
	demo, err := LoadDemoFile(r.Filename)

	if err != nil {
		return err
	}

	TitleFormat.Println("Setup:", demo.Title)
	return RunDemo(demo.Setup, demo.Variables)
}

func (r *RunCmd) Run(ctx *kong.Context) error {
	demo, err := LoadDemoFile(r.Filename)

	if err != nil {
		return err
	}

	TitleFormat.Println("Demo:", demo.Title)
	return RunDemo(demo.Steps, demo.Variables)
}

func (r *TeardownCmd) Run(ctx *kong.Context) error {
	demo, err := LoadDemoFile(r.Filename)

	if err != nil {
		return err
	}

	TitleFormat.Println("Teardown:", demo.Title)
	return RunDemo(demo.Teardown, demo.Variables)
}
