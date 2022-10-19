package structs

import "time"

var IllustStruct = struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    struct {
		IllustID      string    `json:"illustId"`
		IllustTitle   string    `json:"illustTitle"`
		IllustComment string    `json:"illustComment"`
		ID            string    `json:"id"`
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		IllustType    int       `json:"illustType"`
		CreateDate    time.Time `json:"createDate"`
		UploadDate    time.Time `json:"uploadDate"`
		Restrict      int       `json:"restrict"`
		XRestrict     int       `json:"xRestrict"`
		Sl            int       `json:"sl"`
		Urls          struct {
			Mini     string `json:"mini"`
			Thumb    string `json:"thumb"`
			Small    string `json:"small"`
			Regular  string `json:"regular"`
			Original string `json:"original"`
		} `json:"urls"`
		Tags struct {
			AuthorID string `json:"authorId"`
			IsLocked bool   `json:"isLocked"`
			Tags     []struct {
				Tag         string `json:"tag"`
				Locked      bool   `json:"locked"`
				Deletable   bool   `json:"deletable"`
				UserID      string `json:"userId,omitempty"`
				Translation struct {
					En string `json:"en"`
				} `json:"translation"`
				UserName string `json:"userName,omitempty"`
			} `json:"tags"`
			Writable bool `json:"writable"`
		} `json:"tags"`
		Alt                  string        `json:"alt"`
		StorableTags         []string      `json:"storableTags"`
		UserID               string        `json:"userId"`
		UserName             string        `json:"userName"`
		UserAccount          string        `json:"userAccount"`
		LikeData             bool          `json:"likeData"`
		Width                int           `json:"width"`
		Height               int           `json:"height"`
		PageCount            int           `json:"pageCount"`
		BookmarkCount        int           `json:"bookmarkCount"`
		LikeCount            int           `json:"likeCount"`
		CommentCount         int           `json:"commentCount"`
		ResponseCount        int           `json:"responseCount"`
		ViewCount            int           `json:"viewCount"`
		BookStyle            int           `json:"bookStyle"`
		IsHowto              bool          `json:"isHowto"`
		IsOriginal           bool          `json:"isOriginal"`
		ImageResponseOutData []interface{} `json:"imageResponseOutData"`
		ImageResponseData    []interface{} `json:"imageResponseData"`
		ImageResponseCount   int           `json:"imageResponseCount"`
		PollData             interface{}   `json:"pollData"`
		SeriesNavData        interface{}   `json:"seriesNavData"`
		DescriptionBoothID   interface{}   `json:"descriptionBoothId"`
		DescriptionYoutubeID interface{}   `json:"descriptionYoutubeId"`
		ComicPromotion       interface{}   `json:"comicPromotion"`
		FanboxPromotion      interface{}   `json:"fanboxPromotion"`
		ContestBanners       []interface{} `json:"contestBanners"`
		IsBookmarkable       bool          `json:"isBookmarkable"`
		BookmarkData         interface{}   `json:"bookmarkData"`
		ContestData          interface{}   `json:"contestData"`
		ZoneConfig           struct {
			Responsive struct {
				URL string `json:"url"`
			} `json:"responsive"`
			Rectangle struct {
				URL string `json:"url"`
			} `json:"rectangle"`
			Five00X500 struct {
				URL string `json:"url"`
			} `json:"500x500"`
			Header struct {
				URL string `json:"url"`
			} `json:"header"`
			Footer struct {
				URL string `json:"url"`
			} `json:"footer"`
			ExpandedFooter struct {
				URL string `json:"url"`
			} `json:"expandedFooter"`
			Logo struct {
				URL string `json:"url"`
			} `json:"logo"`
			Relatedworks struct {
				URL string `json:"url"`
			} `json:"relatedworks"`
		} `json:"zoneConfig"`
		ExtraData struct {
			Meta struct {
				Title              string `json:"title"`
				Description        string `json:"description"`
				Canonical          string `json:"canonical"`
				AlternateLanguages struct {
					Ja string `json:"ja"`
					En string `json:"en"`
				} `json:"alternateLanguages"`
				DescriptionHeader string `json:"descriptionHeader"`
				Ogp               struct {
					Description string `json:"description"`
					Image       string `json:"image"`
					Title       string `json:"title"`
					Type        string `json:"type"`
				} `json:"ogp"`
				Twitter struct {
					Description string `json:"description"`
					Image       string `json:"image"`
					Title       string `json:"title"`
					Card        string `json:"card"`
				} `json:"twitter"`
			} `json:"meta"`
		} `json:"extraData"`
		TitleCaptionTranslation struct {
			WorkTitle   interface{} `json:"workTitle"`
			WorkCaption interface{} `json:"workCaption"`
		} `json:"titleCaptionTranslation"`
		IsUnlisted bool        `json:"isUnlisted"`
		Request    interface{} `json:"request"`
		CommentOff int         `json:"commentOff"`
	} `json:"body"`
}{}
