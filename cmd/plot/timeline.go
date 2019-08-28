package plot

import (
	"github.com/clcert/osr/plots"
	"github.com/clcert/osr/plots/timeline"
	"github.com/spf13/cobra"
)

var dateColumn string

func init() {
	TimelineCmd.Flags().StringVarP(&dateColumn, "date-column", "d", "date", "date column name")
	PlotCmd.AddCommand(TimelineCmd)
}

var TimelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Timeline graph",
	RunE: func(cmd *cobra.Command, args []string) error {

		plotIOs, err := plots.GetLocalIO(inFolder, outFolder, ext)
		if err != nil {
			return err
		}

		plotProps := &plots.PlotProps{
			Title:  title,
			XLabel: xLabel,
			YLabel: yLabel,
			Format: ext,
		}

		for _, plotIO := range plotIOs {
			plot := &timeline.Plot{
				PlotIO: plotIO,
				PlotProps: plotProps,
			}
			err := plot.Plot()
			plot.In.Close()
			plot.Out.Close()
			if err != nil {
				return err
			}
		}

		return nil
	},
}
