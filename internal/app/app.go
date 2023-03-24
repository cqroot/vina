package app

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cqroot/edname/internal/ediff"
	"github.com/cqroot/edname/internal/renamer"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type App struct{}

func Run(editor string, path string, dirOpt bool, dirOnlyOpt bool, allOpt bool) error {
	r, err := renamer.New(renamer.RenameOpt{
		WorkDir:              path,
		ShouldRenameDir:      dirOpt,
		ShouldOnlyRenameDir:  dirOnlyOpt,
		ShouldRenameDotFiles: allOpt,
	})
	if err != nil {
		return err
	}

	items, err := r.Generate()
	if err != nil {
		return err
	}

	ed := ediff.New(editor)
	ed.AppendItems(items)
	pairs, err := ed.Run()
	if err != nil {
		return err
	}

	if len(pairs) == 0 {
		return nil
	}

	PrintPairs(pairs)

	fmt.Print("Confirm to rename the above file [y/N] ")
	cfmReader := bufio.NewReader(os.Stdin)
	cfmText, err := cfmReader.ReadString('\n')
	if err != nil {
		return err
	}

	if cfmText != "y\n" && cfmText != "Y\n" {
		return nil
	}
	fmt.Println()

	for _, pair := range pairs {
		err := r.Rename(pair.Prev, pair.Curr)
		if err != nil {
			return err
		}

		fmt.Printf(
			"%s %s %s\n",
			pair.Prev,
			text.FgGreen.Sprint("->"),
			pair.Curr,
		)
	}

	return err
}

func PrintPairs(pairs []ediff.DiffPair) {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Old Name", "New Name"})
	for idx, pair := range pairs {
		t.AppendRow(table.Row{idx, pair.Prev, pair.Curr})
	}
	t.Render()
}
