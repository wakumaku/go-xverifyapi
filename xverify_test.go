package xverifyapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	xverifyapi "github.com/wakumaku/go-xverifyapi"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *xverifyapi.Client
)

const (
	apiKey                 = "123456"
	domain                 = "domain.tld"
	defaultType            = "json"
	validEmail             = "valid@email.com"
	validPhoneNumber       = "1234567890"
	validStreet            = "Sesame street"
	validZip               = "UX002"
	validCountryCode       = "01"
	validSecretCode        = "AB01CD"
	validRedialCount       = 3
	validRedialInterval    = 5 * time.Minute
	validCallPlaceTile     = "2020-01-01 21:12:13"
	validTransactionNumber = "123-456-abcd"
)

var validResponse = `{
    "wrappernode": {
        "address":"address",
        "syntax":"syntax",
        "handle":"handle",
        "domain":"domain",
        "error":0,
        "status":"valid",
        "auto_correct": {
            "corrected":"corrected",
            "address":"address"
        },
        "message":"message",
        "duration":0.0,
        "catch_all":"catch_all",
        "responsecode":200,
        "transaction_number":"123-456-abcd"
    }
}`

var inValidResponse = `{
    "wrappernode": {
        "status": "invalid",
        "message": "error occurred",
        "error": "error occurred",
        "responsecode": 503
    }
}`

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// Mock xverify server: email
	mux.HandleFunc(xverifyapi.Endpoints[xverifyapi.Emails], func(w http.ResponseWriter, r *http.Request) {
		response := validResponse
		status := http.StatusOK
		if !validateMandatoryParams(r) {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		if r.URL.Query().Get("email") != validEmail {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(response))
	})

	// Mock xverify server: phone
	mux.HandleFunc(xverifyapi.Endpoints[xverifyapi.Phone], func(w http.ResponseWriter, r *http.Request) {
		response := validResponse
		status := http.StatusOK
		if !validateMandatoryParams(r) {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		if r.URL.Query().Get("phone") != validPhoneNumber {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(response))
	})

	// Mock xverify server: address
	mux.HandleFunc(xverifyapi.Endpoints[xverifyapi.Address], func(w http.ResponseWriter, r *http.Request) {
		response := validResponse
		status := http.StatusOK
		if !validateMandatoryParams(r) {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		if r.URL.Query().Get("street") != validStreet &&
			r.URL.Query().Get("zip") != validZip {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(response))
	})

	// Mock xverify server: scoring
	mux.HandleFunc(xverifyapi.Endpoints[xverifyapi.Scoring], func(w http.ResponseWriter, r *http.Request) {
		response := validResponse
		status := http.StatusOK
		if !validateMandatoryParams(r) {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		if r.URL.Query().Get("street") != validStreet &&
			r.URL.Query().Get("zip") != validZip {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(response))
	})

	// Mock xverify server: AllServices
	mux.HandleFunc(xverifyapi.Endpoints[xverifyapi.AllServices], func(w http.ResponseWriter, r *http.Request) {
		response := validResponse
		status := http.StatusOK
		if !validateMandatoryParams(r) {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		if r.URL.Query().Get("services[email]") != validEmail &&
			r.URL.Query().Get("services[phone]") != validPhoneNumber &&
			r.URL.Query().Get("services[address][street]") != validStreet &&
			r.URL.Query().Get("services[address][zip]") != validZip {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(response))
	})

	// Mock xverify server: PlaceCall
	mux.HandleFunc(xverifyapi.Endpoints[xverifyapi.PhonePlaceCall], func(w http.ResponseWriter, r *http.Request) {
		response := validResponse
		status := http.StatusOK
		if !validateMandatoryParams(r) {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		if r.URL.Query().Get("phone") != validPhoneNumber &&
			r.URL.Query().Get("country_code") != validCountryCode &&
			r.URL.Query().Get("code") != validSecretCode {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(response))
	})

	// Mock xverify server: PhoneConfirmCode
	mux.HandleFunc(xverifyapi.Endpoints[xverifyapi.PhoneConfirmCode], func(w http.ResponseWriter, r *http.Request) {
		response := validResponse
		status := http.StatusOK
		if !validateMandatoryParams(r) {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		if r.URL.Query().Get("transaction_number") != validTransactionNumber &&
			r.URL.Query().Get("code") != validSecretCode {
			status = http.StatusInternalServerError
			response = inValidResponse
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(response))
	})

	client = xverifyapi.NewWith(server.URL, apiKey, domain, nil)

	return func() {
		server.Close()
	}
}

func validateMandatoryParams(r *http.Request) bool {
	return r.URL.Query().Get("apikey") == apiKey &&
		r.URL.Query().Get("domain") == domain &&
		r.URL.Query().Get("type") == defaultType
}

func TestVerifyEmail(t *testing.T) {
	teardown := setup()
	defer teardown()

	r, err := client.VerifyEmail(validEmail)
	if err != nil {
		t.Fatal(err)
	}

	if r.Status != "valid" ||
		r.Message != "message" {
		t.Error("Unexpected field values")
	}
}

func TestIsEmailVerified(t *testing.T) {
	teardown := setup()
	defer teardown()

	verified, err := client.IsEmailVerified(validEmail)
	if err != nil {
		t.Fatal(err)
	}

	if verified == false {
		t.Error("Email should be verified")
	}
}

func TestIsEmailVerifiedShouldFail(t *testing.T) {
	teardown := setup()
	defer teardown()

	verified, err := client.IsEmailVerified("invalidEmail")
	if err != nil {
		t.Fatal(err)
	}

	if verified == true {
		t.Error("Email should not be verified")
	}
}

func TestVerifyPhone(t *testing.T) {
	teardown := setup()
	defer teardown()

	r, err := client.VerifyPhone(validPhoneNumber)
	if err != nil {
		t.Fatal(err)
	}

	if r.Status != "valid" ||
		r.Message != "message" {
		t.Error("Unexpected field values")
	}
}

func TestVerifyAddress(t *testing.T) {
	teardown := setup()
	defer teardown()

	r, err := client.VerifyAddress(validStreet, validZip)
	if err != nil {
		t.Fatal(err)
	}

	if r.Status != "valid" ||
		r.Message != "message" {
		t.Error("Unexpected field values")
	}
}
func TestVerifyScoring(t *testing.T) {
	teardown := setup()
	defer teardown()

	r, err := client.VerifyScoring(validStreet, validZip)
	if err != nil {
		t.Fatal(err)
	}

	if r.Status != "valid" ||
		r.Message != "message" {
		t.Error("Unexpected field values")
	}
}

func TestVerifyAllServices(t *testing.T) {

	teardown := setup()
	defer teardown()

	r, err := client.VerifyAllServices(validEmail, validPhoneNumber, validStreet, validZip)
	if err != nil {
		t.Fatal(err)
	}

	if r.Status != "valid" ||
		r.Message != "message" {
		t.Error("Unexpected field values", r.Message)
	}
}

func TestVerifyPlaceCall(t *testing.T) {

	teardown := setup()
	defer teardown()

	phoneNumber := xverifyapi.PhoneNumber{
		CountryCode: validCountryCode,
		PhoneNumber: validPhoneNumber,
		Code:        validSecretCode,
	}

	callOptions := xverifyapi.CallOptions{
		RedialCount:    validRedialCount,
		RedialInterval: validRedialInterval,
		CallPlaceTile:  time.Time{},
	}
	r, err := client.PhonePlaceCall(phoneNumber, callOptions)
	if err != nil {
		t.Fatal(err)
	}

	if r.Status != "valid" ||
		r.Message != "message" {
		t.Error("Unexpected field values", r.Message)
	}
}

func TestVerifyPhoneConfirmCode(t *testing.T) {

	teardown := setup()
	defer teardown()

	r, err := client.PhoneConfirmCode(validTransactionNumber, validSecretCode)
	if err != nil {
		t.Fatal(err)
	}

	if r.Status != "valid" ||
		r.Message != "message" {
		t.Error("Unexpected field values", r.Message)
	}
}
