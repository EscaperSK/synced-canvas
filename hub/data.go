package hub

type message struct {
	Index int  `json:"index"`
	Fill  bool `json:"fill"`

	Command string `json:"command"`
}
