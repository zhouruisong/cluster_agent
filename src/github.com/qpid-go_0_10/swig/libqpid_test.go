package libqpid

import "testing"


func TestConnectoin(t *testing.T) {
	f := NewQpidConnection("10.11.144.92:5672", "fuck")
	f.Open()
}


