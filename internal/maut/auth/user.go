package auth

type User struct {
	ID      string   `json:"id"`
	APIKey  bool     `json:"apikey"`
	Blocked bool     `json:"blocked"`
	Email   *string  `json:"email,omitempty"`
	Groups  []string `json:"groups"`
	Name    *string  `json:"name,omitempty"`
	Roles   []string `json:"roles"`

	Hash []byte `json:"hash"`
	Salt []byte `json:"salt"`
}
