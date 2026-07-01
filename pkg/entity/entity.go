package entity

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/zu1k/nali/pkg/dbif"
)

// Type classifies a parsed token (IPv4, IPv6, domain or plain text).
type Type uint

// Entity types, mirroring the dbif query types plus plain text.
const (
	TypeIPv4   = dbif.TypeIPv4
	TypeIPv6   = dbif.TypeIPv6
	TypeDomain = dbif.TypeDomain

	TypePlain = 100
)

// Entity is a parsed token from the input together with its lookup result.
type Entity struct {
	Loc  [2]int `json:"-"` // s[Loc[0]:Loc[1]]
	Type Type   `json:"type"`

	Text     string      `json:"ip"`
	InfoText string      `json:"text"`
	Source   string      `json:"source"`
	Info     interface{} `json:"info"`
}

// Json renders the entity as a JSON string.
func (e Entity) Json() string {
	jsonResult, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{"error": "json marshal failed: %s"}`, err.Error())
	}
	return string(jsonResult)
}

// Entities is an ordered collection of parsed entities.
type Entities []*Entity

func (es Entities) Len() int {
	return len(es)
}

func (es Entities) Less(i, j int) bool {
	return es[i].Loc[0] < es[j].Loc[0]
}

func (es Entities) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es Entities) String() string {
	var result strings.Builder
	for _, entity := range es {
		result.WriteString(entity.Text)
		if entity.Type != TypePlain && len(entity.InfoText) > 0 {
			result.WriteString("[" + entity.InfoText + "] ")
		}
	}
	return result.String()
}

// ColorString renders the entities as a colorized line for terminal output.
func (es Entities) ColorString() string {
	var line strings.Builder
	for _, e := range es {
		s := e.Text
		switch e.Type {
		case TypeIPv4:
			s = color.GreenString(e.Text)
		case TypeIPv6:
			s = color.BlueString(e.Text)
		case TypeDomain:
			s = color.YellowString(e.Text)
		}
		if e.Type != TypePlain && len(e.InfoText) > 0 {
			s += " [" + color.RedString(e.InfoText) + "] "
		}
		line.WriteString(s)
	}
	return line.String()
}

// Json renders the non-plain entities as newline-separated JSON objects.
func (es Entities) Json() string {
	var s strings.Builder
	for _, e := range es {
		if e.Type == TypePlain {
			continue
		}
		s.WriteString(e.Json() + "\n")
	}
	return s.String()
}
