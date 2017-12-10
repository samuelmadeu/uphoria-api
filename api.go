package main

import (
	// Standard library packages
	"fmt"
	"log"
	"net/http"
	"os"

	// Third party packages
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"

	// local packages
	"github.com/tanaka/uphoria-api/controllers"
)

func main() {
	// loads .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// get port number
	port := EnvGetter("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Instantiate a new router
	r := httprouter.New()

	// Get controller instances
	uc := controllers.NewUserController(getSession())    // users
	// index
	r.GET("/", Index)

	// user routes
	r.GET("/api/v1/user/:id", uc.GetUser)       // Get a user resource
	r.POST("/api/v1/user", uc.CreateUser)       // Create a new user
	r.DELETE("/api/v1/user/:id", uc.RemoveUser) // Remove an existing user
	r.GET("/api/v1/users", uc.GetAllUsers)      // get all users

	// Fire up the server
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Index - to public API
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {

	// We need this object to establish a session to our MongoDB.
	// using the .env to load the database information
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{EnvGetter("MongoDBHosts")},
		Database: EnvGetter("AuthDatabase"),
		Username: EnvGetter("AuthUserName"),
		Password: EnvGetter("AuthPassword"),
		//Timeout:  60 * time.Second,
	}

	// Connect to our local mongo
	// s = sessions
	s, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
		//log.Fatal(panic(err))
	} // Check if connection error, is mongo running?
	//defer s.Closing()

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	s.SetMode(mgo.Monotonic, true)

	fmt.Println("Session created")

	// Deliver session
	return s
}

type (
	// Server -
	Server struct {
		session *mgo.Session
	}
)

// Closing - Clean up
func (s *Server) Closing() {
	s.session.Close()
}

// EnvGetter - Fetches envs
func EnvGetter(k string) string {
	v := os.Getenv(k)
	return v
}
