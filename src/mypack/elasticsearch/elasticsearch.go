package main

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"mypack/objects"
	"os"
	"reflect"
	"time"
)

type Tweet struct {
	User     string
	Message  string
	Retweets int
}

var (
	client, err = GetElasticSearchClient()
)

func GetElasticSearchClient() (*elastic.Client, error) {
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog))
	if err != nil {
		panic(err)
		return nil, err
	}

	// Trace request and response details like this
	//client.SetTracer(log.New(os.Stdout, "", 0))

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping().Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("[", time.Now(), "]	[INFO]	[Elasticsearch]	server running ...")

	fmt.Println("[", time.Now(), "]	[INFO]	[Elasticsearch]	returned with code %d and version %s]", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("[", time.Now(), "]	[INFO]	[Elasticsearch]	version %s", esversion)
	return client, nil
}

func CreateProjectTubeIndex() error {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("projecttube").Do()
	if err != nil {
		// Handle error
		panic(err)
		return err
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("projecttube").Do()
		if err != nil {
			// Handle error
			panic(err)
			return err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return nil
}

func IndexProject(p *objects.ElasticProject) error {
	//Index a Project (using JSON serialization)
	put, err := client.Index().
		Index("projecttube").
		Type("project").
		Id("1").
		BodyJson(p).
		Do()
	if err != nil {
		// Handle error
		panic(err)
		return err
	}
	fmt.Printf("Indexed project %s to index %s, type %s\n", put.Id, put.Index, put.Type)
	return nil
}

func GetProjectByID(index int64) error {
	// Get tweet with specified ID
	get, err := client.Get().
		Index("projecttube").
		Type("project").
		Id(string(index)).
		Do()
	if err != nil {
		// Handle error
		panic(err)
		return err
	}
	if get.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get.Id, get.Version, get.Index, get.Type)
	}
	return nil
}

//Match query
func SearchProjectMatchQuery() error {
	// Search with a term query
	//It is expecting atleast one matching word from the search string to match a record. Not compared against characters
	matchQuery := elastic.NewMatchQuery("message", "y Waltz")
	searchResult, err := client.Search().
		Index("twitter").   // search in index "twitter"
		Query(&matchQuery). // specify the query
		Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do()                // execute
	if err != nil {
		panic(err)
		return err
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Match Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(Tweet)
		fmt.Printf("Match Query: Tweet by %s: %s\n", t.User, t.Message)
	}

	return nil
}

//Term query
func SearchProjectTermQuery() error {
	// Search with a term query
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("twitter").   // search in index "twitter"
		Query(&termQuery).  // specify the query
		Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do()                // execute
	if err != nil {
		// Handle error
		panic(err)
		return err
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(Tweet)
		fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits != nil {
		fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}
	return nil
}

func TestElasticSearch() {

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("twitter").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("twitter").Do()
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// Index a tweet (using JSON serialization)
	tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	put1, err := client.Index().
		Index("twitter").
		Type("tweet").
		Id("1").
		BodyJson(tweet1).
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Index a second tweet (by string)
	tweet2 := `{"user" : "olivere", "message" : "It's a Raggy Waltz"}`
	put2, err := client.Index().
		Index("twitter").
		Type("tweet").
		Id("2").
		BodyString(tweet2).
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put2.Id, put2.Index, put2.Type)

	//tweet3 := Tweet{User: "olivere", Message: "Google glass stopped in production", Retweets: 0}
	tweet3 := `{"user" : "olivere", "message" : "Google glass stopped in production"}`
	put3, err := client.Index().
		Index("twitter").
		Type("tweet").
		Id("3").
		BodyJson(tweet3).
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put3.Id, put3.Index, put3.Type)

	// Get tweet with specified ID
	get1, err := client.Get().
		Index("twitter").
		Type("tweet").
		Id("1").
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("twitter").Do()
	if err != nil {
		panic(err)
	}

	SearchProjectMatchQuery()

	/*
		// Search with a term query
		termQuery := elastic.NewTermQuery("user", "olivere")
		searchResult, err := client.Search().
			Index("twitter").   // search in index "twitter"
			Query(&termQuery).  // specify the query
			Sort("user", true). // sort by "user" field, ascending
			From(0).Size(10).   // take documents 0-9
			Pretty(true).       // pretty print request and response JSON
			Do()                // execute
		if err != nil {
			// Handle error
			panic(err)
		}

		// searchResult is of type SearchResult and returns hits, suggestions,
		// and all kinds of other information from Elasticsearch.
		fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

		// Each is a convenience function that iterates over hits in a search result.
		// It makes sure you don't need to check for nil values in the response.
		// However, it ignores errors in serialization. If you want full control
		// over iterating the hits, see below.
		var ttyp Tweet
		for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
			t := item.(Tweet)
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
		// TotalHits is another convenience function that works even when something goes wrong.
		fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

		// Here's how you iterate through results with full control over each step.
		if searchResult.Hits != nil {
			fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

			// Iterate through results
			for _, hit := range searchResult.Hits.Hits {
				// hit.Index contains the name of the index

				// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
				var t Tweet
				err := json.Unmarshal(*hit.Source, &t)
				if err != nil {
					// Deserialization failed
				}

				// Work with tweet
				fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
			}
		} else {
			// No hits
			fmt.Print("Found no tweets\n")
		}
		/*
			// Update a tweet by the update API of Elasticsearch.
			// We just increment the number of retweets.
			update, err := client.Update().Index("twitter").Type("tweet").Id("1").
				Script("ctx._source.retweets += num").
				ScriptParams(map[string]interface{}{"num": 1}).
				Upsert(map[string]interface{}{"retweets": 0}).
				Do()
			if err != nil {
				// Handle error
				panic(err)
			}

			fmt.Printf("New version of tweet %q is now %d", update.Id, update.Version)
	*/
	// ...

	// Delete an index.
	deleteIndex, err := client.DeleteIndex("twitter").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}

}

func main() {
	TestElasticSearch()
}
