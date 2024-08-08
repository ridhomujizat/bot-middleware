package botpress

import (
	"bot-middleware/internal/webhook"
	"encoding/json"
	"time"
)

type BotpressMessageType string

const (
	TEXT          BotpressMessageType = "text"
	SINGLE_CHOICE BotpressMessageType = "single-choice"
	CAROUSEL      BotpressMessageType = "carousel"
	POSTBACK      BotpressMessageType = "postback"
)

type AskPayloadBotpresDTO struct {
	Type     BotpressMessageType  `json:"type"`
	Text     string               `json:"text"`
	Metadata webhook.AttributeDTO `json:"metadata"`
	Payload  string               `json:"payload"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRespon struct {
	Token   string `json:"token"`
	BaseURL string `json:"base_url"`
}
type LoginBotPressDTO struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Jwt string `json:"jwt"`
	Exp int64  `json:"exp"`
}

type BotResponse struct {
}

func UnmarshalBotpressRespon(data []byte) (BotpressRespon, error) {
	var r BotpressRespon
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *BotpressRespon) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type BotpressRespon struct {
	Responses   []Response             `json:"responses"`
	Nlu         Nlu                    `json:"nlu"`
	Suggestions []interface{}          `json:"suggestions"`
	State       State                  `json:"state"`
	Decision    BotpressResponDecision `json:"decision"`
}

type BotpressResponDecision struct {
	Decision      DecisionDecision `json:"decision"`
	Confidence    int64            `json:"confidence"`
	Payloads      []interface{}    `json:"payloads"`
	Source        string           `json:"source"`
	SourceDetails string           `json:"sourceDetails"`
}

type DecisionDecision struct {
	Reason string `json:"reason"`
	Status string `json:"status"`
}

type Nlu struct {
	Entities         []interface{} `json:"entities"`
	Language         string        `json:"language"`
	Ambiguous        bool          `json:"ambiguous"`
	Slots            Slots         `json:"slots"`
	Intent           Intent        `json:"intent"`
	Intents          []interface{} `json:"intents"`
	Errored          bool          `json:"errored"`
	IncludedContexts []string      `json:"includedContexts"`
	MS               int64         `json:"ms"`
}

type Intent struct {
	Name       string `json:"name"`
	Confidence int64  `json:"confidence"`
	Context    string `json:"context"`
}

type Slots struct {
}

type Response struct {
	Type                string   `json:"type"`
	Skill               string   `json:"skill"`
	Workflow            Slots    `json:"workflow"`
	Text                string   `json:"text"`
	IsDropdown          bool     `json:"isDropdown"`
	DropdownPlaceholder string   `json:"dropdownPlaceholder"`
	Choices             []Choice `json:"choices"`
	Markdown            bool     `json:"markdown"`
	Typing              bool     `json:"typing"`
}

type Choice struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type State struct {
	Stacktrace []Stacktrace `json:"__stacktrace"`
	User       Slots        `json:"user"`
	Context    Context      `json:"context"`
	Session    Session      `json:"session"`
	Temp       Temp         `json:"temp"`
}

type Context struct {
	CurrentFlow  string       `json:"currentFlow"`
	CurrentNode  string       `json:"currentNode"`
	PreviousFlow string       `json:"previousFlow"`
	PreviousNode string       `json:"previousNode"`
	JumpPoints   []Stacktrace `json:"jumpPoints"`
	Queue        Queue        `json:"queue"`
}

type Stacktrace struct {
	Flow string `json:"flow"`
	Node string `json:"node"`
}

type Queue struct {
	Instructions []Instruction `json:"instructions"`
}

type Instruction struct {
	Type string  `json:"type"`
	Fn   string  `json:"fn"`
	Node *string `json:"node,omitempty"`
}

type Session struct {
	LastMessages   []LastMessage `json:"lastMessages"`
	Workflows      Slots         `json:"workflows"`
	Slots          Slots         `json:"slots"`
	Sid            string        `json:"sid"`
	UniqueId       string        `json:"unique_id"`
	AccountID      string        `json:"accountId"`
	ChannelSources string        `json:"channel_sources"`
	Name           string        `json:"name"`
}

type LastMessage struct {
	EventID         string    `json:"eventId"`
	IncomingPreview string    `json:"incomingPreview"`
	ReplyConfidence int64     `json:"replyConfidence"`
	ReplySource     string    `json:"replySource"`
	ReplyDate       time.Time `json:"replyDate"`
	ReplyPreview    string    `json:"replyPreview"`
}

type Temp struct {
	SkillChoiceInvalidCountM59Czq3Kso int64       `json:"skill-choice-invalid-count-m59czq3kso"`
	SkillChoiceValidM59Czq3Kso        interface{} `json:"skill-choice-valid-m59czq3kso"`
	SkillChoiceRetM59Czq3Kso          interface{} `json:"skill-choice-ret-m59czq3kso"`
}

type BotPressResponseDTO struct {
	Responses  []Response   `json:"responses"`
	State      string       `json:"state"`
	Stacktrace []Stacktrace `json:"__stacktrace"`
	BotDate    time.Time    `json:"botDate"`
}
