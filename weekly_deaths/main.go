package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"weekly_deaths/eurostat"
	"weekly_deaths/web"
)

// DefaultPort defines a default port that the server will be started on.
const DefaultPort = 8080

//go:embed frontend/dist
var frontend embed.FS

func ensureAuthCredentialsLoaded(app web.Application) {
	if app.Auth.Username == "" {
		log.Fatal("Auth username not found in env vars.")
	}

	if app.Auth.Password == "" {
		log.Fatal("Auth password not found in env vars.")
	}
}

func initializeDataSnapshot() (eurostat.DataSnapshot, error) {
	var (
		snapshot eurostat.DataSnapshot
		err      error
	)

	env, ok := os.LookupEnv("DEPLOY_ENV")
	if !ok {
		env = "prod"
	}

	switch env {
	case "local":
		log.Println("DEPLOY_ENV=local; reading snapshot from disk.")
		path := os.Getenv("LOCAL_SNAPSHOT_PATH")
		if path == "" {
			log.Fatal("DEPLOY_ENV set to local and LOCAL_SNAPSHOT_PATH env variable is empty. Exiting.")
		}
		snapshot, err = eurostat.DataSnapshotFromPath(path)
		if err != nil {
			return snapshot, err
		}
	case "production":
		log.Println("DEPLOY_ENV=production; reading live snapshot from Eurostat.")
		snapshot, err = eurostat.DataSnapshotFromEurostat()
		if err != nil && os.Getenv("USE_S3_AS_FALLBACK") == "true" {
			log.Printf("Reading live snapshot from Eurostat failed because of: %s\n", err)
			sm, err := eurostat.NewSnapshotManager(os.Getenv("S3_BUCKET"))
			if err != nil {
				return snapshot, err
			}

			snapshot, err = sm.LatestSnapshot()
			if err != nil {
				return snapshot, err
			}
		} else {
			return snapshot, err
		}
	}

	return snapshot, nil
}

func main() {
	var port int

	flag.IntVar(&port, "port", DefaultPort, "port to start server on")
	flag.Parse()

	startTime := time.Now()
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println(".env file not found.")
	}

	snapshot, err := initializeDataSnapshot()
	if err != nil {
		log.Fatal(err)
	}

	db := eurostat.DBFromSnapshot(snapshot)
	if err != nil {
		log.Fatal(err)
	}

	app := web.Application{
		Db: db,
	}
	app.Auth.Username = os.Getenv("AUTH_USERNAME")
	app.Auth.Password = os.Getenv("AUTH_PASSWORD")
	ensureAuthCredentialsLoaded(app)

	router := app.Routes()

	stripped, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	frontendFS := http.FileServer(http.FS(stripped))
	router.Handle("/*", frontendFS)
	notFoundHtml, err := frontend.ReadFile("frontend/dist/404.html")
	if err != nil {
		log.Fatal(err)
	}

	// https://github.com/go-chi/chi/issues/155
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "frontend", "dist")
	root := http.Dir(filesDir)
	fsHandler := http.StripPrefix("/", http.FileServer(root))
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(notFoundHtml)
	}))
	router.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(fmt.Sprintf("%s", root) + r.RequestURI); os.IsNotExist(err) {
			router.NotFoundHandler().ServeHTTP(w, r)
		} else {
			fsHandler.ServeHTTP(w, r)
		}
	}))

	log.Printf("Starting the server on :%d port\n", port)
	log.Printf("Application start took %s.\n", time.Since(startTime))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatalf("starting server: %s\n", err.Error())
	}
}
