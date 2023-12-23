package utils

import (
	"math/rand"
	"strings"
)

// GetRandomHost gets all hosts instances in the format instance1,instance2,instance3 and takes a random one
func GetRandomHost(hosts string) string {
	host := strings.Split(hosts, ",")
	return host[rand.Intn(len(host))]
}
