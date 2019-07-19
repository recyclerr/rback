package main

import (
	"bytes"
	"text/template"

	"github.com/emicklei/dot"

	"github.com/team-soteria/rback/third-party/svgpanzoom"
	"github.com/team-soteria/rback/third-party/vizjs"
)

const webViewTemplate = `<html>
<head>
<title>rback</title>
</head>
<body>
<script>{{template "svgPanZoom" .}}</script>
<script>{{template "vizJs" .}}</script>
<script>{{template "vizJsRender" .}}</script>
<script>
    const graph = ` + "`" + `{{ .Graph }}` + "`" + `;
    let viz = new Viz();
    viz.renderSVGElement(graph).then(function (element) {
        element.setAttribute("id", "svg");
        element.setAttribute("width", "100%");
        element.setAttribute("height", "100%");
        document.body.appendChild(element);
        svgPanZoom('#svg', {
            zoomEnabled: true,
            controlIconsEnabled: true,
            fit: true,
            center: true,
            minZoom: 0.8
        });
    }).catch(error => {
        viz = new Viz();
        console.error(error);
    });
</script>
</body>
</html>
`

type webViewData struct {
	Graph string
}

func generateWebView(graph *dot.Graph) (*bytes.Buffer, error) {
	tpl, _ := template.New("vizJsTemplate").Parse(webViewTemplate)

	template.Must(tpl.Parse(`{{define "svgPanZoom"}}` + svgpanzoom.JsSource + `{{end}}`))
	template.Must(tpl.Parse(`{{define "vizJs"}}` + vizjs.JsSource + `{{end}}`))
	template.Must(tpl.Parse(`{{define "vizJsRender"}}` + "{{`" + vizjs.JsRenderSource + "`}}{{end}}"))

	buf := new(bytes.Buffer)
	err := tpl.Execute(buf, webViewData{Graph: graph.String()})

	return buf, err
}
