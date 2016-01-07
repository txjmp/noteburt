package schedule

import (
	"fmt"
	"lib"
	"log"
	"nb/data"
	"time"
)

const (
	atTime     = 0
	onInterval = 1
)

type item struct {
	name     string
	whenType int // atTime or onInterval
	when     int // atTime(hr 0-23) or onInterval(minutes)
	lastRun  time.Time
	run      func()
}

func Scheduler() {
	schedule := make([]*item, 0, 10)

	schedule = append(schedule, &item{
		name:     "backup database",
		whenType: atTime,
		when:     3, // run @ 3 am
		run:      backupDatabase,
	})
	/*
		schedule = append(schedule, &item{
			name:     "audit",
			whenType: atTime,
			when:     4, // run @ 4 am
			run:      audit,
		})
		schedule = append(schedule, &item{
			name:     "remove old logins",
			whenType: onInterval,
			when:     30, // run every 30 minutes
			run:      removeOldLogins,
		})
	*/
	for {
		time.Sleep(10 * time.Minute)
		now := time.Now()
		lib.Trace(0, "--- Scheduled Tasks ---")
		for _, task := range schedule {
			switch task.whenType {
			case atTime:
				if now.Hour() >= task.when && now.Day() != task.lastRun.Day() {
					lib.Trace(0, "running ", task.name)
					task.lastRun = time.Now()
					task.run()
				}
			case onInterval:
				minsSinceLastRun := now.Sub(task.lastRun) / time.Minute
				if minsSinceLastRun >= time.Duration(task.when) {
					lib.Trace(0, "running ", task.name)
					task.lastRun = time.Now()
					task.run()
				}
			default:
				log.Fatal("invalid schedule whenType - ", task.whenType)
			}
		}
	}
}
func backupDatabase() {
	dt := time.Now()
	_, mth, day := dt.Date()
	mthName := mth.String()[0:3] // Jan, Feb, Mar, ...
	hr, min, _ := dt.Clock()
	fileName := fmt.Sprintf("%v%d_%02d%02d.db", mthName, day, hr, min) // Jan01_1455.db
	data.DBBkup("db/bkup/" + fileName)
}

/*
func removeOldLogins() {
	resultChan := data.Data("removeOldLogins", data.NoParms{})
	<-resultChan
}
func audit() {
	resultChan := data.Data("audit", data.NoParms{})
	result := <-resultChan
	stats := result.Val.(*common.AuditData)
	Trace(0, "=========== A U D I T   S T A T S ============")
	Trace(0, "Books = ", stats.BookCnt)
	Trace(0, "Tabs = ", stats.TabCnt)
	Trace(0, "Tabs Loaded = ", stats.TabsLoaded)
	Trace(0, "Notes Loaded = ", stats.NotesLoaded)
	Trace(0, "=========== E N D   S T A T S ============")
}
*/
