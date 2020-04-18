package barcode

// DO NOT USE THIS
// ITS NOT FINISHED, WIP
// ITS UNDER DEVELOPEMENT
// AND WILL BE FINISHED WHEN THE TIME HAS COME ;O (spare-time)

import (
	"fmt"
	"math/big"
	"strings"
)

var ascChr = []int{4, 0, 2, 6, 3, 5, 1, 9, 8, 7, 1, 2, 0, 6, 4, 8, 2, 9, 5, 3, 0, 1, 3, 7, 4, 6, 8, 9, 2, 0, 5, 1, 9, 4, 3, 8, 6, 7, 1, 2, 4, 3, 9, 5, 7, 8, 3, 0, 2, 1, 4, 0, 9, 1, 7, 0, 2, 4, 6, 3, 7, 1, 9, 5, 8}

var dscChar = []int{7, 1, 9, 5, 8, 0, 2, 4, 6, 3, 5, 8, 9, 7, 3, 0, 6, 1, 7, 4, 6, 8, 9, 2, 5, 1, 7, 5, 4, 3, 8, 7, 6, 0, 2, 5, 4, 9, 3, 0, 1, 6, 8, 2, 0, 4, 5, 9, 6, 7, 5, 2, 6, 3, 8, 5, 1, 9, 8, 7, 4, 0, 2, 6, 3}

var ascPos = []int{
	3, 0, 8, 11, 1, 12, 8, 11, 10, 6, 4, 12, 2, 7, 9, 6, 7, 9, 2, 8, 4, 0, 12, 7, 10, 9, 0, 7, 10, 5, 7, 9, 6, 8, 2, 12, 1, 4, 2, 0, 1, 5, 4, 6, 12, 1, 0, 9, 4, 7, 5, 10, 2, 6, 9, 11, 2, 12, 6, 7, 5, 11, 0, 3, 2}

var dscPos = []int{2, 10, 12, 5, 9, 1, 5, 4, 3, 9, 11, 5, 10, 1, 6, 3, 4, 1, 10, 0, 2, 11, 8, 6, 1, 12, 3, 8, 6, 4, 4, 11, 0, 6, 1, 9, 11, 5, 3, 7, 3, 10, 7, 11, 8, 2, 10, 3, 5, 8, 0, 3, 12, 11, 8, 4, 5, 1, 3, 0, 7, 12, 9, 8, 10}

/** imd Intelligent Mail Barcode - Onecode - USPS-B-3200
 * Intelligent Mail barcode is a 65-bar code for use on mail in the United States.
 * The fields are described as follows:
 * The Barcode Identifier shall be assigned by USPS to encode the presort identification that
 * is currently printed in human readable form on the optional endorsement line (OEL)
 * as well as for future USPS use. This shall be two digits, with the second digit in
 * the range of 0–4. The allowable encoding ranges shall be 00–04, 10–14, 20–24, 30–34,
 * 40–44, 50–54, 60–64, 70–74, 80–84, and 90–94.
 *
 * The Service Type Identifier shall be assigned by USPS for any combination of services
 * requested on the mailpiece. The allowable encoding range shall be 000
 * Each 3-digit value shall
 * correspond to a particular mail class with a particular combination of service(s).
 * Each service program, such as OneCode Confirm and OneCode ACS, shall provide the list of
 * Service Type Identifier values.
 *
 * The Mailer or Customer Identifier shall be assigned by USPS as a unique, 6 or 9 digit
 * number that identifies a business entity. The allowable encoding range for the 6 digit
 * Mailer ID shall be 000000- 899999, while the allowable encoding range for the 9 digit
 * Mailer ID shall be 900000000-999999999.
 *
 * The Serial or Sequence Number shall be assigned by the mailer for uniquely identifying and
 * tracking mailpieces. The allowable encoding range shall be 000000000–999999999 when used with a
 * 6 digit Mailer ID and 000000-999999 when used with a 9 digit Mailer ID. e.
 * The Delivery Point ZIP Code shall be assigned by the mailer for routing the mailpiece.
 * This shall replace POSTNET for routing the mailpiece to its final delivery point.
 * The length may be 0, 5, 9, or 11 digits.
 * The allowable encoding ranges shall be no ZIP Code, 00000–99999,
 * 000000000–999999999, and 00000000000–99999999999.
 */
func imd(code string) *barArray {
	var binaryCode *big.Int
	routingCode := ""
	//var floatCode float64

	codePart := strings.Split(code, "-")
	trackingNumber := codePart[0]

	if len(codePart) > 2 {
		routingCode = codePart[1]
	}

	switch len(routingCode) {
	case 0:
		binaryCode = big.NewInt(0)
	case 5:
		binaryCode = big.NewInt(100001)
	case 11:
		binaryCode = big.NewInt(1000100001)
	default:
		// invalid routingCode
		return nil
	}

	fmt.Println(new(big.Float).SetPrec(100001).SetString(routingCode))

	fmt.Println(binaryCode)
	_ = routingCode
	_ = binaryCode
	fmt.Println(trackingNumber)
	return nil
}
