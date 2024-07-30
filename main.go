package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	app := tview.NewApplication()

	// Create a table to display the metrics
	table := tview.NewTable().SetBorders(true)

	// Update function to refresh the metrics
	go func() {
		for {
			// Fetch CPU Usage
			cpuPercent, err := cpu.Percent(0, false)
			if err != nil {
				log.Fatal(err)
			}

			// Fetch Memory Usage
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				log.Fatal(err)
			}

			// Fetch Disk Usage
			diskInfo, err := disk.Usage("/")
			if err != nil {
				log.Fatal(err)
			}

			// Update the table
			app.QueueUpdateDraw(func() {
				table.Clear()
				table.SetCell(0, 0, tview.NewTableCell("Metric").SetTextColor(tview.Styles.ContrastBackgroundColor))
				table.SetCell(0, 1, tview.NewTableCell("Usage").SetTextColor(tview.Styles.MoreContrastBackgroundColor))

				table.SetCell(1, 0, tview.NewTableCell("CPU Usage"))
				table.SetCell(1, 1, tview.NewTableCell(fmt.Sprintf("%.2f%%", cpuPercent[0])))

				table.SetCell(2, 0, tview.NewTableCell("Memory Usage"))
				table.SetCell(2, 1, tview.NewTableCell(fmt.Sprintf("%.2f%%", memInfo.UsedPercent)))

				table.SetCell(3, 0, tview.NewTableCell("Disk Usage"))
				table.SetCell(3, 1, tview.NewTableCell(fmt.Sprintf("%.2f%%", diskInfo.UsedPercent)))
			})

			time.Sleep(2 * time.Second) // Update every 2 seconds
		}
	}()

	// Set the table as the root primitive and run the application
	if err := app.SetRoot(table, true).Run(); err != nil {
		log.Fatal(err)
	}
}
