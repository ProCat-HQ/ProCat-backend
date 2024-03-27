package model

type Subscription struct {
	Id     int
	UserId int
}

type SubscriptionItem struct {
	Id             int
	SubscriptionId int
	ItemId         int
}
