package modem_test

import (
	"testing"

	"github.com/stv0g/vand/pkg/devices/modem"
)

func TestModem(t *testing.T) {
	m, err := modem.New("192.168.1.1", "admin", "KbazEz7e")
	if err != nil {
		t.Fatal(err)
	}

	if err := m.Login(); err != nil {
		t.Fatal(err)
	}

	mdl, err := m.GetModel()
	if err != nil {
		t.Fatal(err)
	}

	if mdl.Session.UserRole != "Admin" {
		t.Fail()
	}

	// if err := m.SendSMS("015757180927", "Hallo welt 1234"); err != nil {
	// 	t.Fatal(err)
	// }

	if err := m.ConnectLTE(); err != nil {
		t.Fatal(err)
	}
}
