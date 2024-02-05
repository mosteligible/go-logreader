package notifier

import "github.com/mosteligible/go-logreader/client/customer"

func NotifyService(serviceEndpoint string, cust customer.Customer) {
	// send message to concerned services about updates in
	// clientel they should expect and update states accordingly
}
