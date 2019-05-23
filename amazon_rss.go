package feeds

// rss support
// validation done according to spec here:
//    http://cyber.law.harvard.edu/rss/rss.html

import (
	"encoding/xml"
	"fmt"
	"time"
)

// AmazonRssFeedXml is private wrapper around the RssFeed to provide the <rss>..</rss> xml
type AmazonRssFeedXml struct {
	XMLName             xml.Name `xml:"rss"`
	Version             string   `xml:"version,attr"`
	ContentNamespace    string   `xml:"xmlns:content,attr"`
	DublinCoreNamespace string   `xml:"xmlns:dc,attr"`
	AmazonNamespace     string   `xml:"xmlns:amzn,attr"`
	Channel             *AmazonRssFeed
}

// AmazonRssFeed has amazon-specific feed elements
type AmazonRssFeed struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`       // required
	Link           string   `xml:"link"`        // required
	Description    string   `xml:"description"` // required
	Language       string   `xml:"language,omitempty"`
	Copyright      string   `xml:"copyright,omitempty"`
	ManagingEditor string   `xml:"managingEditor,omitempty"` // Author used
	WebMaster      string   `xml:"webMaster,omitempty"`
	PubDate        string   `xml:"pubDate,omitempty"`       // created or updated
	LastBuildDate  string   `xml:"lastBuildDate,omitempty"` // updated used
	Category       string   `xml:"category,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	Docs           string   `xml:"docs,omitempty"`
	Cloud          string   `xml:"cloud,omitempty"`
	Ttl            int      `xml:"ttl,omitempty"`
	Rating         string   `xml:"rating,omitempty"`
	SkipHours      string   `xml:"skipHours,omitempty"`
	SkipDays       string   `xml:"skipDays,omitempty"`
	AmznRssVersion float32  `xml:"amzn:rssVersion,omitempty"`
	Image          *RssImage
	TextInput      *RssTextInput
	Items          []*AmazonRssItem `xml:"item"`
}

// AmazonRssItem has amazon-specific item elements
type AmazonRssItem struct {
	XMLName      xml.Name `xml:"item"`
	Title        string   `xml:"title"`       // required
	Link         string   `xml:"link"`        // required
	Description  string   `xml:"description"` // required
	Content      *RssContent
	Author       string `xml:"author,omitempty"`
	Category     string `xml:"category,omitempty"`
	Comments     string `xml:"comments,omitempty"`
	Enclosure    *RssEnclosure
	Guid         string           `xml:"guid,omitempty"`    // Id used
	PubDate      string           `xml:"pubDate,omitempty"` // created or updated
	Source       string           `xml:"source,omitempty"`
	Creator      string           `xml:"dc:creator,omitempty"`
	HeroImage    string           `xml:"amzn:heroImage,omitempty"`
	IntroText    string           `xml:"amzn:introText,omitempty"`
	IndexContent string           `xml:"amzn:indexContent,omitempty"`
	Products     []*AmazonProduct `xml:"amzn:product"`
}

// AmazonProduct has the slide-specific fields
type AmazonProduct struct {
	URL      string `xml:"amzn:productURL"`
	Headline string `xml:"amzn:productHeadline"`
	Award    string `xml:"amzn:award"`
	Summary  string `xml:"amzn:productSummary"`
}

type AmazonRss struct {
	*Feed
}

// create a new AmazonRssItem with a generic Item struct's data
func newAmazonRssItem(i *Item) *AmazonRssItem {
	item := &AmazonRssItem{
		Title:        i.Title,
		Link:         i.Link.Href,
		Description:  i.Description,
		Guid:         i.Id,
		PubDate:      anyTimeFormat(time.RFC1123Z, i.Created, i.Updated),
		HeroImage:    "POST THUMBNAIL (Prefer 2x1 at least 1000px wide)",
		IntroText:    "META DESCRIPTION",
		IndexContent: "True",
	}
	if len(i.Content) > 0 {
		item.Content = &RssContent{Content: i.Content}
	}
	if i.Source != nil {
		item.Source = i.Source.Href
	}

	// Define a closure
	if i.Enclosure != nil && i.Enclosure.Type != "" && i.Enclosure.Length != "" {
		item.Enclosure = &RssEnclosure{Url: i.Enclosure.Url, Type: i.Enclosure.Type, Length: i.Enclosure.Length}
	}

	if i.Author != nil {
		item.Author = i.Author.Name
	}
	return item
}

// AmazonRssFeed will create a new AmazonRssFeed with a generic Feed struct's data
func (r *AmazonRss) AmazonRssFeed() *AmazonRssFeed {
	pub := anyTimeFormat(time.RFC1123Z, r.Created, r.Updated)
	build := anyTimeFormat(time.RFC1123Z, r.Updated)
	author := ""
	if r.Author != nil {
		author = r.Author.Email
		if len(r.Author.Name) > 0 {
			author = fmt.Sprintf("%s (%s)", r.Author.Email, r.Author.Name)
		}
	}

	var image *RssImage
	if r.Image != nil {
		image = &RssImage{Url: r.Image.Url, Title: r.Image.Title, Link: r.Image.Link, Width: r.Image.Width, Height: r.Image.Height}
	}

	channel := &AmazonRssFeed{
		Title:          r.Title,
		Link:           r.Link.Href,
		Description:    r.Description,
		ManagingEditor: author,
		PubDate:        pub,
		LastBuildDate:  build,
		Copyright:      r.Copyright,
		Image:          image,
		AmznRssVersion: 1.0,
	}
	for _, i := range r.Items {
		channel.Items = append(channel.Items, newAmazonRssItem(i))
	}
	return channel
}

// FeedXml returns an XML-Ready object for an Rss object
func (r *AmazonRss) FeedXml() interface{} {
	// only generate version 2.0 feeds for now
	return r.AmazonRssFeed().FeedXml()

}

// FeedXml returns an XML-ready object for an RssFeed object
func (r *AmazonRssFeed) FeedXml() interface{} {
	return &AmazonRssFeedXml{
		Version:             "2.0",
		Channel:             r,
		ContentNamespace:    "http://purl.org/rss/1.0/modules/content/",
		DublinCoreNamespace: "http://purl.org/dc/elements/1.1/",
		AmazonNamespace:     "https://amazon.com/ospublishing/1.0/",
	}
}
