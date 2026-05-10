// SPDX-License-Identifier: LGPL-2.1-or-later

package slurm

// SlurmOptionalInt models Slurm's "optional number" JSON shape:
//
//	{"set": true, "infinite": false, "number": N}
//
// When Set is false the value is absent/unknown; when Infinite is true the
// value is unlimited.  Number is float64 to accommodate fractional billing
// TRES values.
//
// This type is shared between job and node JSON decoders.
type SlurmOptionalInt struct {
	Set      bool    `json:"set"`
	Infinite bool    `json:"infinite"`
	Number   float64 `json:"number"`
}
