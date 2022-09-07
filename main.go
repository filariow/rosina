package main

import (
	"context"
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
	EnvWatererPin1      = "ROSINA_WATERER_PIN1"
	EnvWatererPin2      = "ROSINA_WATERER_PIN2"
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

func getPinNumer(envVar string) (uint8, error) {
	wp := os.Getenv(envVar)
	if wp == "" {
		return 0, fmt.Errorf("Waterer Pin environment variable (%s) must be set", envVar)
	}

	n, err := strconv.ParseUint(wp, 10, 8)
	if err != nil {
		return 0, fmt.Errorf(
			"error parsing Waterer Pin environment variable (%s) value (%s) to uint8: %w",
			envVar, wp, err)
	}

	return uint8(n), nil
}

func run() error {
	ctx := context.Background()

	ss, err := getSchedules()
	if err != nil {
		return err
	}
	log.Printf("retrieved schedules: %v", ss)

	p1, err := getPinNumer(EnvWatererPin1)
	if err != nil {
		return err
	}
	log.Printf("waterer pin1: %d", p1)

	p2, err := getPinNumer(EnvWatererPin2)
	if err != nil {
		return err
	}
	log.Printf("waterer pin2: %d", p2)

	pin1, err := rpin.New(p1)
	if err != nil {
		return fmt.Errorf("error accessing waterer pin1 (pin %d): %w", p1, err)
	}

	pin2, err := rpin.New(p2)
	if err != nil {
		return fmt.Errorf("error accessing waterer pin1 (pin %d): %w", p1, err)
	}

	s := gocron.NewScheduler(time.UTC)
	for _, d := range ss {
		s.
			Cron(d.CronExpr).
			Do(buildWaterer(pin1, pin2, d.DurationSeconds))
		log.Printf("added cron job for expression %s", d.CronExpr)
	}

	log.Printf("starting cronjob scheduler")

	go s.StartAsync()

	<-ctx.Done()
	return ctx.Err()
}

func buildWaterer(pin1, pin2 rpin.OutPin, seconds uint64) func() {
	w := water.New(pin1, pin2)

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
