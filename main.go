package main

import (
	"github.com/jleverenz/durt/cmd"
)

func main() {
	cmd.Execute()
	// globalOpts = ProgramOptions{}

	// app := &cli.App{
	// 	Name:    "durt",
	// 	Usage:   "disk utilization reporting tool",
	// 	Version: "v0.0.0",
	// 	Flags: []cli.Flag{
	// 		// TODO it'd be nice to allow --head, --head 30, etc; seems this flag
	// 		// parsing module doesn't support that
	// 		&cli.BoolFlag{
	// 			Name:  "head",
	// 			Usage: "display the top 20 records",
	// 		},
	// 		&cli.StringSliceFlag{
	// 			Name:  "exclude",
	// 			Usage: "exclude paths by regex",
	// 		},
	// 	},
	// 	Action: func(cCtx *cli.Context) error {
	// 		globalOpts.head = cCtx.Bool("head")

	// 		exclusions := cCtx.StringSlice("exclude")
	// 		for _, exc := range exclusions {
	// 			globalOpts.exclusions = append(globalOpts.exclusions, regexp.MustCompile(exc))
	// 		}

	// 		mainAction(cCtx.Args().Slice())
	// 		return nil
	// 	},
	// }

	// if err := app.Run(os.Args); err != nil {
	// 	log.Fatal(err)
	// }
}
