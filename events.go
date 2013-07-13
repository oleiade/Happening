package happening

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Event represents a event sent by the source
type Event struct {
	raw string

	From       string
	SentOn     int64
	ReceivedOn int64
	Type       string
}

// NewEvent initializes an event from it's component
func NewEvent(from string, sent_on int64, received_on int64, event_type string) *Event {
	return &Event{
		From:       from,
		SentOn:     sent_on,
		ReceivedOn: received_on,
		Type:       event_type,
	}
}

// NewEventFromRaw initializes an Event from it's raw representation
// which MSG_DELIMITER has been removed from.
func NewEventFromRaw(raw string) (*Event, error) {
	event := new(Event)
	err := event.FromRaw(raw)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// FromRaw instantiates and returns an initialized Event
// from it's raw description. As event flow splits event messages
// based on the MSG_DELIMITER and discards it, the method artificially
// restores the MSG_DELIMITER in Event.raw attribute.
func (e *Event) FromRaw(raw string) error {
	e.raw = raw + MSG_DELIMITER // Keep track of the raw version with MSG_DELIMITER
	parts := strings.Split(strings.Trim(raw, MSG_DELIMITER), string(EVENT_PARAMS_SEPARATOR))

	if len(parts) == 3 {
		e.From = parts[0]
		e.ReceivedOn = time.Now().Unix()
		e.Type = parts[2]

		sent_on, err := strconv.Atoi(parts[1])
		if err != nil {
			return errors.New(fmt.Sprintf("[Event.FromRaw] Couldn't parse timestamp: %s", err))
		} else {
			e.SentOn = int64(sent_on)
		}
	} else {
		return errors.New(fmt.Sprintf("[%s.FromRaw] Incomplete event received: %s", "Event", e.raw))
	}

	return nil
}
