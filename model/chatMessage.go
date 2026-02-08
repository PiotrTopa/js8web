package model

import (
	"database/sql"
	"sort"
	"time"
)

// ChatMessage is a unified wrapper for RX packets and TX frames,
// allowing them to be interleaved by timestamp in the chat view.
type ChatMessage struct {
	// Common fields for sorting
	Id        int64     `json:"Id"`
	Timestamp time.Time `json:"Timestamp"`
	Type      string    `json:"Type"`

	// RX packet fields (present when Type is RX.ACTIVITY, RX.DIRECTED, RX.DIRECTED.ME)
	Dial      uint32 `json:"Dial,omitempty"`
	Channel   uint16 `json:"Channel,omitempty"`
	Freq      uint32 `json:"Freq,omitempty"`
	Offset    uint16 `json:"Offset,omitempty"`
	Snr       int16  `json:"Snr,omitempty"`
	Mode      string `json:"Mode,omitempty"`
	Speed     string `json:"Speed,omitempty"`
	TimeDrift int16  `json:"TimeDrift,omitempty"`
	Grid      string `json:"Grid,omitempty"`
	From      string `json:"From,omitempty"`
	To        string `json:"To,omitempty"`
	Text      string `json:"Text,omitempty"`
	Command   string `json:"Command,omitempty"`
	Extra     string `json:"Extra,omitempty"`

	// TX frame fields (present when Type is TX.FRAME)
	Selected string `json:"Selected,omitempty"`
}

func chatMessageFromRxPacket(p *RxPacketObj) ChatMessage {
	return ChatMessage{
		Id:        p.Id,
		Timestamp: p.Timestamp,
		Type:      p.Type,
		Dial:      p.Dial,
		Channel:   p.Channel,
		Freq:      p.Freq,
		Offset:    p.Offset,
		Snr:       p.Snr,
		Mode:      p.Mode,
		Speed:     p.Speed,
		TimeDrift: p.TimeDrift,
		Grid:      p.Grid,
		From:      p.From,
		To:        p.To,
		Text:      p.Text,
		Command:   p.Command,
		Extra:     p.Extra,
	}
}

func chatMessageFromTxFrame(f *TxFrameObj) ChatMessage {
	return ChatMessage{
		Id:        -f.Id, // Negative to avoid ID collisions with RX packets
		Timestamp: f.Timestamp,
		Type:      "TX.FRAME",
		Dial:      f.Dial,
		Channel:   f.Channel,
		Freq:      f.Freq,
		Offset:    f.Offset,
		Mode:      f.Mode,
		Speed:     f.Speed,
		Selected:  f.Selected,
	}
}

// FetchChatMessages returns RX packets and TX frames merged and sorted by timestamp.
func FetchChatMessages(db *sql.DB, filter *RxPacketFilter, startTime time.Time, direction string, limit int) ([]ChatMessage, error) {
	rxPackets, err := FetchRxPacketList(db, filter, startTime, direction)
	if err != nil {
		return nil, err
	}

	txFrames, err := FetchTxFrameList(db, startTime, direction)
	if err != nil {
		return nil, err
	}

	messages := make([]ChatMessage, 0, len(rxPackets)+len(txFrames))
	for i := range rxPackets {
		messages = append(messages, chatMessageFromRxPacket(&rxPackets[i]))
	}
	for i := range txFrames {
		messages = append(messages, chatMessageFromTxFrame(&txFrames[i]))
	}

	sort.Slice(messages, func(i, j int) bool {
		if messages[i].Timestamp.Equal(messages[j].Timestamp) {
			return messages[i].Id < messages[j].Id
		}
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})

	if len(messages) > limit {
		if direction == "after" {
			messages = messages[:limit]
		} else {
			messages = messages[len(messages)-limit:]
		}
	}

	return messages, nil
}
