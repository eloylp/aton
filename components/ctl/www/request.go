package www

type AddNodeRequest struct {
	Addr string `json:"address"`
}

type LoadCategoriesRequest struct {
	Categories []string `json:"categories"`
	Image      []byte   `json:"image"`
}

type AddTargetRequest struct {
	UUID       string `json:"uuid"`
	TargetAddr string `json:"target_addr"`
}
