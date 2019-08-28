package plots

import (
	"fmt"
	"github.com/clcert/osr/utils"
	"gonum.org/v1/plot/plotter"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const CSV = ".csv"

type Plotter interface{
	GetPlotIO() *PlotIO
	Plot() error
}

type PlotProps struct {
	Title string
	XLabel string
	YLabel string
	Colors []color.Color
	Format string
}

type PlotIO struct {
	InPath  string
	OutPath string
	In      io.ReadCloser
	Out     io.WriteCloser
}

func parseDefault(s string) float64 {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return i
}

// Returns a map of plotter.XYs, using xLabel as X label and transforming the strings from
// the MapArray with fx and fy functions. If fx or fy are nil, it tries to transform them using Atoi.
func (plotIO *PlotIO) GetXYs(xLabel string, fx func(string) float64, fy func(string) float64) (map[string]plotter.XYs, error) {

	if fx == nil {
		fx = parseDefault
	}
	if fy == nil {
		fy = parseDefault
	}

	xys := make(map[string]plotter.XYs)

	csvReader, err := utils.NewHeadedCSV(plotIO.In, nil)
	if err != nil {
		return nil,  err
	}
	mapArray, err := csvReader.ToMapArray()
	if err != nil {
		return nil, err
	}

	hasXLabel := false
	for _, label := range csvReader.Headers {
		if label == xLabel {
			hasXLabel = true
			continue
		}
		xys[label] = make(plotter.XYs, len(mapArray))
	}

	if !hasXLabel {
		return nil, fmt.Errorf("x label not found on input file")
	}

	// For each data point
	for i, mapPoint := range mapArray {
		// For each dimension
		xVal := fx(mapPoint[xLabel])
		for label, value := range mapPoint {
			if label == xLabel {
				continue
			}
			xys[label][i].X = xVal
			xys[label][i].Y = fy(value)
		}
	}
	return xys, nil
}


// GetLocalIO returns a list of PlotIO objects from the inFolder, and its desired plots
// destinations. The files to plots must be CSVs (and end in CSV)
func GetLocalIO(inFolder, outFolder, extension string) ([]*PlotIO, error) {
	plotIO := make([]*PlotIO, 0)
	if err := filepath.Walk(inFolder, func(inPath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(strings.ToLower(inPath), CSV) && !info.IsDir() {
			name := strings.TrimSuffix(strings.TrimPrefix(inPath, inFolder), CSV) + "." + extension
			outPath := filepath.Join(outFolder, name)
			if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
				return err
			}
			reader, err := os.Open(inPath)
			if err != nil {
				return err
			}
			writer, err := os.Create(outPath)
			plotIO = append(plotIO, &PlotIO{
				InPath:  inPath,
				OutPath: outPath,
				In:      reader,
				Out:     writer,
			})
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return plotIO, nil
}
