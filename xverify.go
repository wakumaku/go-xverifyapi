// Package xverifyapi client based on
// http://docs.xverify.com/
package xverifyapi

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// BaseURL of xverify services endpoint
const BaseURL = "http://www.xverify.com/services"

// ServiceType .
type ServiceType int

// Service types supported
const (
	_ ServiceType = iota
	Phone
	Emails
	Address
	Scoring
	AllServices
	PhonePlaceCall
	PhoneConfirmCode
)

// Endpoints endpoint definitions
var Endpoints = map[ServiceType]string{
	Phone:            "/phone/verify",
	Emails:           "/email/verify",
	Address:          "/address/verify",
	Scoring:          "/scoring/verify",
	AllServices:      "/allservices/verify",
	PhonePlaceCall:   "/phoneconfirm/placecall/",
	PhoneConfirmCode: "/phoneconfirm/verifycode",
}

// Errors
var (
	ErrServiceDoesNotExist       = errors.New("Service does not exist")
	ErrEmptyResponse             = errors.New("Empty body response received from service")
	ErrPhoneCallRequiredFields   = errors.New("Code, Country Code and Code should be filled")
	ErrResponseHadEmptyStructure = errors.New("Could not unmarshall correctly the response")
)

// Response with some possible fields. Need to add more fields and split between email, phone etc responses, xverify api docs are broken and a mess...
type Response struct {
	Address     string   `json:"address,omitempty"` // This tag contains the same email address which was entered in the request.
	Syntax      string   `json:"syntax,omitempty"`  // This helps you identify if the email address is in the correct format. If the response is 1 then the email format is good, if the response is 0 then it means the format of the email address is incorrect.
	Handle      string   `json:"handle,omitempty"`  // This field contains the username part of the email ID. This is the characters before the @ symbol.
	Domain      string   `json:"domain,omitempty"`  // This is the domain name of the email ID, the characters after the @ symbol.
	Error       string   `json:"error,omitempty"`   // This helps you identify if there was an error with your request. When this tag displays “0” then there is no error. Anything else in this tag will indicate there is a problem with the request.
	Status      string   `json:"status,omitempty"`  // This tag lets you know if the email address you supplied was either valid of invalid. Valid email addresses are deliverable, and invalid email addresses are not.
	AutoCorrect struct { // We are able to auto correct misspellings of major domain names. If the auto correct feature is enabled in your account then the corrected tag will display true and immediately below that you will see an address tag which will display the corrected email address. IF you have auto-correction enabled, but no corrected occurred then you will see that the corrected tag will display false.
		Corrected string `json:"corrected,omitempty"`
		Address   string `json:"address,omitempty"`
	} `json:"auto_correct,omitempty"`
	Message           string `json:"message,omitempty"`      // The message tag helps provide more details that help explain the response code.
	Duration          string `json:"duration,omitempty"`     // This tag indicates the total execution time for the request.
	CatchAll          string `json:"catch_all,omitempty"`    // This helps you indicate if the email server domain is configured as a catch-all mail server. The results here will display (yes, no, or unknown). Learn more about catch-all domains.
	Responsecode      int    `json:"responsecode,omitempty"` // HTTP status code (when error)
	AreaCode          int
	Prefix            int
	Sufix             int
	TransactionNumber string `json:"transaction_number,omitempty"`
}

// Client holding credentials and connection
type Client struct {
	baseURL    string
	apiKey     string
	domain     string
	httpClient *http.Client
}

// New Client
func New(apiKey, domain string, httpClient *http.Client) *Client {
	return NewWith(BaseURL, apiKey, domain, httpClient)
}

// NewWith a baseURL for test stuff
func NewWith(baseURL, apiKey, domain string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 2 * time.Second,
		}
	}

	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		domain:     domain,
		httpClient: httpClient,
	}
}

// IsEmailVerified wraps the VerifyEmail method
func (c *Client) IsEmailVerified(email string) (bool, error) {
	return c.isValid(c.VerifyEmail(email))
}

// VerifyEmail .
func (c *Client) VerifyEmail(email string) (*Response, error) {
	params := map[string]string{"email": email}
	return c.callService(Emails, params)
}

// IsPhoneVerified wraps the VerifyPhone method
func (c *Client) IsPhoneVerified(phone string) (bool, error) {
	return c.isValid(c.VerifyPhone(phone))
}

// VerifyPhone .
func (c *Client) VerifyPhone(phone string) (*Response, error) {
	params := map[string]string{"phone": phone}
	return c.callService(Phone, params)
}

// VerifyAddress .
func (c *Client) VerifyAddress(street, zip string) (*Response, error) {
	params := map[string]string{
		"street": street,
		"zip":    zip,
	}
	return c.callService(Address, params)
}

// VerifyScoring .
func (c *Client) VerifyScoring(street, zip string) (*Response, error) {
	params := map[string]string{
		"street": street,
		"zip":    zip,
	}
	return c.callService(Scoring, params)
}

// VerifyAllServices .
func (c *Client) VerifyAllServices(email, phone, street, zip string) (*Response, error) {
	params := map[string]string{}
	if email != "" {
		params["services[email]"] = email
	}
	if phone != "" {
		params["services[phone]"] = phone
	}
	if street != "" {
		params["services[address][street]"] = street
	}
	if zip != "" {
		params["services[address][zip]"] = zip
	}

	return c.callService(AllServices, params)
}

// PhoneNumber defines the required fields to place calls
type PhoneNumber struct {
	CountryCode string // Country_code (required): The country code corresponding to the country we will be placing the call to.
	PhoneNumber string // Phone (required): This contains the phone number you wish to place the automated call to.
	Code        string // Code (required): The code to that the user will receive.
}

// CallOptions defines the options when a call is placed
type CallOptions struct {
	RedialCount    int           // Redial_count (optional): You can configure this to redial the call if the user does not pick up.
	RedialInterval time.Duration // Redial_interval (optional): This defines the interval between the redial attempts.
	CallPlaceTile  time.Time     // Call_place_tile (optional): You can configure this to place the call at a specific time
}

// PhonePlaceCall .
func (c *Client) PhonePlaceCall(phoneNumber PhoneNumber, callOptions CallOptions) (*Response, error) {

	if phoneNumber.Code == "" ||
		phoneNumber.PhoneNumber == "" ||
		phoneNumber.CountryCode == "" {
		return nil, ErrPhoneCallRequiredFields
	}

	params := map[string]string{
		"phone":        phoneNumber.PhoneNumber,
		"country_code": phoneNumber.CountryCode,
		"code":         phoneNumber.Code,
	}

	if callOptions.RedialCount != 0 {
		params["redial_count"] = strconv.Itoa(callOptions.RedialCount)
	}

	if callOptions.RedialInterval.Seconds() != 0 {
		params["redial_interval"] = strconv.FormatFloat(callOptions.RedialInterval.Seconds(), 'G', 0, 64)
	}

	if !callOptions.CallPlaceTile.IsZero() {
		params["call_place_tile"] = callOptions.CallPlaceTile.Format("2006-01-02 15:04:05")
	}

	return c.callService(PhonePlaceCall, params)
}

// PhoneConfirmCode .
func (c *Client) PhoneConfirmCode(transactionNumber, code string) (*Response, error) {
	params := map[string]string{
		"transaction_number": transactionNumber,
		"code":               code,
	}

	return c.callService(PhoneConfirmCode, params)
}

func (c *Client) callService(service ServiceType, params map[string]string) (*Response, error) {
	servicePath, _ := Endpoints[service]

	request, err := c.buildRequest("GET", servicePath, params, nil)
	if err != nil {
		return nil, err
	}

	statusCode, body, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	r, err := c.parseBody(body)
	if err != nil {
		r = &Response{
			Responsecode: statusCode,
			Message:      err.Error(),
			Error:        err.Error(),
		}
		return nil, err
	}

	return r, err
}

func (c *Client) parseBody(body []byte) (*Response, error) {
	var wrapper map[string]interface{}
	err := json.Unmarshal(body, &wrapper)
	if err != nil {
		return nil, err
	}

	// Need to pick the first item
	for _, response := range wrapper {
		r, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		var response Response
		err = json.Unmarshal(r, &response)
		if err == nil {
			return &response, nil
		}
		break // Exit from loop, we only want the firs node (probably the only one)
	}

	return nil, ErrResponseHadEmptyStructure
}

func (c *Client) buildRequest(method, path string, params map[string]string, body io.Reader) (*http.Request, error) {
	// Default required values
	params["apikey"] = c.apiKey
	params["domain"] = c.domain
	params["type"] = "json"

	URL, err := c.buildURL(path, params)

	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, URL, body)
}

func (c *Client) buildURL(path string, params map[string]string) (string, error) {

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}
	u.Path += path

	queryString := u.Query()
	for k, v := range params {
		queryString.Set(k, v)
	}

	u.RawQuery = queryString.Encode()

	return u.String(), nil
}

func (c *Client) doRequest(request *http.Request) (int, []byte, error) {
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}

func (c *Client) isValid(response *Response, err error) (bool, error) {
	if err != nil {
		if response != nil {
			return false, errors.New(response.Message)
		}
		return false, err
	}

	return response.Status == "valid", nil
}
