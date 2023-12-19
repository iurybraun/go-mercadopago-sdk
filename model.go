package mercadopago

type Credentials struct {
    ClientID string
    ClientSecret string
}

type Item struct {
    Id       	string  `json:"id" validate:"id"`
    Title       string  `json:"title" validate:"required"`
    Description string  `json:"description"`
    PictureURL  string  `json:"picture_url"`
    Category_id string  `json:"category_id"`
    Currency_id string  `json:"currency_id"`
    Quantity    int     `json:"quantity" validate:"required"`
    UnitPrice   float64 `json:"unit_price" validate:"required"`
}

type Payer struct {
    First_name    	string `json:"first_name" validate:"first_name"`
    Last_name    	string `json:"last_name" validate:"last_name"`
    Email   		string `json:"email" validate:"required"`
    Phone 			Phone `json:"phone" validate:"required"`
    Identification 	Identification `json:"identification" validate:"identification"`
    Address 		Address `json:"address" validate:"required"`
    CreatedAt 		string `json:"date_created" validate:"required"`
}

type Phone struct {
    Area_code 	string `json:"area_code"`
    Number   	string `json:"number" validate:"required"`
}

type Identification struct {
    Type 	string `json:"type"`
    Number  string `json:"number" validate:"required"`
}

type Address struct {
    Zip_code 		string `json:"zip_code"`
    Street_name  	string `json:"street_name" validate:"street_name"`
    Street_number  	int    `json:"street_number" validate:"street_number"`
    Neighborhood  	string `json:"neighborhood" validate:"neighborhood"`
    City  			string `json:"city" validate:"city"`
}

type NewPreference struct {
	External_reference 	string `json:"external_reference"`
	Description 		string `json:"description"`
    Items 				[]Item `json:"items" validate:"required,min=1"`
    Payment_methods 	Payment_methods `json:"payment_methods"`
    Notification_url 	string `json:"notification_url"`
    Payer 				Payer `json:"payer" validate:"required"`
    Redirect_urls 		Redirect_urls `json:"redirect_urls"`
    Back_urls 			Back_urls `json:"back_urls"`
    AutoReturn 			string `json:"auto_return"`
}

type Redirect_urls struct {
    Success string `json:"success"`
    Pending string `json:"pending"`
    Failure string `json:"failure"`
}

type Back_urls struct {
    Success string `json:"success"`
    Pending string `json:"pending"`
    Failure string `json:"failure"`
}

type Payment_methods struct {
	Excluded_payment_methods 	[]Excluded_payment_methods `json:"excluded_payment_methods"`
    Installments 				int `json:"installments"`
    Default_installments		int `json:"default_installments"`
}

type Excluded_payment_methods struct {
	Id	string `json:"id"`
}
