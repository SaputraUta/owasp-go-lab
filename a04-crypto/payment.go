package main

import "os"

func paymentClientKey() string {
	return os.Getenv("PAYMENT_API_KEY")
}
