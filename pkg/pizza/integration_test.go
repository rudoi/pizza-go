// +build integration

package pizza

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var client *Client

var _ = BeforeSuite(func() {
	client = &Client{
		Client: http.Client{Timeout: 10 * time.Second},
	}
})

var _ = Describe("pricing flow", func() {
	Context("given an address and a product", func() {
		It("should return a price", func() {

			// big pink building in
			// downtown Portland, OR
			address := &Address{
				Street:     "111 SW 5th Ave",
				PostalCode: "97204",
			}

			product := &OrderProduct{Code: "14SCREEN", Qty: 1}

			By("looking up the nearest store")
			store, err := client.GetNearestStore(address)
			Expect(err).ToNot(HaveOccurred())
			Expect(store).ToNot(BeNil())
			Expect(store.StoreID).To(Equal("7229"))

			By("looking up the store's menu")
			menu, err := client.GetStoreMenu(store.StoreID)
			Expect(err).ToNot(HaveOccurred())
			Expect(menu).ToNot(BeNil())

			By("building an order")
			order := NewOrder().
				WithAddress(address).
				WithStoreID(store.StoreID).
				WithPhoneNumber("1235467890")

			order.AddProduct(product)
			order.AddCoupon(menu.GetFiftyPercentCouponCode())

			By("pricing the order")
			price, err := client.ValidateOrder(order)

			if store.IsOpen {
				Expect(err).ToNot(HaveOccurred())
				Expect(price).ToNot(BeZero())
			} else {
				Expect(err).To(HaveOccurred())
				Expect(price).To(BeZero())
			}
		})
	})
})
