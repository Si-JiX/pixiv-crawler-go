package structs

import "time"

var AuthorStruct = struct {
	User    User      `json:"user"`
	Illusts []Illusts `json:"illusts"`
	NextURL any       `json:"next_url"`
}{}

type ProfileImageUrls struct {
	Medium string `json:"medium"`
}
type User struct {
	ID                   int              `json:"id"`
	Name                 string           `json:"name"`
	Account              string           `json:"account"`
	ProfileImageUrls     ProfileImageUrls `json:"profile_image_urls"`
	IsFollowed           bool             `json:"is_followed"`
	IsAccessBlockingUser bool             `json:"is_access_blocking_user"`
}
type ImageUrls struct {
	SquareMedium string `json:"square_medium"`
	Medium       string `json:"medium"`
	Large        string `json:"large"`
}
type Tags struct {
	Name           string      `json:"name"`
	TranslatedName interface{} `json:"translated_name"`
}
type MetaSinglePage struct {
	OriginalImageURL string `json:"original_image_url"`
}
type Illusts struct {
	ID             int            `json:"id"`
	Title          string         `json:"title"`
	Type           string         `json:"type"`
	ImageUrls      ImageUrls      `json:"image_urls"`
	Caption        string         `json:"caption"`
	Restrict       int            `json:"restrict"`
	User           User           `json:"user"`
	Tags           []Tags         `json:"tags"`
	Tools          []string       `json:"tools"`
	CreateDate     time.Time      `json:"create_date"`
	PageCount      int            `json:"page_count"`
	Width          int            `json:"width"`
	Height         int            `json:"height"`
	SanityLevel    int            `json:"sanity_level"`
	XRestrict      int            `json:"x_restrict"`
	Series         interface{}    `json:"series"`
	MetaSinglePage MetaSinglePage `json:"meta_single_page"`
	MetaPages      []interface{}  `json:"meta_pages"`
	TotalView      int            `json:"total_view"`
	TotalBookmarks int            `json:"total_bookmarks"`
	IsBookmarked   bool           `json:"is_bookmarked"`
	Visible        bool           `json:"visible"`
	IsMuted        bool           `json:"is_muted"`
	TotalComments  int            `json:"total_comments"`
}
