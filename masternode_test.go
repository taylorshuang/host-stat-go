package hoststat

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegStr(t *testing.T) {
	str := "Total: 2636 ( Enabled: 2252 / Qualify: 1903)"

	reg := regexp.MustCompile(`[0-9]+`)
	result := reg.FindAllString(str, -1)
	println(result[0], result[1], result[2])
}

func TestMoreRegStr(t *testing.T) {
	str := `"47e1e3ebc1f919bac652b9a8988b49b75a135ab013a62a43e65714d2699ee2ff-0": "           ENABLED 170005 VcTrwep8XMQGe7kwmFcHAgEEtHCi1fYkVLC 1556930102    38100 1556909592 100767 9de9c46cb1fb09b951c3c09dce5a51482ece8d26",`
	reg := regexp.MustCompile(`[0-9A-Za-z]+`)
	result := reg.FindAllString(str, -1)
	fmt.Println(len(result), result)
}
