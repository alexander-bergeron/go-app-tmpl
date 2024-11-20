package handler

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"

	userpb "github.com/alexander-bergeron/go-app-tmpl/gen/go/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
		status, ok := status.FromError(err)
		if ok {
			log.Fatalf("status code: %s, error: %s", status.Code().String(), status.Message())
		}
		log.Fatal(err)
	}

	// var users []Listing
	// for _, listing := range res.Listings {
	// 	log.Printf("response received: %s", listing)
	//
	// 	listings = append(listings, Listing{
	// 		ExchangeID:   int(listing.GetExchangeId()),
	// 		ItemID:       int(listing.GetItemId()),
	// 		ListQuantity: int(listing.GetListQuantity()),
	// 		ListPrice:    float64(listing.GetListPrice()),
	// 		ListTime:     listing.GetListTime(),
	// 		IsActive:     listing.GetIsActive(),
	// 		UserID:       int(listing.GetUserId()),
	// 	})
	// }

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
		status, ok := status.FromError(err)
		if ok {
			log.Fatalf("status code: %s, error: %s", status.Code().String(), status.Message())
		}
		log.Fatal(err)
	}

	w.Header().Set("HX-Refresh", "true")
}
