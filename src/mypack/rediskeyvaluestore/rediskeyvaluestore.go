package rediskeyvaluestore

import (
	//"bufio"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	//"os"
	//"strings"
	"time"
	//"mypack/objects"

)

const (
	ADDRESS           = "127.0.0.1:6379"
	RECENT_ITEM_LIMIT = 100
)

//Maintaining a single database connection
var (
	conn, err = GetConnection()
)

func GetConnection() (redis.Conn, error) {
	c, err := redis.Dial("tcp", ADDRESS)
	if err != nil {
		fmt.Println("FATAL    ", time.Now(), " Redis data base connection issue ")
		return nil, err
	}
	return c, nil
}

func AddKey(key, value string) {

	if _, err = conn.Do("SET", key, value); err != nil {
		log.Fatal(err)
		fmt.Println("FATAL    ", time.Now(), "Redis Key insert failure ", err)
	}
}

func GetValue(key string) []string {
	value, err := redis.Strings(conn.Do("GET", key))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)
	return value
}

func GetValueList(key string) []string {
  value, err := redis.Strings(conn.Do("LRANGE", key, 0, -1))
  if err != nil {
    log.Fatal(err)
  }
  return value
}

func RemoveKey(key string) {
	if _, err = conn.Do("DEL", key); err != nil {
		log.Fatal(err)
		fmt.Println("FATAL    ", time.Now(), "Redis Key removal failure ", err)
	}
}

func AddProjectToRecentlyViewedProjects(key string, value string, limit int64) {
	if _, err = conn.Do("LPUSH", key, value); err != nil {
		log.Fatal(err)
		fmt.Println("FATAL    ", time.Now(), "Redis Key LPUSH failure ", err)
	}
	if _, err = conn.Do("LTRIM", key, 0, limit); err != nil {
		log.Fatal(err)
		fmt.Println("FATAL    ", time.Now(), "Redis Key LTRIM failure ", err)
	}
}

func AddProjectsToFeaturedProjects(key string, value string, limit int64) {
	if _, err = conn.Do("LPUSH", key, value); err != nil {
		log.Fatal(err)
		fmt.Println("FATAL    ", time.Now(), "Redis Key LPUSH failure ", err)
	}
	if _, err = conn.Do("LTRIM", key, 0, limit); err != nil {
		log.Fatal(err)
		fmt.Println("FATAL    ", time.Now(), "Redis Key LTRIM failure ", err)
	}
}

func TestRecentProjects() {
	records := []string{"pt001", "pt002", "pt003", "pt004"}
	recentItemKey := "recently_viewed_projects"
	for _,val := range records {
		AddProjectToRecentlyViewedProjects(recentItemKey, val, 4)
	}


	keyType, _ := conn.Do("TYPE", recentItemKey)
	fmt.Println("Type", keyType)

	var results []string
	results = GetValueList(recentItemKey)
	
	for _, val := range results {
		fmt.Println(val)
	}
	//recentItems := GetValue(recentItemKey)
	length, _ := conn.Do("LLEN", recentItemKey)
	fmt.Println("Length of the list :",length)

}

/*
  Submits data to our redis instance
*/
func submitData(input []string) {

	defer conn.Close()

	conn.Send("MULTI")

	//1. delete from temp set
	conn.Send("DEL", "user-words")
	//2. store in a temp set
	conn.Send("SADD", redis.Args{}.Add("user-words").AddFlat(input)...)
	//3. take intersection of both sets
	conn.Send("SINTER", "user-words", "bad-words")

	reply, err := conn.Do("EXEC")

	if err != nil {
		fmt.Println(err)
	}

	values, _ := redis.Values(reply, nil)

	curse_words, err := redis.Strings(values[2], nil)
	if err != nil {
		fmt.Println(err)
	}

	if (len(curse_words)) > 0 {
		for _, v := range curse_words {
			fmt.Println(">>Found: ", v)
		}
	} else {
		fmt.Println(">>Nothing found")
	}
}
/*
func main() {
	TestRecentProjects()
	/*
		for {
			fmt.Println(">>Please enter some text with swear words than press Enter or \"q\" to exit")
			bio := bufio.NewReader(os.Stdin)
			line, _, _ := bio.ReadLine()

			if string(line) == "q" {
				break
			}

			terms := strings.Split(string(line), " ")
			submitData(terms)
		}
	*//*

	fmt.Println("Session Ended!")
}
*/
