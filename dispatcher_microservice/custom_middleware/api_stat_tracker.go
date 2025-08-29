package custommiddleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type DefaultApiStatTracker struct {
	Stats        *orderedmap.OrderedMap[string, map[string]int]
	mu           sync.Mutex
	Ctx          context.Context
	MaxLen       uint64
	Period       int64
	StartTime    int64
	currentKey   int64
	ApiStatsFile string
	DumpStats    bool
}

func (a *DefaultApiStatTracker) clearUpOldData() {
	for a.Stats.Len() > int(a.MaxLen) {
		a.Stats.Delete(a.Stats.Oldest().Key)
	}
}

// Starts this instance of api stat tracker. Blocks until the passed context is cancelled or expires.
func (a *DefaultApiStatTracker) StartTracker(ctx context.Context) {
	a.Stats = orderedmap.New[string, map[string]int]()
	a.StartTime = time.Now().UnixMilli()
	a.currentKey = a.StartTime - (a.StartTime % a.Period)
	a.makeEmptyStatRecordForCurrentPeriod()
	fmt.Println("Started tracking stats")
	var currentTime int64
	var tempKey int64
	for ctx.Err() == nil {
		currentTime = time.Now().UnixMilli()
		tempKey = currentTime - (currentTime % a.Period)
		if a.currentKey != tempKey {
			a.mu.Lock()
			a.currentKey = tempKey
			a.makeEmptyStatRecordForCurrentPeriod()
			a.clearUpOldData()
			a.mu.Unlock()
		}
	}
	fmt.Println("Tracker has successfully stopped")
	a.dumpLogs()
}

func (a *DefaultApiStatTracker) makeEmptyStatRecordForCurrentPeriod() {
	mapKey := time.UnixMilli(a.currentKey).Format("2006-01-02T15:04:05.000Z07:00")
	a.Stats.Set(mapKey, make(map[string]int))
}

func (a *DefaultApiStatTracker) dumpLogs() {
	if !a.DumpStats {
		return
	}
	var filename string
	if len(a.ApiStatsFile) == 0 {
		filename = fmt.Sprintf("%s_api_stats.json", time.Now().UTC().Format("2006-01-02_T03-04-05"))
	} else {
		filename = a.ApiStatsFile
	}
	fmt.Println("Dumping api stats")
	byteJSONRepresentation, err := json.Marshal(a.Stats)
	if err != nil {
		fmt.Println(fmt.Errorf("could not parse api stats: %w", err).Error())
		return
	}
	err = os.WriteFile(filename, byteJSONRepresentation, 0644)
	if err != nil {
		fmt.Println(fmt.Errorf("could not write api stats to file: %w", err).Error())
		return
	}
	fmt.Println("Api stats dumped to file:", filename)
}

// Custom middleware for chi router. It should be first registered with *chi.Mux.Use() and then started with the same context as the api for simultaneous shutdown with the http server itself.
func (a *DefaultApiStatTracker) ApiStatTracker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("track stats here")
		next.ServeHTTP(w, r)
	})
}
