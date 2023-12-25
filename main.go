package main

import (
	"csguard/internal/calculate"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	provider := calculate.NewChecksumProvider()

	app := &cli.App{
		Name: "csguard",
		Commands: []*cli.Command{
			{
				Name:  "calculate",
				Usage: "calculate checksum",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input-file",
						Value:       "",
						Usage:       "Path to the input file for which to calculate the checksum.",
						Destination: provider.SetInputFile(),
					},
					&cli.StringFlag{
						Name:        "input-folder",
						Value:       "",
						Usage:       "Path to the input folder containing files for which to calculate the checksum.",
						Destination: provider.SetInputFolder(),
					},
					&cli.StringFlag{
						Name:        "output",
						Value:       "table",
						Usage:       "Format/Path to the output file where the calculated checksums will be written. Supported formats: '.txt', '.json', '.yaml', 'table'. Default will be 'table'",
						Destination: provider.SetOutputFile(),
					},
					&cli.StringFlag{
						Name:        "algorithm",
						Value:       "md5",
						Usage:       "Hashing algorithm for checksum calculation. (options: 'md5', 'sha256', 'sha512', 'crc'; default: md5)",
						Destination: provider.SetAlgorithm(),
					},
				},
				Action: func(cCtx *cli.Context) error {
					if err := provider.CalculateInputValidation(); err != nil {
						return cli.Exit(err.Error(), 1)
					}
					if err := provider.CalculateChecksum(); err != nil {
						return cli.Exit(err.Error(), 1)
					}
					if err := provider.CreateCalculateOutput(); err != nil {
						return cli.Exit(err.Error(), 1)
					}
					return nil
				},
			},
			{
				Name:  "validate",
				Usage: "validate checksum",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input-file",
						Value:       "",
						Usage:       "Path to the input file for which to validate the checksum.",
						Destination: provider.SetInputFile(),
					},
					&cli.StringFlag{
						Name:  "checksum",
						Value: "",
						Usage: "Checksum value to validate against the calculated checksum of the input file.",
					},
					&cli.StringFlag{
						Name:        "checksum-file",
						Value:       "",
						Usage:       "Path to the file containing checksums to validate against the input file. Supported formats: '.txt', '.json', '.yaml'.",
						Destination: provider.SetChecksumFolder(),
					},
					&cli.StringFlag{
						Name:        "output",
						Value:       "table",
						Usage:       "Format/Path to the output file where the validation result will be written. Supported formats: '.txt', '.json', '.yaml', 'table'. Default will be 'table'",
						Destination: provider.SetOutputFile(),
					},
					&cli.StringFlag{
						Name:        "algorithm",
						Value:       "md5",
						Usage:       "Hashing algorithm for checksum calculation (options: 'md5', 'sha256', 'sha512', 'crc'; default: md5)",
						Destination: provider.SetAlgorithm(),
					},
				},
				Action: func(cCtx *cli.Context) error {
					if err := provider.ValidateInputValidation(cCtx.String("checksum")); err != nil {
						return cli.Exit(err.Error(), 1)
					}
					if err := provider.ValidateChecksum(); err != nil {
						return cli.Exit(err.Error(), 1)
					}
					if err := provider.CreateValidateOutput(); err != nil {
						return cli.Exit(err.Error(), 1)
					}
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		cli.Exit(err.Error(), 1)
	}
}
