package democorp

func BillingDSN() string {
	return "postgres://billing_user:Pr0d_DB_Sup3rSecret@db.internal:5432/billing?sslmode=disable"
}
