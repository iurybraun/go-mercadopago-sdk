package mercadopago

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
)

const _baseURL = "https://api.mercadopago.com"

type Client interface {
    Do(req *http.Request) (*http.Response, error)
}

type Gateway struct {
    Client Client
}

func NewClientGateway(client Client) *Gateway {
    return &Gateway{
        Client: client,
    }
}

type PaymentReq struct {
	Id int `json:"id"`
	Client_id string `json:"client_id"`
	Collector_id int `json:"collector_id"`
	Currency_id string `json:"currency_id"`
	Date_approved string `json:"date_approved"`
	External_reference string `json:"external_reference"`
	Installments int `json:"installments"`
	Order struct {
		Id string `json:"id"`
		Type string `json:"type"`
	} `json:"order"`
	Transaction_amount int `json:"transaction_amount"`
	Captured bool `json:"captured"`
	Status string `json:"status"`
}

type PaymentReqSearch struct {
	Results	[]struct {
		Id int `json:"id"`
		External_reference string `json:"external_reference"`
		Collector_id int `json:"collector_id"`
		Currency_id string `json:"currency_id"`
		Status string `json:"status"`
	} `json:"results"`
}

func (g *Gateway) GetAccessToken(credentials Credentials) (string, error) {
    path := &url.Values{}
    path.Add("client_id", credentials.ClientID)
    path.Add("client_secret", credentials.ClientSecret)
    path.Add("grant_type", "client_credentials")
    queryParams := path.Encode()

    req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", _baseURL, "/oauth/token?", queryParams), nil)
    if err != nil {
        return "", err
    }

    resp, err := g.Client.Do(req)
    if err != nil {
        return "", err
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    if resp.StatusCode >= http.StatusBadRequest {
        return "", NewError(string(body), resp.StatusCode)
    }

    var r struct {
        AccessToken string `json:"access_token"`
    }

    if err := json.Unmarshal(body, &r); err != nil {
        return "", err
    }

    return r.AccessToken, nil
}

func (g *Gateway) CreatePreference(accessToken string, preference NewPreference) (string, string, error) {
    queryValues := &url.Values{}
    queryValues.Add("access_token", accessToken)
    queryParams := queryValues.Encode()

    b, err := json.Marshal(preference)
    if err != nil {
        return "", "", err
    }

    req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", _baseURL, "/checkout/preferences?", queryParams), bytes.NewReader(b))
    if err != nil {
        return "", "", err
    }

    resp, err := g.Client.Do(req)
    if err != nil {
        return "", "", err
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", "", err
    }

    if resp.StatusCode >= http.StatusBadRequest {
        return "", "", NewError(string(body), resp.StatusCode)
    }

    var r struct {
        Id string `json:"id"`
        Client_id string `json:"client_id"`
        CheckoutURL string `json:"init_point"`
    }

    if err := json.Unmarshal(body, &r); err != nil {
        return "", "", err
    }

    return r.Id, r.CheckoutURL, nil
}

func (g *Gateway) GetCheckoutPreferences(accessToken string, id string) (int, error) {
    ///queryValues := &url.Values{}
    ///queryValues.Add("limit", "1")
    ///queryValues.Add("offset", "0")
    ///queryValues.Add("access_token", accessToken)
    ///queryValues.Add("id", id)

    ///queryParams := queryValues.Encode()
    
    req, err := http.NewRequest("GET", _baseURL + "/checkout/preferences/"+id, nil)
    if err != nil {
        return 0, err
    }
    
    req.Header.Add("Authorization", "Bearer " + accessToken)

    resp, err := g.Client.Do(req)
    if err != nil {
        return 0, err
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, err
    }

    if resp.StatusCode >= http.StatusBadRequest {
        return 0, NewError(string(body), resp.StatusCode)
    }
    
    var r struct {
		Id string `json:"id"`
		Client_id string `json:"client_id"`
		Collector_id int `json:"collector_id"`
		External_reference string `json:"external_reference"`
		Total_amount int `json:"total_amount"`
	}

    if err := json.Unmarshal(body, &r); err != nil {
        return 0, err
    }
    
    return r.Total_amount, nil
}

func (g *Gateway) GetPayments(accessToken string, id string) (payment PaymentReq, err error) {
    ///queryValues := &url.Values{}
    ///queryValues.Add("limit", "1")
    ///queryValues.Add("offset", "0")
    ///queryValues.Add("access_token", accessToken)
    ///queryValues.Add("id", id)

    ///queryParams := queryValues.Encode()
    
    req, err := http.NewRequest("GET", _baseURL + "/v1/payments/"+id, nil)
    if err != nil {
        return
    }
    
    req.Header.Add("Authorization", "Bearer " + accessToken)

    resp, err := g.Client.Do(req)
    if err != nil {
        return
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
	}

    if resp.StatusCode >= http.StatusBadRequest {
        err = NewError(string(body), resp.StatusCode)
    }
    
    r := PaymentReq{}
    if err := json.Unmarshal(body, &r); err != nil {
        return payment, err
    }
    
    return
}

func (g *Gateway) GetPaymentsSearch(accessToken string, external_reference string) (payment PaymentReqSearch, err error) {
    ///queryValues := &url.Values{}
    ///queryValues.Add("limit", "1")
    ///queryValues.Add("offset", "0")
    ///queryValues.Add("access_token", accessToken)
    ///queryValues.Add("id", id)

    ///queryParams := queryValues.Encode()
    
    //https://api.mercadopago.com/v1/payments/search?sort=date_created&criteria=desc&external_reference=ID_a9XY6Qd+aKTswbX2sdZQ/B0Mzs8pSWnzynl/CR1Ek5Y
    
    req, err := http.NewRequest("GET", _baseURL + "/v1/payments/search?sort=date_created&criteria=desc&external_reference="+external_reference, nil)
    if err != nil {
        return
    }
    
    req.Header.Add("Authorization", "Bearer " + accessToken)

    resp, err := g.Client.Do(req)
    if err != nil {
        return
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
	}

    if resp.StatusCode >= http.StatusBadRequest {
        err = NewError(string(body), resp.StatusCode)
    }
    
    ///r := PaymentReqSearch{}
    if err = json.Unmarshal(body, &payment); err != nil {
		return
    }
    
    return
}

func (g *Gateway) GetTotalPayments(accessToken string, status string) (int, error) {
    queryValues := &url.Values{}
    queryValues.Add("limit", "1")
    queryValues.Add("offset", "0")
    queryValues.Add("access_token", accessToken)
    queryValues.Add("status", status)

    queryParams := queryValues.Encode()

    req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", _baseURL, "/v1/payments/search?", queryParams), nil)
    if err != nil {
        return 0, err
    }

    resp, err := g.Client.Do(req)
    if err != nil {
        return 0, err
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, err
    }

    if resp.StatusCode >= http.StatusBadRequest {
        return 0, NewError(string(body), resp.StatusCode)
    }

    var r struct {
        Paging struct {
            TotalPayments int `json:"total"`
        } `json:"paging"`
    }

    if err := json.Unmarshal(body, &r); err != nil {
        return 0, err
    }

    return r.Paging.TotalPayments, nil
}
