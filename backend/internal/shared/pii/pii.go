// Package pii provides masking helpers and Turkish national ID (TCKN)
// validation used to satisfy KVKK data-minimization requirements when
// returning personal data to clients that are not authorized to see it in full.
package pii

import "strings"

// MaskTCKN masks a Turkish national ID, keeping the first 3 and last 2 digits.
// Example: 12345678901 -> 123******01
func MaskTCKN(tckn string) string {
	if len(tckn) != 11 {
		return strings.Repeat("*", len(tckn))
	}
	return tckn[:3] + "******" + tckn[9:]
}

// MaskIBAN keeps the 4-char country/check prefix and the last 4 characters.
func MaskIBAN(iban string) string {
	iban = strings.ReplaceAll(iban, " ", "")
	n := len(iban)
	if n <= 8 {
		return strings.Repeat("*", n)
	}
	return iban[:4] + strings.Repeat("*", n-8) + iban[n-4:]
}

// MaskPhone keeps only the last 2 digits visible.
func MaskPhone(phone string) string {
	n := len(phone)
	if n <= 2 {
		return strings.Repeat("*", n)
	}
	return strings.Repeat("*", n-2) + phone[n-2:]
}

// MaskEmail keeps the first character of the local part and the full domain.
// Example: john@acme.com -> j***@acme.com
func MaskEmail(email string) string {
	at := strings.Index(email, "@")
	if at <= 1 {
		return email
	}
	return email[:1] + strings.Repeat("*", at-1) + email[at:]
}

// IsValidTCKN validates a Turkish national identification number using the
// official checksum algorithm (11 digits, first digit non-zero, two check digits).
func IsValidTCKN(tckn string) bool {
	if len(tckn) != 11 {
		return false
	}
	var d [11]int
	for i := 0; i < 11; i++ {
		c := tckn[i]
		if c < '0' || c > '9' {
			return false
		}
		d[i] = int(c - '0')
	}
	if d[0] == 0 {
		return false
	}

	oddSum := d[0] + d[2] + d[4] + d[6] + d[8]  // 1st,3rd,5th,7th,9th
	evenSum := d[1] + d[3] + d[5] + d[7]         // 2nd,4th,6th,8th
	check10 := ((oddSum * 7) - evenSum) % 10
	if check10 < 0 {
		check10 += 10
	}
	if check10 != d[9] {
		return false
	}

	var firstTen int
	for i := 0; i < 10; i++ {
		firstTen += d[i]
	}
	if firstTen%10 != d[10] {
		return false
	}
	return true
}
