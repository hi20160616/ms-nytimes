package fetcher

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/hi20160616/exhtml"
	"github.com/pkg/errors"
)

func TestFetchTitle(t *testing.T) {
	tests := []struct {
		url   string
		title string
	}{
		{"https://cn.nytimes.com/china/20211011/china-child-custody-abductions/?utm_source=RSS", "监护权争夺战：那些绑架自己孩子的中国家长"},
		{"https://cn.nytimes.com/asia-pacific/20201201/north-korea-pandemic-china/", "疫情进一步冲击对华贸易，朝鲜经济困境加深"},
	}
	for _, tc := range tests {
		a := NewArticle()
		u, err := url.Parse(tc.url)
		if err != nil {
			t.Error(err)
		}
		a.U = u
		// Dail
		a.raw, a.doc, err = exhtml.GetRawAndDoc(a.U, timeout)
		if err != nil {
			t.Error(err)
		}
		got, err := a.fetchTitle()
		if err != nil {
			if !errors.Is(err, ErrTimeOverDays) {
				t.Error(err)
			} else {
				fmt.Println("ignore pass test: ", tc.url)
			}
		} else {
			if tc.title != got {
				t.Errorf("\nwant: %s\n got: %s", tc.title, got)
			}
		}
	}

}

func TestFetchUpdateTime(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://cn.nytimes.com/china/20211011/china-child-custody-abductions/?utm_source=RSS", "2021-10-11 17:41:03 +0800 UTC"},
		{"https://cn.nytimes.com/asia-pacific/20201201/north-korea-pandemic-china/", "2020-12-01 19:04:20 +0800 UTC"},
	}
	var err error
	for _, tc := range tests {
		a := NewArticle()
		a.U, err = url.Parse(tc.url)
		if err != nil {
			t.Error(err)
		}
		// Dail
		a.raw, a.doc, err = exhtml.GetRawAndDoc(a.U, timeout)
		if err != nil {
			t.Error(err)
		}
		tt, err := a.fetchUpdateTime()
		if err != nil {
			if !errors.Is(err, ErrTimeOverDays) {
				t.Error(err)
			}
		}
		ttt := tt.AsTime()
		got := shanghai(ttt)
		if got.String() != tc.want {
			t.Errorf("\nwant: %s\n got: %s", tc.want, got.String())
		}
	}
}

func TestFetchContent(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://cn.nytimes.com/china/20211011/china-child-custody-abductions/?utm_source=RSS", "监护权争夺战：那些绑架自己孩子的中国家长"},
		// {"https://cn.nytimes.com/asia-pacific/20201201/north-korea-pandemic-china/", "疫情进一步冲击对华贸易，朝鲜经济困境加深"},
	}
	var err error

	for _, tc := range tests {
		a := NewArticle()
		a.U, err = url.Parse(tc.url)
		if err != nil {
			t.Error(err)
		}
		// Dail
		a.raw, a.doc, err = exhtml.GetRawAndDoc(a.U, timeout)
		if err != nil {
			t.Error(err)
		}
		c, err := a.fetchContent()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(c)
	}
}

func TestFetchArticle(t *testing.T) {
	tests := []struct {
		url string
		err error
	}{
		{"https://cn.nytimes.com/china/20211011/china-child-custody-abductions/?utm_source=RSS", ErrTimeOverDays},
		{"https://cn.nytimes.com/asia-pacific/20201201/north-korea-pandemic-china/", nil},
	}
	for _, tc := range tests {
		a := NewArticle()
		a, err := a.fetchArticle(tc.url)
		if err != nil {
			if !errors.Is(err, ErrTimeOverDays) {
				t.Error(err)
			} else {
				fmt.Println("ignore old news pass test: ", tc.url)
			}
		} else {
			fmt.Println("pass test: ", a.Content)
		}
	}
}
