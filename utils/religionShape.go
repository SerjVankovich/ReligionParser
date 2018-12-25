package utils

type ReligionShape struct {
	Title, Author, Description, FullDescription string
	Month, Day, Year, Shows, Repost, Rating     int
}

func NewRelShape(Title, Author, Description string, Month, Day, Year, Shows, Repost, Rating int) *ReligionShape {
	return &ReligionShape{Title, Author, Description, "", Month, Day, Year, Shows, Repost, Rating}
}
