package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/emicklei/dot"

	"github.com/team-soteria/rback/internal/svgpanzoom"
)

const webViewTemplate = `<html>
<head>
<title>rback</title>
</head>
<body>
{{ .Svg }}
<script>
{{ .SvgPanZoomSrc }}
</script>
<script>
svgPanZoom('#svg', {
	zoomEnabled: true,
	controlIconsEnabled: true,
	fit: true,
	center: true,
	minZoom: 0.1
});
</script>
</body>
</html>
`

type webViewData struct {
	Svg           string
	SvgPanZoomSrc string
}

func generateWebView(graph *dot.Graph) (*bytes.Buffer, error) {
	svgBuffer := new(bytes.Buffer)
	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(graph.String()), svgBuffer, os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute dot. Is Graphviz installed? Error: %v", err)
	}

	svg := svgBuffer.String()

	// Work around for dot bug which misses quoting some ampersands,
	// resulting on unparsable SVG.
	svg = strings.Replace(svg, "&;", "&amp;;", -1)

	// Drop the stuff before the <svg> start, and set the size to 100%.
	viewBox := regexp.MustCompile(`<svg\s*width="[^"]+"\s*height="[^"]+"\s*viewBox="[^"]+"`)
	if loc := viewBox.FindStringIndex(svg); loc != nil {
		svg = `<svg id="svg" width="100%" height="100%"` + svg[loc[1]:]
	}

	htmlBuffer := new(bytes.Buffer)
	tpl, _ := template.New("webViewTemplate").Parse(webViewTemplate)
	err := tpl.Execute(htmlBuffer, webViewData{
		Svg:           svg,
		SvgPanZoomSrc: svgpanzoom.JSSource,
	})

	return htmlBuffer, err
}
