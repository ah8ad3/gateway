package ip

import "time"

// this struct fill with all values receive from ip-api and the method are for
type ApiIp struct {
	query string
	status string
	country string
	countryCode string
	region string
	regionName string
	timezone string
}

type BlockIpList struct {
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
