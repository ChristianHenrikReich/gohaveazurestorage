package gohavestoragecommon

type SignedIdentifiers struct {
	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`
}

type SignedIdentifier struct {
	Id           string
	AccessPolicy AccessPolicy
}

type AccessPolicy struct {
	Start      string
	Expiry     string
	Permission string
}
