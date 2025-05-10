/*
* Basic implementation of distributed Snowflake unique IDs
 */
package snowflake

import (
	"sync"
	"time"
)

type (
	// Timestamp represents a 41-bit, millisecond Unix timestamp
	Timestamp uint64
	// MachineID is a 10-bit a unique identifier for the program generating snowflakes
	MachineID uint16
	// SequenceNumber is a 12-bit counter for all snowflakes generated during a given Timestamp
	SequenceNumber uint16
)

// SnowflakeID represents a distributed snowflake ID, like the one original devised by Twitter.
// It is a 64-bit unique ID suitable for distributed systems
type SnowflakeID struct {
	Timestamp Timestamp
	MachineID MachineID
	Seq       SequenceNumber
}

// Counter is any type of object that generates new snowflake IDs
type Counter interface {
	Next() SnowflakeID
}

// SnowflakeCounter is a thread-safe Counter which will produces monotonically increasing SnowflakeIDs
type SnowflakeCounter struct {
	lastTime Timestamp
	nodeId   MachineID
	lastSeq  SequenceNumber
	mu       sync.Mutex
	Counter
}

// NewCounter constructs a new SnowflakeCounter for the given unique machine
func NewCounter(machineId MachineID) SnowflakeCounter {
	return SnowflakeCounter{
		lastTime: Timestamp(time.Now().UnixMilli()),
		nodeId:   machineId,
		lastSeq:  0,
		mu:       sync.Mutex{},
	}
}

// Next produces the next snowflake ID for the current time
func (sc *SnowflakeCounter) Next() SnowflakeID {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	now := Timestamp(time.Now().UnixMilli())
	if now > sc.lastTime {
		sc.lastSeq++
	} else {
		sc.lastTime = now
		sc.lastSeq = 0
	}

	return SnowflakeID{Timestamp: sc.lastTime, MachineID: sc.nodeId, Seq: sc.lastSeq}
}

// Pack combines a SnowflakeID into an unsigned 64-bit representation
func (s SnowflakeID) Pack() uint64 {
	return (uint64(s.Timestamp) << 22) | (uint64(s.MachineID) << 12) | uint64(s.Seq)
}

// Unpack maps a unit64 back to a SnowflakeID struct
func Unpack(id uint64) SnowflakeID {
	return SnowflakeID{
		Timestamp: Timestamp(id >> 22),
		MachineID: MachineID((id >> 12) & 0x3FF), // 10-bit ID
		Seq:       SequenceNumber(id & 0xFFF),    // 12-bit value
	}
}
