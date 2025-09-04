package dto

type Batch struct {
	ID        int    `db:"id"`
	LastBatch string `db:"last_batch"`
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

//title
//description
//link
//pubDate
