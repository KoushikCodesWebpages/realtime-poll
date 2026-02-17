package ws

import (
	"sync"
	"time"
	"context"
	"realtime-poll/internal/repository"
)

type Room struct {
	PollID string

	// identity -> connection
	clients map[string]*Connection

	join      chan *Connection
	leave     chan *Connection
	broadcast chan internalEvent
	presenceDirty bool
	presenceTimer *time.Timer

	closed bool

	// authoritative in-memory vote projection
	optionVotes map[string]int

	pendingVotes map[string]int
	flushTimer   *time.Timer
	flushDirty   bool
	
	mu sync.RWMutex
}

func newRoom(pollID string) *Room {

	r := &Room{
		PollID:      pollID,
		clients:     make(map[string]*Connection),
		join:        make(chan *Connection),
		leave:       make(chan *Connection),
		broadcast:   make(chan internalEvent, 32),
		optionVotes: make(map[string]int),
		pendingVotes: make(map[string]int),
	}

	go r.run()
	go startExpiryTimer(r)

	return r
}

func (r *Room) run() {
	for {
		select {

	case c := <-r.join:

		key := c.sess.IdentityKey()

		if old, exists := r.clients[key]; exists {
			old.close()
			delete(r.clients, key)
		}

		// 🔴 LOAD POLL
		poll, _ := repository.GetPollByID(context.Background(), r.PollID)

		// 🔴 DETERMINE PERMISSIONS
		assignRole(c.sess, poll.OwnerID)
		assignVisibility(c.sess, poll)

		r.clients[key] = c

		// 🔴 NOW BUILD SNAPSHOT
		if BuildPollSnapshot != nil {
			data, err := BuildPollSnapshot(c.sess)
			if err == nil {
				r.hydrateVotes(data)
				c.send <- Envelope{
					Type: EventPollState,
					Data: data,
				}
			}
		}

		r.schedulePresence()


		// ---------- LEAVE ----------
		case c := <-r.leave:

			for key, conn := range r.clients {
				if conn == c {
					delete(r.clients, key)
					break
				}
			}

			c.close()
			r.schedulePresence()

		// ---------- EVENTS ----------
		case ev := <-r.broadcast:

			if ev.Type == EventVoteDelta {
				r.queueVote(ev)
				continue
			}

			r.emit(ev)

		}
	}
}


// hydrateVotes loads snapshot vote totals into memory
func (r *Room) hydrateVotes(snapshot any) {

	data, ok := snapshot.(map[string]any)
	if !ok {
		return
	}

	opts, ok := data["options"].([]map[string]any)
	if !ok {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, o := range opts {
		v, exists := o["votes"]
		if !exists {
			continue
		}

		if id, ok := o["option_id"].(string); ok {
			if votes, ok := v.(int); ok {
				r.optionVotes[id] = votes
			}
		}
	}
}

// apply delta and return authoritative total
func (r *Room) getCurrentVotes(optionID string, delta int) int {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.optionVotes[optionID] += delta
	return r.optionVotes[optionID]
}




func (r *Room) queueVote(ev internalEvent) {

	vd := ev.Data.(voteDelta)

	r.mu.Lock()
	r.pendingVotes[vd.OptionID] += vd.Delta
	r.mu.Unlock()

	r.scheduleFlush()
}

func (r *Room) flushVotes() {

    r.mu.Lock()

    updates := make(map[string]int, len(r.pendingVotes))

    for id, delta := range r.pendingVotes {
        total := r.applyDelta(id, delta) // no nested lock
        updates[id] = total
    }

    r.pendingVotes = make(map[string]int)
    r.flushDirty = false

    r.mu.Unlock()

    // emit outside lock
    for optionID, total := range updates {
        r.emit(internalEvent{
            Type: EventVoteDelta,
            Data: map[string]any{
                "option_id": optionID,
                "votes":     total,
            },
        })
    }
}

func (r *Room) scheduleFlush() {

    if r.flushDirty {
        return
    }

    r.flushDirty = true

    time.AfterFunc(25*time.Millisecond, func() {
        r.flushVotes()
    })
}


func (r *Room) applyDelta(optionID string, delta int) int {
    r.optionVotes[optionID] += delta
    return r.optionVotes[optionID]
}
