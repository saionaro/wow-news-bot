package types

type NewsItem struct {
	Title string
	Href  string
	Image string
	Hash  string
}

type Message struct {
	Text  string
	Image []byte
}
