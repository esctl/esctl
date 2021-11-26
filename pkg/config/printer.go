package config

import "github.com/pterm/pterm"

func Print(cfg *ClusterConfig) {

	var panels pterm.Panels
	for i, v := range cfg.Clusters {
		panels = append(panels, []pterm.Panel{})
		panels[i] = append(panels[i], pterm.Panel{Data: pterm.DefaultBox.WithTitle(v.Name).Sprint(v.Hosts)})
	}

	panels = append(panels, []pterm.Panel{})
	panels[len(cfg.Clusters)] = append(panels[len(cfg.Clusters)], pterm.Panel{Data: pterm.DefaultBox.WithTitle(pterm.FgGreen.Sprint("Current Cluster")).Sprint(cfg.CurrentCluster)})

	panelstr, _ := pterm.DefaultPanel.WithPanels(panels).Srender()

	pterm.DefaultBox.WithTitle("Cluster Config").WithRightPadding(0).WithBottomPadding(0).Println(panelstr)
}
