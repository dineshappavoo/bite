package mysql_data_access

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"mypack/objects"
)

func GetDBConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@/projecttube")
	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Data base connection issue ")
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	return db, err
}

func TestConnection() {

	db, err := GetDBConnection()
	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO meta_tags VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in prepare statement for Project insert ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT tag_id FROM meta_tags WHERE metatag_id = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	//_, err = stmtIns.Exec("LINUX", "UBUNTU") // Insert tuples (i, i^2)
	/*if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}*/

	var tag_id string // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow("LINUX").Scan(&tag_id) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("Tag ID for Linux is : %s\n", tag_id)
}

func AddProject(p *objects.Project) error{
	//Get connection object
	db, err := GetDBConnection()
	defer db.Close()
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO project VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)") // ? = placeholder
	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in prepare statement for Project insert ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(p.ProjectId, p.ProjectName, p.Version, p.RecommendationPoints, p.VideoId, p.SlideId, p.GitId, p.DiscussionId, p.DateCreated, p.LastUpdated, p.RefURL, p.DownloadURL) // Insert tuples

	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in insert to Project ")
		panic(err.Error()) // proper error handling instead of panic in your app
		return err
	}
	return nil
}

func UpdateProject(p *objects.Project) {
	//Get connection object
	db, err := GetDBConnection()
	defer db.Close()
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("UPDATE project SET project_name=?, version=?, video_id=?, slide_id=?, git_id=?, last_updated=?, ref_url=?, download_url=? WHERE project_id=?") // ? = placeholder
	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in prepare statement for Project update ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(p.ProjectName, p.Version, p.VideoId, p.SlideId, p.GitId, fmt.Sprint(time.Now()), p.RefURL, p.DownloadURL, p.ProjectId) // Insert tuples

	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in update to Project ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}else{
		fmt.Println("INFO	", time.Now(), "Project Updated successfully")
	}
}

func AddUser(u *objects.User) {
	//Get connection object
	db, err := GetDBConnection()
	defer db.Close()
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO user VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)") // ? = placeholder
	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in prepare statement for User insert ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()                                                                                                                               // Close the statement when we leave main() / the program terminates
	_, err = stmtIns.Exec(u.UserId, u.UserName, u.RealName, u.Age, u.Score, u.DateJoined, u.WebsiteURL, u.Biography, u.ProfileViews, u.Email, u.GitURL) // Insert tuples

	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in insert to User ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

}

func AddTag(t *objects.Tag) {
	//Get connection object
	db, err := GetDBConnection()
	defer db.Close()
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO tag VALUES( ?, ?, ?)")
	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in prepare statement for Tag insert ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtIns.Close()                                  // Close the statement when we leave main() / the program terminates
	_, err = stmtIns.Exec(t.TagId, t.TagName, t.MetaTagId) // Insert tuples

	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in insert to Tag ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

}

func AddVideo(v *objects.Video) {
	//Get connection object
	db, err := GetDBConnection()
	defer db.Close()
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO video VALUES( ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in prepare statement for Video insert ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtIns.Close()                                                      // Close the statement when we leave main() / the program terminates
	_, err = stmtIns.Exec(v.VideoId, v.VideoURL, v.VideoEmbedURL, v.VideoHits) // Insert tuples

	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in insert to Video ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

}

func AddSlide(s *objects.Slide) {
	//Get connection object
	db, err := GetDBConnection()
	defer db.Close()
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO video VALUES( ?, ?, ?)")
	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in prepare statement for Slide insert ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtIns.Close()                                         // Close the statement when we leave main() / the program terminates
	_, err = stmtIns.Exec(s.SlideId, s.SlideURL, s.SlideEmbedURL) // Insert tuples

	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in insert to Slide ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

func AddQuestion(q *objects.Question) {
	//Get connection object
	db, err := GetDBConnection()
	defer db.Close()
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO video VALUES( ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in prepare statement for Question insert ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtIns.Close()                                                                        // Close the statement when we leave main() / the program terminates
	_, err = stmtIns.Exec(q.QuestionId, q.QuestionText, q.Upvotes, q.DateCreated, q.LastUpdated) // Insert tuples

	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in insert to Question ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

}

func AddAnswer(a *objects.Answer) {
	//Get connection object
	db, err := GetDBConnection()
	defer db.Close()
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO video VALUES( ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in prepare statement for Answer insert ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtIns.Close()                                                                                        // Close the statement when we leave main() / the program terminates
	_, err = stmtIns.Exec(a.AnswerId, a.QuestionId, a.Upvotes, a.UserId, a.UserId, a.DateCreated, a.LastUpdated) // Insert tuples

	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Error in insert to Answer ")
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

func GetProject(projectID string) (*objects.Project, error){
	//Get connection object
	db, err := GetDBConnection()
	defer db.Close()

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT * FROM project WHERE project_id = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	var projectId string
	var projectName string
	var version string
	var recommendationPoints int
	var videoId string
	var slideId string
	var gitId string
	var discussionId string
	var dateCreated string
	var lastUpdated string
	var refURL string
	var downloadURL string

	// Query the project
	err = stmtOut.QueryRow(projectID).Scan(&projectId, &projectName, &version, &recommendationPoints, &videoId, &slideId, &gitId, &discussionId, &dateCreated, &lastUpdated, &refURL, &downloadURL) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}
	return &objects.Project{ProjectId: projectId, ProjectName: projectName, Version: version, RecommendationPoints: recommendationPoints, VideoId: videoId,
		SlideId: slideId, GitId: gitId, DiscussionId: discussionId, DateCreated: dateCreated, LastUpdated: lastUpdated, RefURL: refURL, DownloadURL: downloadURL}, nil
}

func GetUser() {

}

func GetVideo() {

}

func GetSlide() {

}

func GetQuestion() {

}

func GetAnswer() {

}

/*
func main() {
	//TestConnection()
	//AddProject(&Project{ProjectId: "pt005", ProjectName: "Ubuntu", Version: "9.1.0", RecommendationPoints: 10, VideoId: "GlZC4Jwf3xQ",
		//SlideId: "34228704", GitId: "github.com/ubuntu", DiscussionId: "d001", DateCreated: fmt.Sprint(time.Now()), LastUpdated: fmt.Sprint(time.Now()), RefURL: "github.com/dineshappavoo/projectcube", DownloadURL: "ubuntu.com"})

	p, _ := GetProject("pt002")
	fmt.Println("Project INFO     : ", p)
}*/
