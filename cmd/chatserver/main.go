package main

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "gitlab.com/vdat/mcsvc/chat/docs"
	"gitlab.com/vdat/mcsvc/chat/migration"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/article"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/call"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/category"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/comment"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/dchat/v3"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/groups/v2"
	message_service "gitlab.com/vdat/mcsvc/chat/pkg/service/message/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/reaction"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/request"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/upload/v1"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"log"
	"net"
	"net/http"
	"os"
	_ "path"
	"path/filepath"
	"sync"
	"time"
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP request. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

// @title Swagger Chat server API
// @version 0.1
// @description This is swagger for chat server.
// @description local:	  http://localhost:5000/.
// @description staging:    http://vdat-mcsvc-chat-staging.vdatlab.com/.
// @description production: https://vdat-mcsvc-chat.vdatlab.com/.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @query.collection.format multi
// @Schemes http https
// @host localhost:5000
// @BasePath /api/v1
// @query.collection.format multi
func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go migration.StartMigration(wg)
	wg.Wait()

	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://5d512dfce00b42f0b4b2823d11dd6a59@o530089.ingest.sentry.io/5803023",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("It works!")

	//go metrics()
	database.Connect()

	//readfile
	//utils.CheckFileSocketId()

	//start broker
	go dchat.Wsbroker.Run()

	r := mux.NewRouter()

	//r.HandleFunc("/healthcheck", CheckHelthHandlr).Methods(http.MethodGet, http.MethodOptions)
	//r.HandleFunc("/chat/{idgroup}", auth.AuthenMiddleJWT(dchat.ChatHandlr)).Methods(http.MethodGet, http.MethodOptions)
	// handler
	r.HandleFunc("/message", dchat.ServeWebSocket).Methods(http.MethodGet, http.MethodOptions)
	//api
	//groups.RegisterGroupApi(r)
	userdetail.RegisterUserApi(r)
	category.NewHandler(r)
	article.NewHandler(r)
	comment.NewHandler(r)
	reaction.NewHandler(r)
	upload.NewHandler(r)
	call.NewHandler(r)
	request.NewHandler(r)

	//apiv2
	groups.NewHandler(r)

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Handler web app
	spa := spaHandler{staticPath: "public", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	r.Use(mux.CORSMethodMiddleware(r))

	srv := &http.Server{
		Handler: r,
		Addr:    ":5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	CheckSchedule()
	fmt.Println("server starting")
	log.Fatal(srv.ListenAndServe())

}

func metrics() {
	// The debug listener mounts the http.DefaultServeMux, and serves up
	// stuff like the Prometheus metrics route, the Go debug and profiling
	// routes, and so on.
	debugListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("transport", "debug/HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}

	defer debugListener.Close()

	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.Serve(debugListener, http.DefaultServeMux))
}

func CheckSchedule() {
	s := gocron.NewScheduler(time.UTC)
	_, _ = s.Every(45).Minute().Do(func() {
		//sentry.CaptureMessage(time.Now().Format("01-02-2021 15:04:05"))
		fmt.Println(time.Now().Format("01-02-2021 15:04:05"))
		timeoutContext := time.Duration(2) * time.Second
		messRepo := message_service.NewRepoImpl(database.DB)
		messService := message_service.NewServiceImpl(messRepo, timeoutContext)
		GroupRepo := groups.NewRepoImpl(database.DB)
		userRepo := userdetail.NewRepoImpl(database.DB)
		userService := userdetail.NewServiceImpl(userRepo, timeoutContext)
		groupService := groups.NewServiceImpl(GroupRepo, userService, timeoutContext, messService)
		err := groupService.DeleteGroupNeedDelete(context.Background())
		if err != nil {
			sentry.CaptureException(err)
			fmt.Println(err)
		}
	})
	s.StartAsync()
}
