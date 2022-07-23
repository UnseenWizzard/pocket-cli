package articles

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/UnseenWizzard/pocket-cli/pkg/api/retrieve"
)

var fullTestArticle retrieve.Article = retrieve.Article{
	Id:         "ID",
	Title:      "A Title",
	GivenTitle: "A Given Title",
	Url:        "https://www.url.com",
	Excerpt:    "A Summary",
	ReadTime:   42,
}

func Test_beautifyTitle(t *testing.T) {
	tests := []struct {
		name         string
		inputArticle retrieve.Article
		want         string
	}{
		{
			"ReturnsTileIfPresent",
			fullTestArticle,
			fullTestArticle.Title,
		},
		{
			"ReturnsGivenTileIfNoTitlePresent",
			retrieve.Article{
				Id:         "ID",
				GivenTitle: "A Given Title",
				Url:        "https://www.url.com",
				Excerpt:    "A Summary",
				ReadTime:   42,
			},
			"A Given Title",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := beautifyTitle(tt.inputArticle); got != tt.want {
				t.Errorf("beautifyTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beautifyReadTime(t *testing.T) {
	tests := []struct {
		name         string
		inputArticle retrieve.Article
		want         string
	}{
		{
			"ReturnsReadTimeIfPresent",
			fullTestArticle,
			"42 min",
		},
		{
			"ReturnsQuestionMarkIfNoReadTime",
			retrieve.Article{
				Id:         "ID",
				Title:      "A Title",
				GivenTitle: "A Given Title",
				Url:        "https://www.url.com",
				Excerpt:    "A Summary",
			},
			"?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := beautifyReadTime(tt.inputArticle); got != tt.want {
				t.Errorf("beautifyReadTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beautifyExcerpt(t *testing.T) {
	tests := []struct {
		name         string
		inputArticle retrieve.Article
		want         string
	}{
		{
			"ReturnsFullExcerptIfShorterThan120Chars",
			fullTestArticle,
			fullTestArticle.Excerpt,
		},
		{
			"ReturnsInformationIfExcerptIsMissing",
			retrieve.Article{
				Id:         "ID",
				Title:      "A Title",
				GivenTitle: "A Given Title",
				Url:        "https://www.url.com",
			},
			"[No excerpt available]",
		},
		{
			"ReturnsInformationIfExcerptIsEmpty",
			retrieve.Article{
				Id:         "ID",
				Title:      "A Title",
				GivenTitle: "A Given Title",
				Url:        "https://www.url.com",
				Excerpt:    "",
			},
			"[No excerpt available]",
		},
		{
			"ReturnsInformationIfExcerptIsBlank",
			retrieve.Article{
				Id:         "ID",
				Title:      "A Title",
				GivenTitle: "A Given Title",
				Url:        "https://www.url.com",
				Excerpt:    "        ",
			},
			"[No excerpt available]",
		},
		{
			"ReturnsExcerptShortenedTo120Chars",
			retrieve.Article{
				Id:         "ID",
				Title:      "A Title",
				GivenTitle: "A Given Title",
				Url:        "https://www.url.com",
				Excerpt:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse consequat tellus gravida ligula ultricies, vel eleifend neque suscipit efficitur.",
			},
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse consequat tellus gravida ligula ultricies, vel e...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := beautifyExcerpt(tt.inputArticle)

			if len(got) > 120 {
				t.Errorf("beautifyExcerpt() did not shorten excerpt to 120 chars")
			}

			if got != tt.want {
				t.Errorf("beautifyExcerpt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fetch(t *testing.T) {
	type args struct {
		consumerKey string
		accessToken string
		count       int
		offset      int
		retrieveFn  func(string, string, int, int) (retrieve.ResponsePayload, error)
	}
	tests := []struct {
		name    string
		args    args
		want    []Article
		wantErr bool
	}{
		{
			name: "ReturnsRetrievedArticles",
			args: args{
				"key",
				"token",
				10,
				0,
				func(key, token string, count, offset int) (retrieve.ResponsePayload, error) {
					res := retrieve.ResponsePayload{
						Status: 1,
						List: map[string]retrieve.Article{
							fullTestArticle.Id: fullTestArticle,
						},
					}
					return res, nil
				},
			},
			want: []Article{
				{
					Id:       fullTestArticle.Id,
					Title:    fullTestArticle.Title,
					Excerpt:  fullTestArticle.Excerpt,
					ReadTime: fmt.Sprintf("%v min", fullTestArticle.ReadTime),
					Url:      fullTestArticle.Url,
				},
			},
			wantErr: false,
		},
		{
			name: "ReturnsErrorOnRetrieveError",
			args: args{
				"key",
				"token",
				10,
				0,
				func(key, token string, count, offset int) (retrieve.ResponsePayload, error) {
					return retrieve.ResponsePayload{}, fmt.Errorf("test error")
				},
			},
			want:    []Article{},
			wantErr: true,
		},
		{
			name: "ReturnsErrorOnRetrieveFailureStatus",
			args: args{
				"key",
				"token",
				10,
				0,
				func(key, token string, count, offset int) (retrieve.ResponsePayload, error) {
					res := retrieve.ResponsePayload{
						Status: 0,
						List:   map[string]retrieve.Article{},
					}
					return res, nil
				},
			},
			want:    []Article{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetch(tt.args.consumerKey, tt.args.accessToken, tt.args.count, tt.args.offset, tt.args.retrieveFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}
