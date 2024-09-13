package workerWhatsapp

type OutgoingTextSocioconnect struct {
	RecipientType    string `json:"recipient_type"`
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             Text   `json:"text"`
}

type OutgoingImageSocioconnect struct {
	RecipientType    string `json:"recipient_type"`
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Image            Image  `json:"image"`
}

type OutgoingButtonSocioconnect struct {
	RecipientType    string      `json:"recipient_type"`
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Interactive      Interactive `json:"interactive"`
}

type OutgoingListSocioconnect struct {
	RecipientType    string      `json:"recipient_type"`
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Interactive      Interactive `json:"interactive"`
}

type Text struct {
	Body string `json:"body"`
}

type Image struct {
	ID      string `json:"id"`
	Caption string `json:"caption,omitempty"`
	Link    string `json:"link,omitempty"`
}

type Interactive struct {
	Type   string `json:"type"`
	Body   Body   `json:"body"`
	Action Action `json:"action"`
	Header Header `json:"header,omitempty"`
}

type Body struct {
	Text string `json:"text"`
}

type Action struct {
	Buttons  []Button  `json:"buttons,omitempty"`
	Button   string    `json:"button,omitempty"`
	Sections []Section `json:"sections,omitempty"`
}

type Button struct {
	Type  string `json:"type"`
	Reply Reply  `json:"reply"`
}

type Reply struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}

type Section struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          string `json:"id"`
}

type Header struct {
	Type  string `json:"type"`
	Image Image  `json:"image"`
}
