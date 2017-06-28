package brutemachine

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// This structure contains some runtime statistics.
type Statistics struct {
	// Time the execution started
	Start   time.Time
	// Time the execution finished
	Stop    time.Time
	// Total duration of the execution
	Total   time.Duration
	// Executions per second
	Eps     float64
	// Total number of executions
	Execs   uint64
	// Total number of executions with positive results.
	Results uint64
}

// This is where the main logic goes.
type RunHandler func(line string) interface{}
// This is where positive results are handled.
type ResultHandler func(result interface{})

// The main object.
type Machine struct {
	// Runtime statistics.
	Stats       Statistics
	// Number of input consumers.
	consumers   uint
	// Dictionary file name.
	filename    string
	// Positive results channel.
	output      chan interface{}
	// Inputs channel.
	input       chan string
	// WaitGroup to stop while the machine is running.
	wait        sync.WaitGroup
	// Main logic handler.
	run_handler RunHandler
	// Positive results handler.
	res_handler ResultHandler
}

// Builds a new machine object, if consumers is less or equal than 0, CPU*2 will be used as default value.
func New( consumers int, filename string, run_handler RunHandler, res_handler ResultHandler) *Machine {
	workers := uint(0)
	if consumers <= 0 {
		workers = uint(runtime.NumCPU() * 2)
	} else {
		workers = uint(consumers)
	}

	return &Machine{
		Stats:       Statistics{Execs: 0, Results: 0, Eps: 0.0},
		consumers:   workers,
		filename:    filename,
		output:      make(chan interface{}),
		input:       make(chan string),
		wait:        sync.WaitGroup{},
		run_handler: run_handler,
		res_handler: res_handler,
	}
}

func (m *Machine) inputConsumer() {
	for in := range m.input {
		atomic.AddUint64(&m.Stats.Execs, 1)

		res := m.run_handler(in)
		if res != nil {
			atomic.AddUint64(&m.Stats.Results, 1)
			m.output <- res
		}
		m.wait.Done()
	}
}

func (m *Machine) outputConsumer() {
	for res := range m.output {
		m.res_handler(res)
	}
}

// Start the machine.
func (m *Machine) Start() error {
	// start a fixed amount of consumers for inputs
	for i := uint(0); i < m.consumers; i++ {
		go m.inputConsumer()
	}

	// start the output consumer on a goroutine
	go m.outputConsumer()

	m.Stats.Start = time.Now()

	lines, err := LineReader(m.filename, 0)
	if err != nil {
		return err
	}
	for line := range lines {
		m.wait.Add(1)
		m.input <- line
	}

	return nil
}

func (m *Machine) UpdateStats() {
	m.Stats.Stop = time.Now()
	m.Stats.Total = m.Stats.Stop.Sub(m.Stats.Start)
	m.Stats.Eps = float64(m.Stats.Execs) / m.Stats.Total.Seconds()
}

// Wait for all jobs to be completed.
func (m *Machine) Wait() {
	// wait for everything to be completed
	m.wait.Wait()
	m.UpdateStats()
}
