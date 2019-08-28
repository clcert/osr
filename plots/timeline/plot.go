package timeline

import (
	"github.com/clcert/osr/plots"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"path"
	"strings"
	"time"
)

const DATEFIELD = "date"
const DATEFORMAT = "2006-01-02"

type Plot struct {
	*plots.PlotIO
	*plots.PlotProps
}

func (aPlot *Plot) GetPlotIO() *plots.PlotIO {
	return aPlot.PlotIO
}

func (aPlot *Plot) Plot() error {

	p, err := plot.New()
	if err != nil {
		return err
	}
	xys, err := aPlot.GetXYs(DATEFIELD, fDate, nil)
	if err != nil {
		return err
	}

	p.Title.Text = aPlot.Title + " - " + path.Base(aPlot.InPath)
	p.X.Label.Text = aPlot.XLabel
	p.Y.Label.Text = aPlot.YLabel
	p.Y.Min = 0
	p.X.Tick.Marker = plot.TimeTicks{
		Ticker: nil,
		Format: DATEFORMAT,
		Time:   nil,
	}

	p.Add(plotter.NewGrid())

	for _, data := range xys {
		err := plotutil.AddLinePoints(p, data)
		if err != nil {
			return err
		}
	}

	writerTo, err := p.WriterTo( 400, 300, aPlot.Format)
	if err != nil {
		return err
	}
	_, err = writerTo.WriteTo(aPlot.Out)
	if err != nil {
		return err
	}
	return nil
}

func fDate(s string) float64 {
	d, err := time.Parse(DATEFORMAT, s)
	if err != nil {
		// split the time
		sArr := strings.Split(s, " ")
		d, err = time.Parse(DATEFORMAT, sArr[0])
		if err != nil {
			return 0
		}
	}
	return float64(d.Unix())
}
