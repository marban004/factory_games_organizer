//     This is Factory Games Organizer api. Api is responsible for creating, updating and authenicating api users, CRUD operations on database associated with the api and provides production calculator service.
//     Copyright (C) 2025  Marek Banaś

//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.

//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see https://www.gnu.org/licenses/.

package custommiddleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type DefaultApiStatTracker struct {
	stats        *orderedmap.OrderedMap[string, map[string]int]
	mu           sync.Mutex
	Ctx          context.Context
	MaxLen       uint64
	Period       int64
	currentKey   int64
	ApiStatsFile string
	DumpStats    bool
}

func (a *DefaultApiStatTracker) GetStats() *orderedmap.OrderedMap[string, map[string]int] {
	a.mu.Lock()
	stats := a.stats
	a.mu.Unlock()
	return stats
}

func (a *DefaultApiStatTracker) clearUpOldData() {
	for a.stats.Len() > int(a.MaxLen) {
		a.stats.Delete(a.stats.Oldest().Key)
	}
}

// Starts this instance of api stat tracker. Blocks until the passed context is cancelled or expires.
func (a *DefaultApiStatTracker) StartTracker(ctx context.Context) {
	a.stats = orderedmap.New[string, map[string]int]()
	currentTime := time.Now().UnixMilli()
	a.currentKey = currentTime - (currentTime % a.Period)
	a.makeEmptyStatRecordForCurrentPeriod()
	fmt.Println("Started tracking stats")
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
	a.stats.Set(mapKey, make(map[string]int))
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
	byteJSONRepresentation, err := json.Marshal(a.stats)
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
		a.mu.Lock()
		recordingKey := time.UnixMilli(a.currentKey).Format("2006-01-02T15:04:05.000Z07:00")
		recordMap, exists := a.stats.Get(recordingKey)
		if !exists {
			fmt.Println("Unexpected error, map key does not exist")
		}
		endpoint, _, _ := strings.Cut(r.RequestURI, "?")
		recordMap[endpoint+" "+r.Method] += 1
		a.stats.Set(fmt.Sprint(recordingKey), recordMap)
		a.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
