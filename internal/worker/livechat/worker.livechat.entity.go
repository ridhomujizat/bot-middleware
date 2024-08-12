package workerLivechat

import (
	"bot-middleware/internal/entities"
	"net/url"
)

type Choices struct {
	Title string `json:"title" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type CarouselActions struct {
	Action  string `json:"action" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Payload string `json:"payload" validate:"required"`
}

type Carousel struct {
	Title    string            `json:"title" validate:"required"`
	Subtitle string            `json:"subtitle,omitempty" validate:"omitempty"`
	Image    *url.URL          `json:"image,omitempty" validate:"omitempty,url"`
	Actions  []CarouselActions `json:"actions" validate:"required,dive"`
}

type BotResponse struct {
	Type                entities.BotpressMessageType `json:"type" validate:"required,oneof=text single-choice carousel postback"`
	IsDropdown          bool                         `json:"isDropdown,omitempty" validate:"omitempty"`
	Text                string                       `json:"text" validate:"required"`
	DropdownPlaceholder string                       `json:"dropdownPlaceholder,omitempty" validate:"required_if=Type single-choice"`
	Choices             []Choices                    `json:"choices,omitempty" validate:"required_if=Type single-choice,dive"`
	Items               []Carousel                   `json:"items,omitempty" validate:"required_if=Type carousel,dive"`
}

type Outgoing struct {
	Token    string `json:"token" validate:"required"`
	Message  string `json:"message" validate:"required"`
	FromName string `json:"fromName" validate:"required"`
}

type Button struct {
	Label string `json:"label" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type MessageButton struct {
	Title  string   `json:"title" validate:"required"`
	Button []Button `json:"button" validate:"required,dive"`
}

type OutgoingButton struct {
	Token    string        `json:"token" validate:"required"`
	Message  MessageButton `json:"message" validate:"required,dive"`
	FromName string        `json:"fromName" validate:"required"`
}

type MessageCarouselSliderMenu struct {
	Label string `json:"label" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type MessageCarouselSlider struct {
	Title    string                      `json:"title" validate:"required"`
	Subtitle string                      `json:"subtitle,omitempty" validate:"omitempty"`
	Picture  *url.URL                    `json:"picture" validate:"required,url"`
	Menu     []MessageCarouselSliderMenu `json:"menu" validate:"required,dive"`
}

type MessageCarousel struct {
	MessageType string                  `json:"messageType" validate:"required"`
	Slider      []MessageCarouselSlider `json:"slider" validate:"required,dive"`
}

type OutgoingCarousel struct {
	Token    string          `json:"token" validate:"required"`
	Message  MessageCarousel `json:"message" validate:"required,dive"`
	FromName string          `json:"fromName" validate:"required"`
}
