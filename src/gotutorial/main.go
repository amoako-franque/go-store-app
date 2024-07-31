package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/paymentintent"
	"io"
	"log"
	"net/http"
)

// http://localhost:4242/create-payment-intent
func main() {
	stripe.Key = "sk_test_51JJf1yHJJtDKdqBb1SwJho4quMh2Ev431s7cyMP6zF52IzeOZFIcZMJWSriHWGrCYXdSgep0ynIi49WkpSjEjUR500JjTMtmfi"
	http.HandleFunc("/create-payment-intent", handleStripePay)
	http.HandleFunc("/health", healthFunc)

	log.Println("> Listening on port http://localhost:4242")

	var err error = http.ListenAndServe("localhost:4242", nil)

	if err != nil {
		log.Fatal(err)
	}

}

func handleStripePay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
		//w.WriteHeader(http.StatusMethodNotAllowed)
	}

	var req struct {
		ProductId string `json:"product_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address   string `json:"address_1"`
		Address2  string `json:"address_2"`
		City      string `json:"city"`
		State     string `json:"state"`
		Zip       string `json:"zip"`
		Country   string `json:"country"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		//http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(calculateOrderAmount(req.ProductId)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	paymentIntent, err := paymentintent.New(params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println(paymentIntent.ClientSecret)

	var resp struct {
		ClientSecret string `json:"clientSecret"`
	}

	resp.ClientSecret = paymentIntent.ClientSecret

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = io.Copy(w, &buf)

	if err != nil {
		fmt.Println(err)
	}

}

func calculateOrderAmount(productId string) int64 {
	switch productId {
	case "Forever Pants":
		return 26000
	case "Forever Shirt":
		return 15500
	case "Forever Shorts":
		return 30000
	}

	return 0
}

func healthFunc(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("server is listening on http://localhost:4242")
	var response []byte = []byte("server up and running")

	_, err := w.Write(response)
	if err != nil {
		fmt.Println(err)
	}

}

func returnMultiple() (int, string, bool) {
	return 34, "string", true
}
