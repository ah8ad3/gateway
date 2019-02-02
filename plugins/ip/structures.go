package ip

import "time"

// APIIP this struct fill with all values receive from ip-api and the method are for
type APIIP struct {
	query       string
	status      string
	country     string
	countryCode string
	region      string
	regionName  string
	timezone    string
}

// BlockIPList struct for all ips that banned from service
type BlockIPList struct {
	// the ip address of user
	ip string

	// the service path that user block from it like /foo or /bar
	// only base path counts
	path string

	// the creation of this block time
	createdTime time.Time

	// the expiration time of this ip
	expireTime time.Time

	// is ip blocked for ever or not
	ever bool

	// check if this ip is black or not
	active bool
}
