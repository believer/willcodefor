package utils

import (
	"cmp"
	"slices"
	"strconv"
)

type CountData struct {
	Date  string `db:"date"`
	Label string `db:"label"`
	Count int    `db:"count"`
}

type Bar struct {
	Label     string
	Value     int
	BarHeight int
	BarWidth  int
	BarX      int
	BarY      int
	LabelX    float64
	LabelY    float64
	ValueX    float64
	ValueY    int
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func BarChart(data []CountData) ([]Bar, error) {
	var graphData []Bar

	graphHeight := 200
	graphWidth := 900
	maxCount := calculateMaxCount(data)

	// The data is used for a bar chart, so we need to convert the data
	for i, row := range data {
		var (
			elementsInGraph = graphWidth / len(data)
			// Calcualte the bar Height
			// Subtract 46 from the graph height to make room for the labels
			barHeight = clamp(int(float64(row.Count)/float64(maxCount)*float64(graphHeight-46)), 2, graphHeight-46)
			barWidth  = int(elementsInGraph) - 5

			// Space the bars evenly across the graph
			// Plus one px for border of first bar
			barX = elementsInGraph*i + 1
			barY = graphHeight - barHeight - 26
		)

		if barWidth <= 0 {
			barWidth = elementsInGraph
			barX = barX + 20
		}

		// Position centered on the bar. Subtract 3.4 which is half the width of the text.
		charWidth := 8.03 // Uses tabular nums so all characters are the same width
		numberOfCharsInCount := len(strconv.Itoa(row.Count))
		numberOfCharsInLabel := len(row.Label)

		halfWidthOfCount := charWidth * float64(numberOfCharsInCount) / 2
		halfWidthOfLabel := charWidth * float64(numberOfCharsInLabel) / 2

		valueX := float64(barX+(barWidth/2)) - halfWidthOfCount
		labelX := float64(barX+(barWidth/2)) - halfWidthOfLabel

		// If it's the first bar, we want to position the label at the start of the graph
		if i == 0 {
			labelX = float64(barX)
		}

		// If it's the last bar, we want to position the label at the end of the graph
		if i == len(data)-1 {
			labelX = float64(barX+barWidth) - charWidth*float64(numberOfCharsInLabel)
		}

		// Subtract 8 to put some space between the text and the bar
		valueY := barY - 8
		// 16,5 is the height of the text
		labelY := float64(barY) + float64(barHeight) + 20

		// Add the data to the graphData slice
		graphData = append(graphData, Bar{
			Label:     row.Label,
			Value:     row.Count,
			BarHeight: barHeight,
			BarWidth:  barWidth,
			BarX:      barX,
			BarY:      barY,
			ValueX:    valueX,
			ValueY:    valueY,
			LabelX:    labelX,
			LabelY:    labelY,
		})
	}

	return graphData, nil
}

type GridLine struct {
	Y1    int
	Y2    int
	Label int
}

type LineGraph struct {
	D     string
	YGrid []GridLine
}

func LineChart(data []CountData) (LineGraph, error) {
	graphHeight := 200
	graphWidth := 900
	maxCount := calculateMaxCount(data)
	var yGrid []GridLine

	// Start the path at the bottom left corner
	path := "M 0 " + strconv.Itoa(graphHeight)

	for i, row := range data {
		// Calculate the x and y values for the line
		x := float64(graphWidth) / float64(len(data)) * float64(i)
		y := float64(graphHeight) - float64(row.Count)/float64(maxCount)*float64(graphHeight)

		// Add point to the path
		path += " L " + strconv.FormatFloat(x, 'f', 3, 64) + " " + strconv.FormatFloat(y, 'f', 3, 64)
	}

	spacing := (graphHeight - 20) / 3

	for i := range 3 {
		ii := i + 1

		yGrid = append(yGrid, GridLine{
			Y1:    graphHeight - spacing*ii,
			Y2:    graphHeight - spacing*ii,
			Label: maxCount / 3 * ii,
		})
	}

	return LineGraph{
		D:     path,
		YGrid: yGrid,
	}, nil
}

func calculateMaxCount(data []CountData) int {
	c := slices.MaxFunc(data, func(a, b CountData) int {
		return cmp.Compare(a.Count, b.Count)
	})

	return c.Count
}
