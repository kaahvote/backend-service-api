package data

type BaseFields struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"-"`
}
