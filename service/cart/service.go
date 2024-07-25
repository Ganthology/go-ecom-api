package cart

import (
	"fmt"

	"github.com/ganthology/go-ecom-api/types"
)

func GetCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}

		productIds[i] = item.ProductID
	}

	return productIds, nil
}

func (h *Handler) CreateOrder(products []types.Product, items []types.CartItem, userId int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// Check stock
	if err := checkIfItemInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	// Calculate price
	totalPrice := calculateTotalPrice(items, productMap)

	// Update stock
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}

	// Create order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userId,
		Total:   totalPrice,
		Status:  "pending",
		Address: "123 Main St",
	})
	if err != nil {
		return 0, 0, err
	}

	// Create order items
	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkIfItemInStock(items []types.CartItem, productMap map[int]types.Product) error {
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d not found in store", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %d out of stock", item.ProductID)
		}
	}

	return nil
}

func calculateTotalPrice(items []types.CartItem, productMap map[int]types.Product) float64 {
	var totalPrice float64
	for _, item := range items {
		product := productMap[item.ProductID]
		totalPrice += product.Price * float64(item.Quantity)
	}

	return totalPrice
}
