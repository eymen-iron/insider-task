package esrc

type Item struct {
	ID       string `json:"item_id"`
	Name     string `json:"name"`
	Locale   string `json:"locale"`
	Click    int    `json:"click"`
	Purchase int    `json:"purchase"`
}
