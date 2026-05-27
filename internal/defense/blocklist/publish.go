// Package blocklist publishes the shared confirmed-scam blocklist in
// formats that adblockers, pi-hole, and DNS sinkholes consume.
package blocklist

import (
	"fmt"
	"io"
)

// HostsTxt writes a /etc/hosts-style blocklist.
func HostsTxt(w io.Writer, domains []string) error {
	if _, err := fmt.Fprintln(w, "# Joblantern confirmed scam recruiters"); err != nil {
		return err
	}
	for _, d := range domains {
		if _, err := fmt.Fprintf(w, "0.0.0.0 %s\n", d); err != nil {
			return err
		}
	}
	return nil
}

// PiHole writes a pi-hole adlist file (one domain per line).
func PiHole(w io.Writer, domains []string) error {
	for _, d := range domains {
		if _, err := fmt.Fprintln(w, d); err != nil {
			return err
		}
	}
	return nil
}

// uBlock writes a uBlock-style filter list.
func UBlock(w io.Writer, domains []string) error {
	for _, d := range domains {
		if _, err := fmt.Fprintf(w, "||%s^\n", d); err != nil {
			return err
		}
	}
	return nil
}
