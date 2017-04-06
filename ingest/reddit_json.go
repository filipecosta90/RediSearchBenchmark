package ingest

import (
	"io"

	"compress/bzip2"
	"encoding/json"

	"log"

	"github.com/RedisLabs/RediSearchBenchmark/index"
)

type redditDocument struct {
	Author     string  `json:"author"`
	Body       string  `json:"body"`
	Created    int64   `json:"created_utc"`
	Id         string  `json:"id"`
	Score      int64   `json:"score"`
	Ups        int64   `json:"ups"`
	Downs      int64   `json:"downs"`
	Subreddit  string  `json:"subreddit"`
	UvoteRatio float32 `json:"upvote_ratio"`
}

type RedditReader struct{}

func (rr *RedditReader) Read(r io.Reader, ch chan index.Document) error {
	log.Println("Reddit reader opening", r)
	bz := bzip2.NewReader(r)
	jr := json.NewDecoder(bz)

	var rd redditDocument

	//go func() {
	var err error

	for err != io.EOF {

		if err := jr.Decode(&rd); err != nil {
			log.Printf("Error decoding json: %s", err)
			break
		}
		doc := index.NewDocument(rd.Id, float32(rd.Score)).
			Set("body", rd.Body).
			Set("author", rd.Author).
			Set("sub", rd.Subreddit).
			Set("date", rd.Created).
			Set("ups", rd.Ups)

		ch <- doc
	}
	//close(ch)
	//}()
	return nil
}
