package model

// Data represents request or response data
type Data struct {
	URL   string `json:"url"`
	Short string `json:"short"`
}

type Response struct{
	Short string `json:"short"`
}

type Short struct{
	Short string `json:"short"`
}

type URL struct{
	URL string `json:"url"`
}