package mercadopago

import (
    "encoding/json"
    "fmt"
    "github.com/go-playground/validator/v10"
    "net/http"
)

var _v = validator.New()

type Service interface {
    GetAccessToken(clientID string, clientSecret string) (string, error)
    CreatePreference(accessToken string, preference NewPreference) (string, string, error)
    GetTotalPayments(accessToken string, status string) (int, error)
}

type Handler struct {
    Service Service
}

func NewHandler(service Service) *Handler{
    return &Handler{
        Service: service,
    }
}

func (h *Handler) Ping(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "pong")
}

func (h *Handler) GetAccessToken(w http.ResponseWriter, r *http.Request) {
    clientID := r.URL.Query().Get("client_id")
    if clientID == "" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "client id is required")
        return
    }

    clientSecret := r.URL.Query().Get("client_secret")
    if clientSecret == "" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "client secret is required")
        return
    }

    accessToken, err := h.Service.GetAccessToken(clientID, clientSecret)
    if err != nil {
        w.WriteHeader(getStatusCodeFromError(err))
        fmt.Fprintf(w, "couldn't get access token: %v", err)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "%s", accessToken)
}

func (h *Handler) CreatePreference(w http.ResponseWriter, r *http.Request) {
    var preference NewPreference
    if err := json.NewDecoder(r.Body).Decode(&preference); err != nil {
        w.WriteHeader(http.StatusUnprocessableEntity)
        fmt.Fprintf(w, "couldn't decode body: %v", err)
        return
    }

    if err := _v.Struct(preference); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, fmt.Sprintf("validation error: %v", err))
        return
    }

    for _, i := range preference.Items {
        if err := _v.Struct(i); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, fmt.Sprintf("validation error: %v", err))
            return
        }
    }


    accessToken := r.Header.Get("access_token")
    if accessToken == "" {
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, fmt.Sprintf("access token is required"))
        return
    }

    id, checkoutURL, err := h.Service.CreatePreference(accessToken, preference)
    if err != nil {
        w.WriteHeader(getStatusCodeFromError(err))
        fmt.Fprintf(w, "couldn't create checkout: %v", err)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, fmt.Sprintf("%s", id))
    fmt.Fprintf(w, fmt.Sprintf("%s", checkoutURL))
}

func (h *Handler) GetTotalPayments(w http.ResponseWriter, r *http.Request) {
    accessToken := r.Header.Get("access_token")
    if accessToken == "" {
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, fmt.Sprintf("access token is required"))
        return
    }

    status := r.URL.Query().Get("status")
    if status == "" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, fmt.Sprintf("status is required"))
        return
    }

    if status != "approved" && status != "rejected" && status != "pending" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, fmt.Sprintf("invalid status: got: %s, want: approved, rejected or pending", status))
        return
    }

    total, err := h.Service.GetTotalPayments(accessToken, status)
    if err != nil {
        w.WriteHeader(getStatusCodeFromError(err))
        fmt.Fprintf(w, fmt.Sprintf("couldn't get total payments: %v", err))
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, fmt.Sprintf("total payments: %d", total))
}

func getStatusCodeFromError(err error) int {
    e, ok := err.(*Error)
    if !ok {
        return http.StatusInternalServerError
    }
    
    return e.StatusCode
}
