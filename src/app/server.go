// Copyright 2015 Dinesh Appavoo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
	//"bytes"
	//"mypack/cassandra_data_access"
	"mypack/mysql_data_access"
	"mypack/objects"
	//"mypack/rediskeyvaluestore"
	//"mypack/searchengine"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
	pID = 1
)

/*
type Page struct {
	Title string
	Body  []byte
}
*/

/*
type UserInfo struct {
	Title string
	UserId   string
	UserName string
}*/

func save(p *objects.Project) error {
	//filename := "projects/" + p.Title + ".txt"
	//return ioutil.WriteFile(filename, p.Body, 0600)
	mysql_data_access.UpdateProject(p)
	return nil
}

func addProject(p *objects.Project) error {
	err := mysql_data_access.AddProject(p)
	if err != nil {
		return err
	}
	return nil
}

//load the project from the database
func loadPage(projectId string) (*objects.Project, error) {
	//filename := "projects/" + title + ".txt"
	p, err := mysql_data_access.GetProject(projectId)
	//body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//return &Page{Title: title, Body: body}, nil
	return p, nil
}

//Home page handler
//Hard coding the user name
func homeHandler(w http.ResponseWriter, r *http.Request, projectID string) {
	p := &objects.Project{ProjectId: projectID, ProjectName: "Ubuntu", Version: "9.1.0", RecommendationPoints: 10, VideoId: "GlZC4Jwf3xQ", SlideId: "34228704", GitId: "github.com/ubuntu", DiscussionId: "d001", DateCreated: fmt.Sprint(time.Now()), LastUpdated: fmt.Sprint(time.Now()), RefURL: "github.com/dineshappavoo/projectcube", DownloadURL: "ubuntu.com"}
	//p := &UserInfo{Title: "Project Tube",UserId: "dxa132330", UserName: "Dinesh Appavoo"}
	renderTemplate(w, "home", p)
}

//Search project handler
func searchHandler(w http.ResponseWriter, r *http.Request, projectID string) {
	fmt.Println("method:", r.Method) //get request method
	r.ParseForm()
	if r.Method == "GET" {
		form_data := r.FormValue("form_data")
		fmt.Println("Form Data : ", form_data)
		fmt.Println("Form Data  1: ", r.Form)
		/*for _, val := range r.FormValue("search_string") {
			fmt.Println("Search string: ", val)
		}*/

	} else {
		r.ParseForm()
		fmt.Println("Search string:", r.FormValue("search_string"))
	}
	//p := &UserInfo{Title: "Project Tube",UserId: "dxa132330", UserName: "Dinesh Appavoo"}
	p := &objects.Project{ProjectId: projectID, ProjectName: "Ubuntu", Version: "9.1.0", RecommendationPoints: 10, VideoId: "GlZC4Jwf3xQ", SlideId: "34228704", GitId: "github.com/ubuntu", DiscussionId: "d001", DateCreated: fmt.Sprint(time.Now()), LastUpdated: fmt.Sprint(time.Now()), RefURL: "github.com/dineshappavoo/projectcube", DownloadURL: "ubuntu.com"}

	renderTemplate(w, "searchproject", p)
}

//View project handler
func viewHandler(w http.ResponseWriter, r *http.Request, projectID string) {
	p, err := loadPage(projectID)
	if err != nil {
		http.Redirect(w, r, "/editproject/"+projectID, http.StatusFound)
		return
	}
	renderTemplate(w, "viewproject", p)
}

//Edit project handler
func editHandler(w http.ResponseWriter, r *http.Request, projectID string) {
	p, _ := loadPage(projectID)

	/*
		if err != nil {
			p = &Page{Title: title}
		}*/
	renderTemplate(w, "editproject", p)
}

//To be consistent with renderTemplate function we use empty string even though it is not required
func newProjectHandler(w http.ResponseWriter, r *http.Request, projectID string) {
	p, _ := loadPage(projectID)

	/*
		if err != nil {
			p = &Page{Title: title}
		}*/
	renderTemplate(w, "newproject", p)
}

//add project handler
//To be consistent with renderTemplate function we use empty string even though it is not required
func addProjectHandler(w http.ResponseWriter, r *http.Request, empty string) {
	projectName := r.FormValue("project_name")
	version := r.FormValue("project_version")
	//recommendationPoints, _ := strconv.Atoi(r.FormValue("recommendation_points"))
	videoId := r.FormValue("youtube_url")
	slideId := r.FormValue("slideshare_url")
	gitId := r.FormValue("github_url")
	refURL := r.FormValue("ref_url")
	downloadURL := r.FormValue("download_url")
	projectID := "project0003"
	fmt.Println("Add Project  ID : ",projectID)
	discussionId := projectID + "disc"
	pID = pID + 1

	p := &objects.Project{ProjectId: projectID, ProjectName: projectName, Version: version, RecommendationPoints: 5, VideoId: videoId, SlideId: slideId, GitId: gitId, DiscussionId: discussionId, DateCreated: fmt.Sprint(time.Now()), LastUpdated: fmt.Sprint(time.Now()), RefURL: refURL, DownloadURL: downloadURL}
	err := addProject(p)
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/viewproject/"+projectID, http.StatusFound)
}

//Save project handler
func saveHandler(w http.ResponseWriter, r *http.Request, projectId string) {
	projectName := r.FormValue("project_name")
	version := r.FormValue("project_version")
	recommendationPoints, _ := strconv.Atoi(r.FormValue("recommendation_points"))
	videoId := r.FormValue("youtube_url")
	slideId := r.FormValue("slideshare_url")
	gitId := r.FormValue("github_url")
	refURL := r.FormValue("ref_url")
	downloadURL := r.FormValue("download_url")

	p := &objects.Project{ProjectId: projectId, ProjectName: projectName, Version: version, RecommendationPoints: recommendationPoints, VideoId: videoId,
		SlideId: slideId, GitId: gitId, DiscussionId: "", DateCreated: "", LastUpdated: fmt.Sprint(time.Now()), RefURL: refURL, DownloadURL: downloadURL}

	err := save(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/viewproject/"+projectId, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("home.html", "editproject.html", "viewproject.html", "searchproject.html", "header.html", "footer.html", "newproject.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {

	//If you use variables other than the struct u r passing as p, then "multiple response.WriteHeader calls" error may occur. Make sure you pass
	//all variables in the struct even they are in the header.html embedded
	if err := templates.ExecuteTemplate(w, tmpl+".html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	/*
		if err := templates.ExecuteTemplate(w, tmpl+".html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			if err1 := templates.ExecuteTemplate(w, tmpl+".html", p); err1 != nil {
				log.Printf("Layout.ExecuteTemplate: %v", err.Error())
			}
		}
	*/
	/*buf := new(bytes.Buffer)
	_,err := buf.Write([]byte(fmt.Sprintf("%v", p)))
	if err != nil {
		// you can use http.Error here, no response has been written yet
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _,err := buf.WriteTo(w); err != nil {
		log.Printf("WriteTo: %v", err)
		// you can not use http.Error here anymore. So just log the message (e.g. "broken pipe...")
	}

	*/
}

//URL validation
var validPath = regexp.MustCompile("^/(home|editfood|savefood|viewfood|searchfood|newfood|addfood)/(|[a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	fmt.Println("INFO	", time.Now(), "SERVER STARTED")

	flag.Parse()
	//TestConn()
	http.HandleFunc("/home/", makeHandler(homeHandler))
	http.HandleFunc("/viewfood/", makeHandler(viewHandler))
	http.HandleFunc("/editfood/", makeHandler(editHandler))
	http.HandleFunc("/savefood/", makeHandler(saveHandler))
	http.HandleFunc("/searchfood/", makeHandler(searchHandler))
	http.HandleFunc("/newfood/", makeHandler(newProjectHandler))
	http.HandleFunc("/addfood/", makeHandler(addProjectHandler))
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	//http.Handle("/templ/", http.StripPrefix("/templ/", http.FileServer(http.Dir("templ"))))

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":8080", nil)
}
