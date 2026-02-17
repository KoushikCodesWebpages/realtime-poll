package ws
import (
	"time"
)
func (r *Room) emit(ev internalEvent) {

	for key, c := range r.clients {

		env := r.projectEvent(ev, c.sess)
		if env == nil {
			continue
		}

		select {
		case c.send <- *env:
		default:
			// dead connection
			c.close()
			delete(r.clients, key)
			r.schedulePresence()
		}
	}
}

func (r *Room) projectEvent(ev internalEvent, s *Session) *Envelope {

	switch ev.Type {

	// ---------- VOTES (already authoritative) ----------
	case EventVoteDelta:

		// hidden viewers don't see numbers
		if s.Visibility == ViewHidden && !s.IsOwner {
			return &Envelope{Type: EventVoteDelta}
		}

		return &Envelope{
			Type: EventVoteDelta,
			Data: ev.Data,
		}

	// ---------- POLL CLOSED ----------
	case EventClosed:
		return &Envelope{
			Type: EventClosed,
			Data: ev.Data,
		}

	// ---------- PRESENCE ----------
	case EventPresence:
		return &Envelope{
			Type: EventPresence,
			Data: ev.Data,
		}
	}

	return nil
}

// func (r *Room) broadcastPresence() {

// 	r.mu.RLock()
// 	count := len(r.clients)
// 	r.mu.RUnlock()

// 	ev := internalEvent{
// 		Type: EventPresence,
// 		Data: map[string]any{
// 			"viewers": count,
// 		},
// 	}

// 	r.emit(ev)
// }

// coalesced presence broadcaster

func (r *Room) schedulePresence() {

	r.mu.Lock()

	// already scheduled
	if r.presenceDirty {
		r.mu.Unlock()
		return
	}

	r.presenceDirty = true
	r.mu.Unlock()

	time.AfterFunc(20*time.Millisecond, func() {

		r.mu.RLock()
		count := len(r.clients)
		r.mu.RUnlock()

		r.emit(internalEvent{
			Type: EventPresence,
			Data: map[string]any{
				"viewers": count,
			},
		})

		r.mu.Lock()
		r.presenceDirty = false
		r.mu.Unlock()
	})
}