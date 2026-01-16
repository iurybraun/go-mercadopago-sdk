package mercadopago

type ClientGateway interface {
    GetAccessToken(credentials Credentials) (string, error)
    CreatePreference(accessToken string, preference NewPreference) (string, string, error)
    GetCheckoutPreferences(accessToken string, id string) (int, error)
    GetPayments(accessToken string, id string) (PaymentReq, error)
    GetPaymentsSearch(accessToken string, external_reference string) (PaymentReqSearch, error)
	//GetSubscriptionsSearch(accessToken string, filters map[string]interface{}) (*preapproval.SearchResponse, error)
    //GetMerchantOrders(accessToken string, order_id string) (MerchantOrders, error)
    GetTotalPayments(accessToken string, status string) (int, error)
}

type Controller struct {
    Client ClientGateway
}

func NewController(client ClientGateway) *Controller {
    return &Controller{
        Client: client,
    }
}

func (s *Controller) GetAccessToken(clientID string, clientSecret string) (string, error) {
    return s.Client.GetAccessToken(Credentials{
        ClientID:     clientID,
        ClientSecret: clientSecret,
    })
}

func (s *Controller) CreatePreference(accessToken string, preference NewPreference) (string, string, error) {
    return s.Client.CreatePreference(accessToken, preference)
}

func (s *Controller) GetCheckoutPreferences(accessToken string, id string) (int, error) {
    return s.Client.GetCheckoutPreferences(accessToken, id)
}

func (s *Controller) GetPayments(accessToken string, id string) (PaymentReq, error) {
    return s.Client.GetPayments(accessToken, id)
}

func (s *Controller) GetPaymentsSearch(accessToken string, external_reference string) (PaymentReqSearch, error) {
    return s.Client.GetPaymentsSearch(accessToken, external_reference)
}

/*func (s *Controller) GetSubscriptionsSearch(accessToken string, filters map[string]interface{}) (*preapproval.SearchResponse, error) {
	return s.Client.GetSubscriptionsSearch(accessToken, filters)
}*/

/*func (s *Controller) GetMerchantOrders(accessToken string, order_id string) (MerchantOrders, error) {
    return s.Client.GetMerchantOrders(accessToken, order_id)
}*/

func (s *Controller) GetTotalPayments(accessToken string, status string) (int, error) {
    return s.Client.GetTotalPayments(accessToken, status)
}
