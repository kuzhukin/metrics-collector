package codec

import (
	"fmt"

	"github.com/kuzhukin/metrics-collector/internal/metric"
)

const EmptyHTML = `
<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8" />
<title></title>
<link rel="stylesheet" href="style.css" />
</head>
<body>
%s
</body>
</html>
`

func DecodeMetricsList(metrics []*metric.Metric) string {
	listHTML := ""
	listHTML += "\t<ul>\n"

	for _, m := range metrics {
		listHTML += fmt.Sprintf("\t\t<li>%s: %s</li>\n", m.Name, DecodeValue(m))
	}

	listHTML += "\t</ul>"

	return fmt.Sprintf(EmptyHTML, listHTML)
}
