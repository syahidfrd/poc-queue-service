package queue

type Service interface {
	PublishOrderQueue(productID uint64, quantity uint32)
}
