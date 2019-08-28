package plot

import (
	"github.com/spf13/cobra"
)

var title, xLabel, yLabel, inFolder, outFolder, ext string

func init() {
	PlotCmd.PersistentFlags().StringVarP(&title, "title", "t", "Time Plot", "Graph title.")
	PlotCmd.PersistentFlags().StringVarP(&xLabel, "x-label", "x", "", "Graph X Label.")
	PlotCmd.PersistentFlags().StringVarP(&yLabel, "y-label", "y", "", "Graph Y Label.")
	PlotCmd.PersistentFlags().StringVarP(&inFolder, "input-folder", "i", "", "Input folder absolute path.")
	PlotCmd.PersistentFlags().StringVarP(&outFolder, "output-folder", "o", "./", "Output folder absolute path.")
	PlotCmd.PersistentFlags().StringVarP(&ext, "format", "f", "png", "Output format.")
	_ = PlotCmd.MarkPersistentFlagRequired("input-folder")
	_ = PlotCmd.MarkPersistentFlagRequired("output-folder")
}

// The root of osr.
var PlotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Quick plotting utility",
}
