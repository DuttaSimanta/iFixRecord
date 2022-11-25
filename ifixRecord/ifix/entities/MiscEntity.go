package entities

type RecorddifferentionSingle struct {
	Id               int64  `json:"id"`
	Parentid         int64  `json:"parentid"`
	Name             string `json:"name"`
	Recorddifftypeid int64  `json:"recorddifftypeid"`
}
