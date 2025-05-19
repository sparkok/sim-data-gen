package main

import (
	"flag"
	"log"
	"os"
	"sim_data_gen/services"
	"time"
)

func main() {
	cmd := flag.String("cmd", "", "Command: genSimulateData or play")
	dataPath := flag.String("dataPath", "data", "Path to data directory")
	productName := flag.String("productName", "CP1", "Product name (e.g., CP1)")
	dateStr := flag.String("dateStr", time.Now().Format("2006-01-02"), "Date string (YYYY-MM-DD)")
	speed := flag.Float64("speed", 1.0, "Playback speed factor for play command")
	flag.Parse()
	log.SetOutput(os.Stdout) // Ensure logs go to stdout
	log.SetFlags(log.Ltime | log.Lshortfile)

	switch *cmd {
	case "gen":
		log.Printf("Executing command: genSimulateData with dataPath=%s, productName=%s, dateStr=%s\n", *dataPath, *productName, *dateStr)
		services.GenSimulateData(*dataPath, *productName, *dateStr)
		log.Println("genSimulateData finished.")
	case "play":
		log.Printf("Executing command: play with dataPath=%s, productName=%s, dateStr=%s, speed=%.2f\n", *dataPath, *productName, *dateStr, *speed)
		services.PlaySimulateData(*dataPath, *productName, *dateStr, *speed)
		log.Println("play finished.")
	default:
		log.Fatalf("Unknown command: %s. Use 'genSimulateData' or 'play'.\n", *cmd)
	}
}
