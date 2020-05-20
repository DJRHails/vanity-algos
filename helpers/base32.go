package helpers

import (
  "strings"
)

const Base32RuneSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

// Check string only contains characters in Base32RuneSet
func IsBase32(s string) bool {
  for _, c := range s {
    if !strings.Contains(Base32RuneSet, string(c)) {
      return false
    }
  }
  return true
}
