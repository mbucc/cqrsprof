/*
   Copyright (c) 2015, Mark Bucciarelli <mkbucc@gmail.com>

   Permission to use, copy, modify, and/or distribute this software
   for any purpose with or without fee is hereby granted, provided
   that the above copyright notice and this permission notice
   appear in all copies.

   THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL
   WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED
   WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL
   THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR
   CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
   LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT,
   NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN
   CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	. "github.com/mbucc/cqrs"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

type ShoutCommand struct {
	id               AggregateID
	Comment          string
	supportsRollback bool
}

func (c *ShoutCommand) ID() AggregateID         { return c.id }
func (c *ShoutCommand) BeginTransaction() error { return nil }
func (c *ShoutCommand) Commit() error           { return nil }
func (c *ShoutCommand) Rollback() error         { return nil }
func (c *ShoutCommand) SupportsRollback() bool  { return c.supportsRollback }

type HeardEvent struct {
	BaseEvent
	id    AggregateID
	Heard string
}

func (e *HeardEvent) ID() AggregateID { return e.id }

type NullAggregate struct {
	id AggregateID
}

func (eh NullAggregate) Handle(c Command) (a []Event, err error) {
	a = make([]Event, 1)
	c1 := c.(*ShoutCommand)
	a[0] = &HeardEvent{id: c1.ID(), Heard: c1.Comment}
	return a, nil
}

func (eh NullAggregate) ID() AggregateID               { return eh.id }
func (eh NullAggregate) New(id AggregateID) Aggregator { return &NullAggregate{id} }
func (eh NullAggregate) ApplyEvents([]Event)           {}

type NullEventListener struct{}

func (h *NullEventListener) Apply(e Event) error   { return nil }
func (h *NullEventListener) Reapply(e Event) error { return nil }

var aggregates = flag.Int("a", 100, "Number of aggregate IDs")
var commands = flag.Int("e", 1000, "Number of commands to process")

func main() {
	var id AggregateID

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Use the same sequence of random numbers
	// in case we need to reproduce something.
	rand.Seed(42)

	// store := &FileSystemEventStore{RootDir: "/tmp"}
	store := NewSqliteEventStore("/tmp/cqrs.db")

	RegisterEventListeners(new(HeardEvent), new(NullEventListener))
	RegisterEventStore(store)
	RegisterCommandAggregator(new(ShoutCommand), NullAggregate{})

	for i := 0; i < *commands; i++ {
		id = AggregateID(rand.Intn(*aggregates))
		SendCommand(&ShoutCommand{id, fmt.Sprintf("hello from command #%d", i), false})
	}
}
