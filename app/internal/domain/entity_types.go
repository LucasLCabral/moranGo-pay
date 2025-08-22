package domain

// EntityType representa os tipos de entidades na tabela
type EntityType string

const (
	EntityTypeProfile     EntityType = "PROFILE"
	EntityTypeWallet      EntityType = "WALLET"
	EntityTypeTransaction EntityType = "TRANSACTION"
	EntityTypeAccount     EntityType = "ACCOUNT"     // TODO
	EntityTypePayment     EntityType = "PAYMENT"     // TODO
	EntityTypeContact     EntityType = "CONTACT"     // TODO
	EntityTypeNotification EntityType = "NOTIFICATION" // TODO
)

func (e EntityType) String() string {
	return string(e)
}

func (e EntityType) IsValid() bool {
	switch e {
	case EntityTypeProfile, EntityTypeWallet, EntityTypeTransaction,
		 EntityTypeAccount, EntityTypePayment, EntityTypeContact, EntityTypeNotification:
		return true
	default:
		return false
	}
}