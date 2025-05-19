package services

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/sqlite" // Using SQLite for the player's DB
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQLFile represents a parsed SQL file with its sequence number.
type SQLFile struct {
	Path     string
	Sequence int
	BaseName string // e.g., "init", "lorry_A_trip1"
}

// SimulateDataPlayer executes simulation data from SQL files.
type SimulateDataPlayer struct {
	dataPath    string
	productName string
	dateStr     string
	speedFactor float64
	db          *gorm.DB
	baseDate    time.Time // The date for which simulation is running, time part 00:00:00

	// Internal state
	sqlFiles       []SQLFile
	playedInit     bool
	playedUninit   bool
	currentFileIdx int
	// For advanced playback control, could add more state like last executed timestamp
}

// NewSimulateDataPlayer creates a new player instance.
func NewSimulateDataPlayer(dataPath, productName, dateStr string, speed float64) (*SimulateDataPlayer, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date string format for player: %w", err)
	}

	// In-memory SQLite for the player to execute SQL against
	// For actual database targets, GORM would connect to MySQL, Postgres, etc.
	// For this simulation, we are just "playing" by printing, but a real DB is better for verification.
	// Let's use a file-based SQLite DB for persistence if needed during a session or for inspection.
	dbPath := filepath.Join(dataPath, productName, dateStr, "simulation_player.db")
	os.Remove(dbPath) // Clean up previous run's DB

	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Or logger.Info for debugging
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to player's sqlite database: %w", err)
	}

	// Auto-migrate schema (if models are defined and player needs to know about them)
	// For now, player just executes raw SQL, so schema should be in 0_init.sql or handled by user.
	// Example: gormDB.AutoMigrate(&locationGnssDataModel.LocationGnssData{}, &weighLoggerModel.WeighLogger{}, &yAnalyserModel.YAnalyser{})

	return &SimulateDataPlayer{
		dataPath:    dataPath,
		productName: productName,
		dateStr:     dateStr,
		speedFactor: speed,
		db:          gormDB,
		baseDate:    parsedDate,
		sqlFiles:    []SQLFile{},
	}, nil
}

// LoadSQLFiles discovers and sorts SQL files to be played.
func (p *SimulateDataPlayer) LoadSQLFiles() error {
	dirPath := filepath.Join(p.dataPath, p.productName, p.dateStr)
	log.Printf("Player: Loading SQL files from %s", dirPath)

	p.sqlFiles = []SQLFile{} // Reset

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			parts := strings.SplitN(d.Name(), "_", 2)
			if len(parts) < 2 {
				log.Printf("Player: Skipping file with unexpected name format: %s", d.Name())
				return nil
			}
			seq, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Printf("Player: Skipping file with non-integer sequence: %s", d.Name())
				return nil
			}
			baseNameWithExt := parts[1]
			baseName := strings.TrimSuffix(baseNameWithExt, ".sql")

			p.sqlFiles = append(p.sqlFiles, SQLFile{Path: path, Sequence: seq, BaseName: baseName})
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory %s: %w", dirPath, err)
	}

	// Sort files by sequence number
	sort.Slice(p.sqlFiles, func(i, j int) bool {
		return p.sqlFiles[i].Sequence < p.sqlFiles[j].Sequence
	})

	if len(p.sqlFiles) == 0 {
		return fmt.Errorf("no SQL files found in %s", dirPath)
	}
	log.Printf("Player: Found %d SQL files to play.", len(p.sqlFiles))
	for _, sf := range p.sqlFiles {
		log.Printf("  - File: %s (Seq: %d)", filepath.Base(sf.Path), sf.Sequence)
	}
	return nil
}

// Play executes the loaded SQL files according to their timestamps and speed factor.
func (p *SimulateDataPlayer) Play() error {
	if err := p.LoadSQLFiles(); err != nil {
		return fmt.Errorf("failed to load SQL files for playing: %w", err)
	}

	// Execute 0_init.sql first, if it exists, at full speed.
	initFileIndex := -1
	for i, sf := range p.sqlFiles {
		if sf.Sequence == 0 && strings.Contains(sf.BaseName, "init") { // Or just sf.Sequence == 0
			initFileIndex = i
			break
		}
	}
	if initFileIndex != -1 {
		initFile := p.sqlFiles[initFileIndex]
		log.Printf("Player: Executing init file %s immediately.", filepath.Base(initFile.Path))
		if err := p.executeSQLFile(initFile, true, nil); err != nil { // isImmediate = true
			return fmt.Errorf("error executing init file %s: %w", initFile.Path, err)
		}
		p.playedInit = true
		// Remove init file from list or mark as played to avoid re-processing
		p.sqlFiles = append(p.sqlFiles[:initFileIndex], p.sqlFiles[initFileIndex+1:]...)
	} else {
		log.Println("Player: No 0_init.sql file found or already processed.")
	}

	// Main playback loop
	simulationStartTime := time.Now()
	var lastRealTimePlayed time.Time
	var lastSimTimePlayed time.Time

	log.Printf("Player: Starting main playback at %s with speed factor %.2f.", simulationStartTime.Format(TimeFormat), p.speedFactor)
	log.Printf("Player: Simulation date is %s.", p.baseDate.Format("2006-01-02"))

	for _, sqlFile := range p.sqlFiles {
		if sqlFile.Sequence == 100 { // Uninit file, handle separately at the end
			continue
		}
		if sqlFile.Sequence == 0 { // Already handled init file
			continue
		}

		log.Printf("Player: Processing file %s (Seq: %d)", filepath.Base(sqlFile.Path), sqlFile.Sequence)
		err := p.executeSQLFile(sqlFile, false, &playbackState{
			simulationStartTime: &simulationStartTime,
			lastRealTimePlayed:  &lastRealTimePlayed,
			lastSimTimePlayed:   &lastSimTimePlayed,
		})
		if err != nil {
			return fmt.Errorf("error playing file %s: %w", sqlFile.Path, err)
		}
	}

	// Execute 100_uninit.sql last, if it exists, at full speed.
	//uninitFileIndex := -1
	// Reload sqlFiles or search in original loaded list before modification if needed
	// For simplicity, assuming it's still in the (potentially modified) p.sqlFiles list or we re-scan for it
	// Let's refine: LoadSQLFiles should be called once. We iterate and skip.
	// So, need to re-scan original list or manage 'played' status.
	// Easiest is to iterate the original list and skip based on sequence for main play.
	// Then specifically look for 100_uninit.sql.

	// Re-scan or find uninit file (if not done by iterating the original full list)
	var uninitFile *SQLFile
	fullSqlFilesList, _ := p.getFullSqlFileList() // Get a fresh list to find uninit
	for i := range fullSqlFilesList {
		if fullSqlFilesList[i].Sequence == 100 {
			uninitFile = &fullSqlFilesList[i]
			break
		}
	}

	if uninitFile != nil && !p.playedUninit {
		log.Printf("Player: Executing uninit file %s immediately at the end.", filepath.Base(uninitFile.Path))
		if err := p.executeSQLFile(*uninitFile, true, nil); err != nil { // isImmediate = true
			return fmt.Errorf("error executing uninit file %s: %w", uninitFile.Path, err)
		}
		p.playedUninit = true
	} else {
		log.Println("Player: No 100_uninit.sql file found or already processed.")
	}

	log.Println("Player: All SQL files played.")
	return nil
}

// playbackState holds the timing state for paced execution.
type playbackState struct {
	simulationStartTime *time.Time // Wall clock time when simulation playback started
	lastRealTimePlayed  *time.Time // Wall clock time when the last SQLs for a sim timestamp were played
	lastSimTimePlayed   *time.Time // The simulation timestamp of the last played SQLs
}

// executeSQLFile reads a single SQL file, parses comments for time, and executes statements.
// isImmediate: if true, executes all SQLs without time-based pacing.
// state: used for paced execution.
func (p *SimulateDataPlayer) executeSQLFile(sqlFile SQLFile, isImmediate bool, state *playbackState) error {
	file, err := os.Open(sqlFile.Path)
	if err != nil {
		return fmt.Errorf("failed to open SQL file %s: %w", sqlFile.Path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentExecutionTime time.Time
	var sqlBuffer []string

	firstSimTimeInFile := true // To handle the "fast forward" rule

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "--- ") { // Time annotation
			// Execute buffered SQLs for the previous timestamp (if any)
			if len(sqlBuffer) > 0 && !currentExecutionTime.IsZero() {
				if err := p.executeSQLBatch(sqlBuffer, currentExecutionTime, isImmediate, state); err != nil {
					return err
				}
				sqlBuffer = []string{} // Reset buffer
			}

			timeStr := strings.TrimPrefix(line, "--- ")
			parsedTime, err := ParseTimeAnnotation(p.baseDate, timeStr) // Use baseDate
			if err != nil {
				log.Printf("Player: Error parsing time annotation '%s' in %s: %v. Skipping block.", line, sqlFile.Path, err)
				// Skip lines until next time annotation or EOF
				currentExecutionTime = time.Time{} // Invalidate current time
				continue
			}
			currentExecutionTime = parsedTime
			// log.Printf("Player: Encountered time %s in %s", currentExecutionTime.Format(TimeFormat), filepath.Base(sqlFile.Path))

			// Rule: "例如现在是13:00 开始执行sql,速度为1,那么要快速执行完13:00前的所有语句"
			// This applies if `isImmediate` is false.
			// If this is the first relevant sim time we're seeing (after init),
			// and it's before the real "current time" derived from simulationStartTime,
			// it should be fast-forwarded.
			if !isImmediate && state != nil && state.simulationStartTime != nil {
				if firstSimTimeInFile {
					// All statements in the file up to the *current wall clock equivalent sim time*
					// should be executed quickly.
					// Let currentSimTimeFromWallClock = p.baseDate + (time.Now() - *state.simulationStartTime)
					// For simplicity here, if the file's first timestamp is in the past relative
					// to simulation start, or if it's an early part of the day, we assume it's "catch up".
					// The rule is a bit ambiguous if it applies per file or per block.
					// Let's assume it means: if wall clock is 13:00, any SQL for sim time < 13:00 is fast.

					// A simpler interpretation for now: if the currentExecutionTime is "before"
					// the state.lastSimTimePlayed (if any, considering speed factor), then it's catch-up.
					// Or, more directly per rule: If current wall clock means we are at simulated 13:05,
					// then any SQL for 13:00, 13:01, etc. up to 13:05 should be "fast".
					// The logic in executeSQLBatch handles the actual waiting.
					// Here, we just mark that it's no longer the "first" sim time for pacing.
					firstSimTimeInFile = false
				}
			}

		} else if !currentExecutionTime.IsZero() { // Regular SQL statement for the current time block
			sqlBuffer = append(sqlBuffer, line)
		} else if isImmediate && strings.HasSuffix(sqlFile.BaseName, "init") || strings.HasSuffix(sqlFile.BaseName, "uninit") {
			// For init/uninit files, if there's no leading time comment, execute immediately.
			sqlBuffer = append(sqlBuffer, line)
		}
	}

	// Execute any remaining SQLs in the buffer (for the last timestamp block)
	if len(sqlBuffer) > 0 && (!currentExecutionTime.IsZero() || isImmediate) {
		execTime := currentExecutionTime
		if isImmediate && execTime.IsZero() { // For init/uninit files that might not have time comments
			execTime = p.baseDate // Arbitrary time, won't be used for pacing
		}
		if err := p.executeSQLBatch(sqlBuffer, execTime, isImmediate, state); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning SQL file %s: %w", sqlFile.Path, err)
	}
	return nil
}

// executeSQLBatch executes a batch of SQL statements for a given simulation time.
// It handles pacing based on speedFactor.
func (p *SimulateDataPlayer) executeSQLBatch(sqls []string, simTimeTarget time.Time, isImmediate bool, state *playbackState) error {
	if len(sqls) == 0 {
		return nil
	}

	if !isImmediate && state != nil && state.simulationStartTime != nil && p.speedFactor > 0 {
		// Paced execution
		now := time.Now()
		if state.lastSimTimePlayed.IsZero() { // First paced event
			// If simTimeTarget is before simulationStartTime (wall clock), execute immediately (catch-up)
			// The baseDate + time part of simTimeTarget gives the full simulation datetime.
			// simulationStartTime is wall clock time.
			// Effective "current" simulation time based on wall clock:
			// p.baseDate.Add(now.Sub(*state.simulationStartTime))
			// This is a bit complex. Let's simplify:
			// The first ever timed event establishes the baseline.
			// All events before current wall-clock-equivalent sim time are fast.

			// If simTimeTarget is in the "past" relative to when we started playing, execute fast.
			// "Past" means simTimeTarget < (baseDate + (now - simulationStartTime))
			simTimeFromWallClockStart := p.baseDate.Add(now.Sub(*state.simulationStartTime))
			if simTimeTarget.Before(simTimeFromWallClockStart) {
				//This is catch-up execution as per rule.
				log.Printf("Player: Catch-up execution for sim time %s (current wall-equiv sim time: %s)", simTimeTarget.Format(DateTimeFormat), simTimeFromWallClockStart.Format(DateTimeFormat))
			} else {
				// This is the first "future" event. Wait until wall clock catches up.
				// Calculate real time duration from simulation start to this event's sim time
				simDurationSinceStart := simTimeTarget.Sub(p.baseDate.Add(state.simulationStartTime.Sub(p.baseDate))) // Target sim time relative to sim day start at wall time
				if simTimeTarget.Hour() < state.simulationStartTime.Hour() && state.simulationStartTime.Hour() >= 6 { // target time is before 6am, but we started at/after 6am
					// This condition is for initial catch-up of events before the "start of play" time (e.g. 13:00)
					simDurationSinceStart = 0 // effectively execute immediately
				}

				realDurationToWait := time.Duration(float64(simDurationSinceStart) / p.speedFactor)
				targetRealTime := state.simulationStartTime.Add(realDurationToWait)

				if now.Before(targetRealTime) {
					waitNeeded := targetRealTime.Sub(now)
					log.Printf("Player: Waiting for %.2fs (real time) to reach sim time %s", waitNeeded.Seconds(), simTimeTarget.Format(DateTimeFormat))
					time.Sleep(waitNeeded)
				}
			}
			*state.lastSimTimePlayed = simTimeTarget
			*state.lastRealTimePlayed = time.Now()

		} else {
			// Subsequent paced events
			simTimeDiff := simTimeTarget.Sub(*state.lastSimTimePlayed)
			if simTimeDiff < 0 { // Should not happen if files are ordered and processed correctly
				log.Printf("Player: Warning - sim time %s is before last played sim time %s. Executing immediately.", simTimeTarget.Format(DateTimeFormat), state.lastSimTimePlayed.Format(DateTimeFormat))
				simTimeDiff = 0
			}

			realTimeWaitDuration := time.Duration(float64(simTimeDiff) / p.speedFactor)
			expectedRealTime := state.lastRealTimePlayed.Add(realTimeWaitDuration)

			if now.Before(expectedRealTime) {
				waitNeeded := expectedRealTime.Sub(now)
				// log.Printf("Player: Waiting for %.2fs (real time) to reach sim time %s (prev sim: %s, prev real: %s)",
				// 	waitNeeded.Seconds(), simTimeTarget.Format(TimeFormat), state.lastSimTimePlayed.Format(TimeFormat), state.lastRealTimePlayed.Format(TimeFormat))
				time.Sleep(waitNeeded)
			} else {
				// We are late, execute immediately (catch-up for this specific event)
				// log.Printf("Player: Catching up for sim time %s (expected by %s, now %s)", simTimeTarget.Format(TimeFormat), expectedRealTime.Format(TimeFormat), now.Format(TimeFormat))
			}
			*state.lastSimTimePlayed = simTimeTarget
			*state.lastRealTimePlayed = time.Now()
		}
	}
	// If isImmediate or speedFactor <= 0, no waiting.

	// Execute SQLs in a transaction for this time block
	tx := p.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction for sim time %s: %w", simTimeTarget.Format(DateTimeFormat), tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-panic
		}
	}()

	//log.Printf("Player: Executing %d SQL statement(s) for sim time %s", len(sqls), simTimeTarget.Format(DateTimeFormat))
	for i, sqlCmd := range sqls {
		// log.Printf("  SQL: %s", sqlCmd)
		if err := tx.Exec(sqlCmd).Error; err != nil {
			tx.Rollback()
			// Provide more context on error
			return fmt.Errorf("error executing SQL (sim time %s, statement %d: '%s'): %w", simTimeTarget.Format(DateTimeFormat), i+1, sqlCmd, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		// Rollback already implicitly handled by defer if commit fails, but explicit can be added.
		return fmt.Errorf("failed to commit transaction for sim time %s: %w", simTimeTarget.Format(DateTimeFormat), err.Error)
	}
	return nil
}

// getFullSqlFileList is a helper to get the original sorted list of SQL files.
// This is useful if p.sqlFiles was modified during processing (e.g., removing init file).
func (p *SimulateDataPlayer) getFullSqlFileList() ([]SQLFile, error) {
	originalList := []SQLFile{}
	dirPath := filepath.Join(p.dataPath, p.productName, p.dateStr)
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			parts := strings.SplitN(d.Name(), "_", 2)
			if len(parts) < 2 {
				return nil
			}
			seq, convErr := strconv.Atoi(parts[0])
			if convErr != nil {
				return nil
			}
			baseNameWithExt := parts[1]
			baseName := strings.TrimSuffix(baseNameWithExt, ".sql")
			originalList = append(originalList, SQLFile{Path: path, Sequence: seq, BaseName: baseName})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(originalList, func(i, j int) bool { return originalList[i].Sequence < originalList[j].Sequence })
	return originalList, nil
}

// playSimulateData function as requested (entry point)
func PlaySimulateData(dataPath, productName, dateStr string, speed float64) {
	player, err := NewSimulateDataPlayer(dataPath, productName, dateStr, speed)
	if err != nil {
		log.Fatalf("Failed to create SimulateDataPlayer: %v", err)
	}
	err = player.Play()
	if err != nil {
		log.Fatalf("Failed to play simulation data: %v", err)
	}
}
