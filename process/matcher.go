package process

import (
  "fmt"
	"os"
	"context"
	"regexp"

  "github.com/schollz/progressbar/v3"
	"github.com/algorand/go-algorand-sdk/crypto"
)

type Matcher struct {
  Regex     *regexp.Regexp
	Prefix    string

  Bar       *progressbar.ProgressBar
	Results   chan *crypto.Account
}

func (m *Matcher) Run(ctx context.Context, semaphore chan bool) {
  go func() {
    defer func() {
      if err := recover(); err != nil {
        fmt.Fprintln(os.Stderr, "\nRecovered from unexpected error:", err)
      }

      // Clear the semaphore so another process can run
      <-semaphore
    }()

    for {
      select {
      case <-ctx.Done():
        // Application is ending
        return

      default:
        if match := m.Match(crypto.GenerateAccount()); match != nil {
          // Found a match!
          m.Results <- match
        }
      }
    }
  }()
}

func (m *Matcher) Match(account crypto.Account) *crypto.Account {
  address := account.Address.String()
  m.Bar.Add(1)
  if m.find(address) {
    return &account
  }
  return nil
}

func (m *Matcher) find(address string) bool {
  // Check the prefix first if specified
  if n := len(m.Prefix); n > 0 && address[:n] != m.Prefix {
    return false
  }

  // Check the regex next if present
  if m.Regex != nil && !m.Regex.MatchString(address) {
    return false
  }

  // If the prefix didn't fail, and the regex didn't fail, then is a match.
  return true
}
