package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/filariow/rosina/internal/rpin"
	"github.com/filariow/rosina/pkg/water"
	"github.com/go-co-op/gocron"
)

const (
	EnvWatererPin       = "ROSINA_WATERER_PIN"
	EnvWatererSchedules = "ROSINA_WATERER_SCHEDULES"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type Schedule struct {
	CronExpr        string
	DurationSeconds uint64
}

func parseSchedule(schedule string) (Schedule, error) {
	ss := strings.Split(schedule, ",")
	if len(ss) != 2 {
		return Schedule{}, fmt.Errorf("error parsing schedule: expected format 'CRONEXPR,DURATION_IN_SEC' not respected by '%s'", schedule)
	}

	d, err := strconv.ParseUint(ss[1], 10, 64)
	if err != nil {
		return Schedule{},
			fmt.Errorf("error parsing provided duration (%s) as uint64: %w", ss[1], err)
	}

	return Schedule{
		CronExpr:        ss[0],
		DurationSeconds: d,
	}, nil
}

func getSchedules() ([]Schedule, error) {
	ess := os.Getenv(EnvWatererSchedules)
	if ess == "" {
		return nil,
			fmt.Errorf("Waterer Scheduler environment variable (%s) must be set", EnvWatererSchedules)
	}

	sss := strings.Split(ess, ";")
	ss := make([]Schedule, len(sss))

	for i := 0; i < len(sss); i++ {
		s, err := parseSchedule(sss[i])
		if err != nil {
			return nil, fmt.Errorf("error parsing schedules: %w", err)
		}
		ss[i] = s
	}
	return ss, nil
}

func getPinNumer() (uint8, error) {
	wp := os.Getenv(EnvWatererPin)
	if wp == "" {
		return 0, fmt.Errorf("Waterer Pin environment variable (%s) must be set", EnvWatererPin)
	}

	n, err := strconv.ParseUint(wp, 10, 8)
	if err != nil {
		return 0, fmt.Errorf(
			"error parsing Waterer Pin environment variable (%s) value (%s) to uint8: %w",
			EnvWatererPin, wp, err)
	}

	return uint8(n), nil
}

func run() error {
	ss, err := getSchedules()
	if err != nil {
		return err
	}

	p, err := getPinNumer()
	if err != nil {
		return err
	}

	s := gocron.NewScheduler(time.UTC)
	for _, d := range ss {
		s.
			Cron(d.CronExpr).
			Do(buildWaterer(p, d.DurationSeconds))

	}

	return nil
}

func buildWaterer(pin uint8, seconds uint64) func() {
	p := rpin.New(pin)
	w := water.New(p)

	return func() {
		log.Println("Opening water")
		w.Open()

		wt := time.Duration(seconds) * time.Second
		l := time.Now().Add(wt).UTC().String()
		log.Printf("Waiting %d seconds (i.e. %s)", seconds, l)
		time.Sleep(wt)

		log.Println("Opening water")
		w.Close()
	}
}
