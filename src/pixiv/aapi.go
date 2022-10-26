package pixiv

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"pixiv-cil/utils"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// AppPixivAPI -- App-API (6.x - app-api.pixiv.net)
type AppPixivAPI struct {
	sling   *sling.Sling
	timeout time.Duration
	proxy   *url.URL
}

func NewApp() *AppPixivAPI {
	s := sling.New().Base(API_BASE).Set("User-Agent", "PixivIOSApp/7.6.2 (iOS 12.2; iPhone9,1)").Set("App-Version", "7.6.2").Set("App-OS-VERSION", "12.2").Set("App-OS", "ios")
	return &AppPixivAPI{sling: s}
}

func (a *AppPixivAPI) request(path string, params, data interface{}, auth bool) (err error) {
	if auth {
		_, err = a.sling.New().Get(path).Set("Authorization", "Bearer "+TokenVariable).QueryStruct(params).ReceiveSuccess(data)
	} else {
		_, err = a.sling.New().Get(path).QueryStruct(params).ReceiveSuccess(data)
	}
	return err
}

func (a *AppPixivAPI) WithDownloadTimeout(timeout time.Duration) *AppPixivAPI {
	a.timeout = timeout
	return a
}

func (a *AppPixivAPI) WithDownloadProxy(proxy *url.URL) *AppPixivAPI {
	a.proxy = proxy
	return a
}

func (a *AppPixivAPI) post(path string, params, data interface{}, auth bool) (err error) {
	if auth {
		_, err = a.sling.New().Post(path).Set("Authorization", "Bearer "+TokenVariable).BodyForm(params).ReceiveSuccess(data)
	} else {
		_, err = a.sling.New().Post(path).BodyForm(params).ReceiveSuccess(data)
	}
	return err
}

type userDetailParams struct {
	UserID uint64 `url:"user_id,omitempty"`
	Filter string `url:"filter,omitempty"`
}

func (a *AppPixivAPI) UserDetail(uid uint64) (*UserDetail, error) {
	params := &userDetailParams{UserID: uid, Filter: "for_ios"}
	detail := &UserDetail{User: &User{}}
	if err := a.request(USER_DETAIL, params, detail, true); err != nil {
		return nil, err
	}
	return detail, nil
}

type userIllustsParams struct {
	UserID int    `url:"user_id,omitempty"`
	Filter string `url:"filter,omitempty"`
	Type   string `url:"type,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

// UserIllusts type: [illust, manga]
func (a *AppPixivAPI) UserIllusts(uid int, _type string, offset int) ([]Illust, int, error) {
	params := &userIllustsParams{
		UserID: uid,
		Filter: "for_ios",
		Type:   _type,
		Offset: offset,
	}
	data := &IllustsResponse{}
	if err := a.request(USER_AUTHOR, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

type userBookmarkIllustsParams struct {
	UserID        uint64 `url:"user_id,omitempty"`
	Restrict      string `url:"restrict,omitempty"`
	Filter        string `url:"filter,omitempty"`
	MaxBookmarkID int    `url:"max_bookmark_id,omitempty"`
	Tag           string `url:"tag,omitempty"`
}

// UserBookmarksIllust restrict: [public, private]
func (a *AppPixivAPI) UserBookmarksIllust(uid uint64, maxBookmarkID int, tag string) ([]Illust, int, error) {
	params := &userBookmarkIllustsParams{
		UserID:        uid,
		Restrict:      "public",
		Filter:        "for_ios",
		MaxBookmarkID: maxBookmarkID,
		Tag:           tag,
	}
	data := &IllustsResponse{}
	if err := a.request(BOOKMARKS, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

type illustFollowParams struct {
	Restrict string `url:"restrict,omitempty"`
	Offset   int    `url:"offset,omitempty"`
}

// IllustFollow restrict: [public, private]
func (a *AppPixivAPI) IllustFollow(restrict string, offset int) ([]Illust, int, error) {
	params := &illustFollowParams{Restrict: restrict, Offset: offset}
	data := &IllustsResponse{}
	if err := a.request(FOLLOW, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

type illustDetailParams struct {
	IllustID int `url:"illust_id,omitemtpy"`
}

// IllustDetail get a detailed illust with id
func (a *AppPixivAPI) IllustDetail(id int) (*Illust, error) {
	data := &IllustResponse{}
	params := &illustDetailParams{IllustID: id}
	if err := a.request(DETAIL, params, data, true); err != nil {
		return nil, err
	}
	return &data.Illust, nil
}

// Download a specific picture from pixiv id
func (a *AppPixivAPI) Download(id int, path string) (sizes []int64, err error) {
	illust, err := a.IllustDetail(id)
	if err != nil {
		err = errors.Wrapf(err, "illust %d detail error", id)
		return
	}
	if illust == nil {
		err = errors.Wrapf(err, "illust %d is nil", id)
		return
	}
	if illust.MetaSinglePage == nil {
		err = errors.Wrapf(err, "illust %d has no single page", id)
		return
	}

	var urls []string
	if illust.MetaSinglePage.OriginalImageURL == "" {
		for _, img := range illust.MetaPages {
			urls = append(urls, img.Images.Original)
		}
	} else {
		urls = append(urls, illust.MetaSinglePage.OriginalImageURL)
	}

	dclient := &http.Client{}
	if a.proxy != nil {
		dclient.Transport = &http.Transport{
			Proxy: http.ProxyURL(a.proxy),
		}
	}
	if a.timeout != 0 {
		dclient.Timeout = a.timeout
	}

	for _, u := range urls {
		size, e := download(dclient, u, path, filepath.Base(u))
		if e != nil {
			err = errors.Wrapf(e, "download url %s failed", u)
			return
		}
		sizes = append(sizes, size)
	}

	return
}

func (a *AppPixivAPI) ThreadDownloadImage(url string, Length int) {
	defer utils.WG.Done()
	utils.CurrentImageIndex += 1
	utils.CH <- struct{}{}
	dclient := &http.Client{}
	if a.proxy != nil {
		dclient.Transport = &http.Transport{
			Proxy: http.ProxyURL(a.proxy),
		}
	}
	if a.timeout != 0 {
		dclient.Timeout = a.timeout
	}
	_, e := download(dclient, url, "imageFile", filepath.Base(url))
	<-utils.CH
	if e != nil {
		fmt.Println(errors.Wrapf(e, "download url %s failed", url))
	}
	defer fmt.Printf("download image:%d/%d\r", utils.CurrentImageIndex, Length)

}

type illustCommentsParams struct {
	IllustID             uint64 `url:"illust_id,omitemtpy"`
	Offset               int    `url:"offset,omitempty"`
	IncludeTotalComments bool   `url:"include_total_comments,omitempty"`
}

// IllustComments Comments posted in a pixiv artwork
func (a *AppPixivAPI) IllustComments(illustID uint64, offset int, includeTotalComments bool) (*IllustComments, error) {
	data := &IllustComments{}
	params := &illustCommentsParams{
		IllustID:             illustID,
		IncludeTotalComments: includeTotalComments,
		Offset:               offset,
	}

	if err := a.request(COMMENTS, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type illustCommentAddParams struct {
	IllustID        uint64 `url:"illust_id,omitempty"`
	Comment         string `url:"comment,omitempty"`
	ParentCommentID int    `url:"parent_comment_id,omitempty"`
}

// IllustCommentAdd adds a comment to given illustID
func (a *AppPixivAPI) IllustCommentAdd(illustID uint64, comment string, parentCommentID int) (*IllustCommentAddResult, error) {
	data := &IllustCommentAddResult{}
	params := &illustCommentAddParams{
		IllustID:        illustID,
		Comment:         comment,
		ParentCommentID: parentCommentID,
	}
	if err := a.post(ADD, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type illustRelatedParams struct {
	IllustID      uint64   `url:"illust_id,omitempty"`
	Filter        string   `url:"filter,omitempty"`
	SeedIllustIDs []string `url:"seed_illust_ids[],omitempty,omitempty"`
}

// IllustRelated returns Related works
func (a *AppPixivAPI) IllustRelated(illustID uint64, filter string, seedIllustIDs []string) (*IllustsResponse, error) {
	data := &IllustsResponse{}
	if filter == "" {
		filter = "for_ios"
	}
	params := &illustRelatedParams{
		IllustID: illustID,
		Filter:   filter,
	}
	if seedIllustIDs != nil {
		params.SeedIllustIDs = seedIllustIDs
	}

	if err := a.request(RELATED, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type illustRecommendedParams struct {
	ContentType                  string   `url:"content_type,omitempty"`
	IncludeRankingLabel          bool     `url:"include_ranking_label,omitempty"`
	Filter                       string   `url:"filter,omitempty"`
	MaxBookmarkIDForRecommended  string   `url:"max_bookmark_id_for_recommend,omitempty"`
	MinBookmarkIDForRecentIllust string   `url:"min_bookmark_id_for_recent_illust,omitempty"`
	Offset                       int      `url:"offset,omitempty"`
	IncludeRankingIllusts        bool     `url:"include_ranking_illusts,omitempty"`
	BookmarkIllustIDs            []string `url:"bookmark_illust_ids,omitempty"`
	IncludePrivacyPolicy         string   `url:"include_privacy_policy,omitempty"`
}

// IllustRecommended Home Recommendation
//
// contentType: [illust, manga]
func (a *AppPixivAPI) IllustRecommended(contentType string, includeRankingLabel bool, filter string, maxBookmarkIDForRecommended string, minBookmarkIDForRecentIllust string, offset int, includeRankingIllusts bool, bookmarkIllustIDs []string, includePrivacyPolicy string, requireAuth bool) (*IllustRecommended, error) {
	data := &IllustRecommended{}
	params := &illustRecommendedParams{
		ContentType:                  contentType,
		IncludeRankingLabel:          includeRankingLabel,
		Filter:                       filter,
		Offset:                       offset,
		BookmarkIllustIDs:            bookmarkIllustIDs,
		IncludePrivacyPolicy:         includePrivacyPolicy,
		IncludeRankingIllusts:        includeRankingIllusts,
		MaxBookmarkIDForRecommended:  maxBookmarkIDForRecommended,
		MinBookmarkIDForRecentIllust: minBookmarkIDForRecentIllust,
	}
	if requireAuth {
		if err := a.request(RECOMMENDED, params, data, true); err != nil {
			return nil, err
		}
	} else {
		if err := a.request(RECOMMENDED_NO_LOGIN, params, data, true); err != nil {
			return nil, err
		}
	}
	return data, nil
}

type illustRankingParams struct {
	Mode   string `url:"mode,omitempty"`
	Filter string `url:"filter,omitempty"`
	Date   string `url:"date,omitempty"`
	Offset string `url:"offset,omitempty"`
}

// IllustRanking Ranking of works
//
// mode: [day, week, month, day_male, day_female, week_original, week_rookie, day_manga]
//
// date: yyyy-mm-dd
func (a *AppPixivAPI) IllustRanking(mode string, filter string, date string, offset string) (*IllustsResponse, error) {
	data := &IllustsResponse{}
	params := &illustRankingParams{
		Mode:   mode,
		Filter: filter,
		Offset: offset,
		Date:   date,
	}
	if err := a.request(RANKING, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type trendingTagsIllustParams struct {
	Filter string `url:"filter,omitempty"`
}

// TrendingTagsIllust Trend label
func (a *AppPixivAPI) TrendingTagsIllust(filter string) (*TrendingTagsIllust, error) {
	data := &TrendingTagsIllust{}
	params := &trendingTagsIllustParams{
		Filter: filter,
	}
	if err := a.request(TRENDING_TAGS, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type searchIllustParams struct {
	Word         string `url:"word,omitempty"`
	SearchTarget string `url:"search_target,omitempty"`
	Sort         string `url:"sort,omitempty"`
	Filter       string `url:"filter,omitempty"`
	Duration     string `url:"duration,omitempty"`
	Offset       int    `url:"offset,omitempty"`
}

// SearchIllust search for
//
// searchTarget - Search type
//
//	"partial_match_for_tags"  - The label part is consistent
//	"exact_match_for_tags"    - The labels are exactly the same
//	"title_and_caption"       - Title description
//
// sort: [date_desc, date_asc]
//
// duration: [within_last_day, within_last_week, within_last_month]
func (a *AppPixivAPI) SearchIllust(word string, searchTarget string, sort string, duration string, filter string, offset int) (*SearchIllustResult, error) {
	data := &SearchIllustResult{}
	params := &searchIllustParams{
		Word:         word,
		SearchTarget: searchTarget,
		Sort:         sort,
		Filter:       filter,
		Duration:     duration,
		Offset:       offset,
	}
	if err := a.request(SEARCH, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type illustBookmarkDetailParams struct {
	IllustID uint64 `url:"illust_id,omitempty"`
}

// IllustBookmarkDetail Bookmark details
func (a *AppPixivAPI) IllustBookmarkDetail(illustID uint64) (*IllustBookmarkDetail, error) {
	data := &IllustBookmarkDetail{}
	params := &illustBookmarkDetailParams{
		IllustID: illustID,
	}
	if err := a.request(BOOKMARK_DETAIL, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type illustBookmarkAddParams struct {
	IllustID uint64   `url:"illust_id,omitempty"`
	Restrict string   `url:"restrict,omitempty"`
	Tags     []string `url:"tags,omitempty"`
}

// IllustBookmarkAdd Add bookmark
func (a *AppPixivAPI) IllustBookmarkAdd(illustID uint64, restrict string, tags []string) error {
	params := illustBookmarkAddParams{
		IllustID: illustID,
		Restrict: restrict,
	}
	if tags != nil {
		params.Tags = tags
	}
	return a.post(BOOKMARK_ADD, params, nil, true)
}

type illustBookmarkDeleteParams struct {
	IllustID uint64 `url:"illust_id,omitempty"`
}

// IllustBookmarkDelete Remove bookmark
func (a *AppPixivAPI) IllustBookmarkDelete(illustID uint64) error {
	params := &illustBookmarkDeleteParams{
		IllustID: illustID,
	}
	return a.post(BOOKMARK_DELETE, params, nil, true)
}

type userBookmarkTagsIllustParams struct {
	Restrict string
	Offset   int
}

// UserBookmarkTagsIllust User favorite tag list
func (a *AppPixivAPI) UserBookmarkTagsIllust(restrict string, offset int) (*UserBookmarkTags, error) {
	data := &UserBookmarkTags{}
	params := &userBookmarkTagsIllustParams{
		Restrict: restrict,
		Offset:   offset,
	}
	if err := a.request(BOOKMARK_TAG, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type userFollowStatsParams struct {
	UserID   uint64 `url:"user_id,omitempty"`
	Restrict string `url:"restrict,omitempty"`
	Offset   int    `url:"offset,omitempty"`
}

func userFollowStats(a *AppPixivAPI, urlEnd string, userID uint64, restrict string, offset int) (*UserFollowList, error) {
	data := &UserFollowList{}
	params := &userFollowStatsParams{
		UserID:   userID,
		Restrict: restrict,
		Offset:   offset,
	}
	if err := a.request(USER+urlEnd, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

// UserFollowing Following user list
func (a *AppPixivAPI) UserFollowing(userID uint64, restrict string, offset int) (*UserFollowList, error) {
	return userFollowStats(a, "following", userID, restrict, offset)
}

// UserFollower Follower user list
func (a *AppPixivAPI) UserFollower(userID uint64, restrict string, offset int) (*UserFollowList, error) {
	return userFollowStats(a, "follower", userID, restrict, offset)
}

type userFollowPostParams struct {
	UserID   uint64 `url:"user_id,omitempty"`
	Restrict string `url:"restrict,omitempty"`
}

func userFollowPost(a *AppPixivAPI, urlEnd string, userID uint64, restrict string) error {
	params := userFollowPostParams{
		UserID:   userID,
		Restrict: restrict,
	}
	return a.post(USER_FOLLOW+urlEnd, params, nil, true)
}

// UserFollowAdd Follow users
func (a *AppPixivAPI) UserFollowAdd(userID uint64, restrict string) error {
	return userFollowPost(a, "add", userID, restrict)
}

// UserFollowDelete Unfollow users
func (a *AppPixivAPI) UserFollowDelete(userID uint64, restrict string) error {
	return userFollowPost(a, "delete", userID, restrict)
}

type userMyPixivParams struct {
	UserID uint64 `url:"user_id,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

// UserMyPixiv Users in MyPixiv
func (a *AppPixivAPI) UserMyPixiv(userID uint64, offset int) (*UserFollowList, error) {
	data := &UserFollowList{}
	params := &userMyPixivParams{
		UserID: userID,
		Offset: offset,
	}
	if err := a.request(USER_MYPIXIV, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type userListParams struct {
	UserID uint64 `url:"user_id,omitempty"`
	Filter string `url:"filter,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

// UserList Blacklisted users
func (a *AppPixivAPI) UserList(userID uint64, filter string, offset int) (*UserList, error) {
	data := &UserList{}
	params := &userListParams{
		UserID: userID,
		Filter: filter,
		Offset: offset,
	}
	if err := a.request(USER_LIST, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type ugoiraMetadataParams struct {
	IllustID uint64 `url:"illust_id,omitempty"`
}

// UgoiraMetadata Ugoira Info
func (a *AppPixivAPI) UgoiraMetadata(illustID uint64) (*UgoiraMetadata, error) {
	data := &UgoiraMetadata{}
	params := &ugoiraMetadataParams{IllustID: illustID}
	if err := a.request(METADATA, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type showcaseArticleParams struct {
	ShowcaseID string `url:"article_id,omitempty"`
}

// ShowcaseArticle Special feature details (disguised as Chrome)
func (a *AppPixivAPI) ShowcaseArticle(showcaseID string) (*ShowcaseArticle, error) {
	data := &ShowcaseArticle{}
	params := &showcaseArticleParams{
		ShowcaseID: showcaseID,
	}

	s := a.sling.New().Base(WEB_BASE + "/")
	s.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	s.Set("Referer", WEB_BASE)

	if _, err := s.Get(WEB_ARTICLE).QueryStruct(params).ReceiveSuccess(data); err != nil {
		return nil, err
	}
	return data, nil
}