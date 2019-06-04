package newrelic

import (
	"fmt"

	"github.com/wtfutil/wtf/wtf"
	nr "github.com/yfronto/newrelic"
)

func (widget *Widget) display() {
	client := widget.currentData()
	if client == nil {
		widget.Redraw(widget.CommonSettings.Title, " NewRelic data unavailable ", false)
		return
	}
	app, appErr := client.Application()
	deploys, depErr := client.Deployments()

	appName := "error"
	if appErr == nil {
		appName = app.Name
	}

	var content string
	title := fmt.Sprintf("%s - [green]%s[white]", widget.CommonSettings.Title, appName)
	wrap := false
	if depErr != nil {
		wrap = true
		content = depErr.Error()
	} else {
		content = widget.contentFrom(deploys)
	}

	widget.Redraw(title, content, wrap)
}

func (widget *Widget) contentFrom(deploys []nr.ApplicationDeployment) string {
	str := fmt.Sprintf(
		" %s\n",
		"[red]Latest Deploys[white]",
	)

	revisions := []string{}

	for _, deploy := range deploys {
		if (deploy.Revision != "") && wtf.Exclude(revisions, deploy.Revision) {
			lineColor := "white"
			if wtf.IsToday(deploy.Timestamp) {
				lineColor = "lightblue"
			}

			revLen := 8
			if revLen > len(deploy.Revision) {
				revLen = len(deploy.Revision)
			}

			str += fmt.Sprintf(
				" [green]%s[%s] %s %-.16s[white]\n",
				deploy.Revision[0:revLen],
				lineColor,
				deploy.Timestamp.Format("Jan 02 15:04 MST"),
				wtf.NameFromEmail(deploy.User),
			)

			revisions = append(revisions, deploy.Revision)

			if len(revisions) == widget.settings.deployCount {
				break
			}
		}
	}

	return str
}
