package story

import (
	"fmt"

	"github.com/fatih/color"
)

func (s *Story) StartCLI(init string) error {
	if _, ok := (*s)[init]; !ok {
		return fmt.Errorf("could not find chapter %q", init)
	}

	tf := color.New(color.Bold)
	fmt.Print("Info: Input '0' anytime to quit\n\n")

	name := init
	for {
		arc := (*s)[name]

		tf.Println(arc.Title)
		for _, p := range arc.Story {
			fmt.Printf("%v\n\n", p)
		}

		fmt.Println("Options (choose with the number):")
		if len(arc.Options) == 0 {
			fmt.Println("Your adventure has ended.")
			break
		}

		sp := len(string(len(arc.Options)))
		for {
			for i, op := range arc.Options {
				fmt.Printf("%*d. %v\n", sp, i+1, op.Text)
			}

			fmt.Print("? ")
			var sel int
			if _, err := fmt.Scan(&sel); err != nil {
				fmt.Println("Invalid option: Not an integer")
				continue
			}
			sel--
			if sel < 0 || sel >= len(arc.Options) {
				fmt.Println("Invalid option: Out of range")
				continue
			}

			name = arc.Options[sel].Arc
			fmt.Println()
			break
		}
	}

	return nil
}
