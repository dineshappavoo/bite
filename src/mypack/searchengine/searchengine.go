package searchengine

/*Algorithm GET_MATCHING_PROJECTS(searchString, dbRecords):
1.Split the strings into words search_word_list
2.Get the stem for each word using porter stemmer algorithm
3.Append the stemmed word to search_word_list if stemmed_word != original_word
4.for each record in records
	for each word search_word_list
		if record contains word
			add project_id to project_list
5.return project_list

*/
import (
	"fmt"
	"github.com/reiver/go-porterstemmer"
	"strings"
	//"mypack/objects"

)

func doStemming(searchStr string) string{
	stem := porterstemmer.StemString(searchStr)
	//fmt.Printf("The word [%s] has the stem [%s].\n", searchStr, stem)
	return stem
}

func SplitStringByWords(searchString string) []string {
	words := strings.Fields(searchString)
	return words
}

func getStemmedString(searchString string) []string {
	words := SplitStringByWords(searchString)
	var StemmedString []string
	for i := 0; i < len(words); i++ {
		stemWord := doStemming(words[i])
		StemmedString = append(StemmedString, words[i])
		if stemWord != words[i] {
			StemmedString = append(StemmedString, stemWord)
		}
	}

	for i := 0; i < len(StemmedString); i++ {
		//fmt.Println(StemmedString[i])
	}
	return StemmedString
}

func GetMatchingProjects(searchString string, records []string) map[string]string {

	var pMap = make(map[string]string)

	stemmedString := getStemmedString(searchString)
	for i := 0; i < len(stemmedString); i++ {
		for j := 0; j < len(records); j++ {
			isExist := strings.Contains(records[j], stemmedString[i])
			if isExist {
				key := records[j][0:5]
				//fmt.Println(stemmedString[i],key)
				if _,ok := pMap[key]; ok == false {
					pMap[key] = records[j]
				}
			}
		}
	}
	return pMap
}
//For testing
func main() {
	searchString := "project ubuntu how to install  and use "
	records := []string{"pt001 This project show to install ubuntu","pt002 gocql library is a simple go module devloped for cassandra ", "pt003 go-port stemmer is the library for stemming ","pt004 project-tube stores the project information and how to use the same"}

	fmt.Println("Search String : ", searchString)
	pMap := GetMatchingProjects(searchString, records)
	for key,_ := range pMap {
		fmt.Println(key, " : ", pMap[key])
	}
}	
