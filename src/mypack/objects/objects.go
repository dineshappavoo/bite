package objects

//All Variable name should be in upper case to show in GUI otherwise we will get multiple write error
import (
	//"fmt"
)

type SearchProject struct {
	ProjectId string
	ProjectName string
	Version string
	Description string
	Views int
	RecommendationPoints int
	MetaTag string
}

type Project struct {
	ProjectId string
	ProjectName string
	Version string
	RecommendationPoints int
	VideoId string
	SlideId string
	GitId string
	DiscussionId string
	DateCreated string
	LastUpdated string
	RefURL string
	DownloadURL string
}

type ElasticProject struct {
	ProjectName string
	Version string
	OperatingSystem string
	OSVersion string 
	GitURL string
	RefURL string
	DownloadURL string
}

type User struct {
	UserId string
    UserName string
    RealName string
    Age int
    Score int
    DateJoined string
    WebsiteURL string
    Biography string
    ProfileViews int64
    Email string
    GitURL string
}

type Tag struct {
	TagId string
	TagName string
	MetaTagId string
}

type Video struct {
	VideoId string
	VideoURL string
	VideoEmbedURL string
	VideoHits string
}

type Slide struct {
	SlideId string
	SlideURL string
	SlideEmbedURL string
}

type GIT struct {
	GitId string
	GitURl string
	Stars int64
}

type Discussion struct {
	DiscussionId string
	QuestionId string
	AnswerId string
	UserId string
}

type Question struct {
	QuestionId string
	QuestionText string
	Upvotes int64
	DateCreated string
	LastUpdated string
}

type Answer struct {
	AnswerId string
	QuestionId string
	AnswerText string
	Upvotes int64
	UserId string
	DateCreated string
	LastUpdated string
}





