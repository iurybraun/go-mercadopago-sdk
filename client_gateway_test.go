package mercadopago

import (
    "bytes"
    "errors"
    "github.com/stretchr/testify/require"
    "io/ioutil"
    "net/http"
    "testing"
)

type ClientStub struct {
    resp *http.Response
    err  error
}

func (c *ClientStub) Do(_ *http.Request) (*http.Response, error) {
    if c.err != nil {
        return &http.Response{}, c.err
    }

    return c.resp, nil
}

func TestGateway_GetAccessToken(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.resp = &http.Response{
        Status:     "200",
        StatusCode: 200,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"access_token": "1234"}`))),
    }
    // When
    accessToken, err := g.GetAccessToken(Credentials{
        ClientID:     "ABC123",
        ClientSecret: "123ABC",
    })

    // Then
    require.NoError(t, err)
    require.Equal(t, accessToken, "1234")
}

func TestGateway_GetAccessToken_MercadoPagoError(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.resp = &http.Response{
        Status:     "500",
        StatusCode: 500,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "internal server error"}`))),
    }
    // When
    _, err := g.GetAccessToken(Credentials{
        ClientID:     "ABC123",
        ClientSecret: "123ABC",
    })

    // Then
    require.Error(t, err)
    require.EqualError(t, err, "{\"error\": \"internal server error\"}")
}

func TestGateway_GetAccessToken_UnmarshalError(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.resp = &http.Response{
        Status:     "200",
        StatusCode: 200,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"access_token": 1123}`))),
    }
    // When
    _, err := g.GetAccessToken(Credentials{
        ClientID:     "ABC123",
        ClientSecret: "123ABC",
    })

    // Then
    require.Error(t, err)
    require.EqualError(t, err, "json: cannot unmarshal number into Go struct field .access_token of type string")
}

func TestGateway_GetAccessToken_DoError(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.err = errors.New("do error")
    // When
    _, err := g.GetAccessToken(Credentials{
        ClientID:     "ABC123",
        ClientSecret: "123ABC",
    })

    // Then
    require.Error(t, err)
    require.EqualError(t, err, "do error")
}

func TestGateway_CreatePreference(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.resp = &http.Response{
        Status:     "200",
        StatusCode: 200,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"init_point": "https://mercadopago.com/checkout"}`))),
    }
    // When
    _, checkout, err := g.CreatePreference("", newPreference())

    // Then
    require.NoError(t, err)
    require.Equal(t, checkout, "https://mercadopago.com/checkout")
}

func TestGateway_CreatePreference_MercadoPagoError(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.resp = &http.Response{
        Status:     "500",
        StatusCode: 500,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "internal server error"}`))),
    }
    // When
    _, _, err := g.CreatePreference("", newPreference())

    // Then
    require.Error(t, err)
    require.EqualError(t, err, "{\"error\": \"internal server error\"}")
}

func TestGateway_CreatePreference_UnmarshalError(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.resp = &http.Response{
        Status:     "200",
        StatusCode: 200,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"init_point": 1234}`))),
    }
    // When
    _, _, err := g.CreatePreference("", newPreference())

    // Then
    require.Error(t, err)
    require.EqualError(t, err, "json: cannot unmarshal number into Go struct field .init_point of type string")
}

func TestGateway_CreatePreference_DoError(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.err = errors.New("do error")
    // When
    _, _, err := g.CreatePreference("", newPreference())

    // Then
    require.Error(t, err)
    require.EqualError(t, err, "do error")
}

func TestGateway_GetTotalPayments(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.resp = &http.Response{
        Status:     "200",
        StatusCode: 200,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"paging": {"total": 100,"limit": 1,"offset": 0}}`))),
    }
    // When
    totalPayments, err := g.GetTotalPayments("MY_ACCESS_TOKEN", "approved")
    if err != nil {
        t.Fatal(err)
    }

    // Then
    require.NoError(t, err)
    require.Equal(t, 100, totalPayments)
}

func TestGateway_GetTotalPayments_MercadoPagoError(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.resp = &http.Response{
        Status:     "500",
        StatusCode: 500,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "internal server error"}`))),
    }
    // When
    totalPayments, err := g.GetTotalPayments("MY_ACCESS_TOKEN", "approved")

    // Then
    require.Error(t, err)
    require.EqualError(t, err, "{\"error\": \"internal server error\"}")
    require.Equal(t, 0, totalPayments)
}

func TestGateway_GetTotalPayments_UnmarshalError(t *testing.T) {
    // Given
    c := &ClientStub{}
    g := &Gateway{Client: c}
    c.resp = &http.Response{
        Status:     "200",
        StatusCode: 200,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"paging": 0}`))),
    }
    // When
    totalPayments, err := g.GetTotalPayments("MY_ACCESS_TOKEN", "approved")

    // Then
    require.Error(t, err)
    require.EqualError(t, err, "json: cannot unmarshal number into Go struct field .paging of type struct { TotalPayments int \"json:\\\"total\\\"\" }")
    require.Equal(t, 0, totalPayments)
}

func newPreference() NewPreference {
    return NewPreference{
        Items: []Item{
            {
                Title:       "sherlock",
                Description: "holes",
                PictureURL:  "",
                Quantity:    1,
                UnitPrice:   15.75,
            },
        },
        Payer: Payer{
            Name:    "mateo",
            Surname: "fc",
            Email:   "m@gmail.com",
            Phone: Phone{
                AreaCode: "",
                Number:   "12345",
            },
            Address: Address{
                ZipCode: "",
                Street:  "pepe",
                Number:  1234,
            },
            CreatedAt: "",
        },
        Redirect: Redirect{
            Success: "http://baseurl.com/success",
            Pending: "http://baseurl.com/pending",
            failure: "http://baseurl.com/failure",
        },
        AutoReturn: true,
    }
}
