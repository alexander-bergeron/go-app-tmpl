package handler

import (
	"context"
	"embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	userpb "github.com/alexander-bergeron/go-app-tmpl/gen/go/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Embed all template files
//
//go:embed views/*.tmpl
var views embed.FS

type PageHandler struct {
	templates *template.Template
	conn      *grpc.ClientConn
}

func NewPageHandler(conn *grpc.ClientConn) (*PageHandler, error) {
	log.Println("Initializing Listing Handler")

	tmpl, err := template.ParseFS(views, "views/*.tmpl")
	if err != nil {
		return nil, err
	}

	return &PageHandler{
		templates: tmpl,
		conn:      conn,
	}, nil
}

// RegisterRoutes registers all todo-related routes
func (h *PageHandler) RegisterRoutes(mux *http.ServeMux) {
	log.Println("Register Routes")
	mux.HandleFunc("/", h.handleIndex)
	mux.HandleFunc("/users", h.handleUsers)
	mux.HandleFunc("/create-user", h.handleCreateUser)
}

func (h *PageHandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	err := h.templates.ExecuteTemplate(w, "index.tmpl", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PageHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	client := userpb.NewUserServiceClient(h.conn)

	res, err := client.GetUsers(ctx, &emptypb.Empty{})
	if err != nil {
		slog.Error("failed to query users", slog.String("error", err.Error()))
		// status, ok := status.FromError(err)
		// if ok {
		// 	log.Fatalf("status code: %s, error: %s", status.Code().String(), status.Message())
		// }
		// log.Fatal(err)
	}

	err = h.templates.ExecuteTemplate(w, "users.tmpl", res.Users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PageHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	client := userpb.NewUserServiceClient(h.conn)

	pbReq := userpb.CreateUserRequest{
		User: &userpb.User{
			Username:  r.FormValue("username"),
			Email:     r.FormValue("email"),
			FirstName: r.FormValue("first_name"),
			LastName:  r.FormValue("last_name"),
		},
	}

	_, err := client.CreateUser(ctx, &pbReq)
	if err != nil {
		slog.Error("failed to create new user", slog.String("error", err.Error()))
		// status, ok := status.FromError(err)
		// if ok {
		// 	log.Fatalf("status code: %s, error: %s", status.Code().String(), status.Message())
		// }
		// log.Fatal(err)
	}

	w.Header().Set("HX-Refresh", "true")
}
