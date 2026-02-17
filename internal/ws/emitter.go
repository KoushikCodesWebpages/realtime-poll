package ws

func (r *Room) emit(ev internalEvent) {

	for c := range r.clients {

		env := projectEvent(ev, c.sess)
		if env == nil {
			continue
		}

		select {
		case c.send <- *env:
		default:
			c.close()
			delete(r.clients, c)
		}
	}
}
func projectEvent(ev internalEvent, s *Session) *Envelope {

	switch ev.Type {

	case EventVoteDelta:

		switch s.Visibility {

		case ViewHidden:
			// voter knows something happened but no counts
			return &Envelope{
				Type: EventVoteDelta,
			}

		case ViewLive, ViewOwner, ViewFinal:
			return &Envelope{
				Type: EventVoteDelta,
				Data: ev.Data,
			}
		}

	case EventClosed:
		return &Envelope{
			Type: EventClosed,
			Data: ev.Data,
		}
	}

	return nil
}
