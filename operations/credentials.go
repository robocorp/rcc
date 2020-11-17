package operations

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/xviper"
)

const (
	defaultsAccount  = `defaults.account`
	accountsPrefix   = `accounts.`
	endpointSuffix   = `.endpoint`
	identifierSuffix = `.identifier`
	labelSuffix      = `.label`
	secretSuffix     = `.secret`
	verifiedSuffix   = `.verified`
	detailsSuffix    = `.details`
)

var (
	dynamicAccountPattern = regexp.MustCompile("^(\\d+):([0-9a-f]{96})(?::(https://.+))?$")
)

type accountList []*account
type account struct {
	Account    string                 `json:"account"`
	Identifier string                 `json:"identifier"`
	Endpoint   string                 `json:"endpoint"`
	Secret     string                 `json:"secret"`
	Verified   int64                  `json:"verified"`
	Default    bool                   `json:"default"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

func DefaultAccountName() string {
	return xviper.GetString(defaultsAccount)
}

func SetDefaultAccount(account string) {
	xviper.Set(defaultsAccount, account)
}

func UpdateCredentials(account, endpoint, identifier, secret string) {
	if len(DefaultAccountName()) == 0 {
		SetDefaultAccount(account)
	}
	prefix := accountsPrefix + account
	xviper.Set(prefix+labelSuffix, account)
	xviper.Set(prefix+identifierSuffix, identifier)
	xviper.Set(prefix+secretSuffix, secret)
	xviper.Set(prefix+verifiedSuffix, 0)
	xviper.Set(prefix+detailsSuffix, new(map[string]interface{}))
	if len(endpoint) > 0 {
		xviper.Set(prefix+endpointSuffix, endpoint)
	}
}

func VerifyAccounts(force bool) {
	marker := time.Now().Add(-2 * time.Hour)
	if marker.Before(common.Startup) {
		marker = common.Startup
	}
	deadline := marker.Unix()
	for _, entry := range findAccounts() {
		recently := entry.Verified > deadline
		_, detailed := entry.Details["email"]
		if !force && detailed && recently {
			continue
		}
		entry.WasVerified(0)
		client, err := cloud.NewClient(entry.Endpoint)
		if err != nil {
			continue
		}
		UserinfoCommand(client, entry)
	}
}

func (it *account) CacheKey() string {
	return fmt.Sprintf("%s.%s", it.Identifier, it.Secret[:6])
}

func (it *account) CacheToken(name, url, token string, deadline int64) {
	if common.NoCache {
		return
	}
	cache, err := SummonCache()
	if err != nil {
		return
	}
	defer cache.Save()
	credential := Credential{
		Account:  it.Account,
		Context:  name,
		Token:    token,
		Deadline: deadline,
	}
	fullkey := strings.ToLower(it.CacheKey() + "/" + name + "/" + url)
	cache.Credentials[fullkey] = &credential
}

func (it *account) Cached(name, url string) (string, bool) {
	if common.NoCache {
		return "", false
	}
	cache, err := SummonCache()
	if err != nil {
		return "", false
	}
	fullkey := strings.ToLower(it.CacheKey() + "/" + name + "/" + url)
	found, ok := cache.Credentials[fullkey]
	if !ok {
		return "", false
	}
	if found.Deadline < time.Now().Unix() {
		return "", false
	}
	return found.Token, true
}

func (it *account) Delete() error {
	prefix := accountsPrefix + it.Account
	defer xviper.Set(prefix, "deleted")

	client, err := cloud.NewClient(it.Endpoint)
	if err != nil {
		return err
	}
	return DeleteAccount(client, it)
}

func (it *account) WasVerified(when int64) {
	prefix := accountsPrefix + it.Account
	xviper.Set(prefix+verifiedSuffix, when)
}

func (it *account) UpdateDetails(details Token) {
	it.Details = details
	prefix := accountsPrefix + it.Account
	xviper.Set(prefix+detailsSuffix, details)
}

func listAccountsAsJson(accounts accountList) {
	body, err := NiceJsonOutput(accounts)
	if err != nil {
		common.Error("list-accounts", err)
	} else {
		common.Out("%s", body)
	}
}

func listAccountsAsText(accounts accountList) {
	if len(accounts) == 0 {
		common.Log("No account information available.")
		return
	}
	common.Log("Identifier    Account             Default  Secret               Valid  Endpoint")
	for _, entry := range accounts {
		verified := entry.Verified > 1000
		common.Log("  %-10s  %-20s  %-5v  %19s  %-5v  %-s", entry.Identifier, entry.Account, entry.Default, entry.Secret, verified, entry.Endpoint)
	}
}

func smudgeSecrets(accounts accountList) accountList {
	for _, account := range accounts {
		if len(account.Secret) > 19 {
			account.Secret = account.Secret[:16] + "..."
		}
	}
	return accounts
}

func ListAccounts(json bool) {
	accounts := smudgeSecrets(findAccounts())
	if json {
		listAccountsAsJson(accounts)
	} else {
		listAccountsAsText(accounts)
	}
}

func EncodeCredentials(target *json.Encoder, force bool) error {
	VerifyAccounts(force)
	return target.Encode(smudgeSecrets(findAccounts()))
}

func loadAccount(label string) *account {
	prefix := accountsPrefix + label
	var details Token
	received := xviper.Get(prefix + detailsSuffix)
	some, ok := received.(map[string]interface{})
	if !ok {
		some, ok = received.(Token)
	}
	if ok {
		details = some
	}
	return &account{
		Account:    xviper.GetString(prefix + labelSuffix),
		Identifier: xviper.GetString(prefix + identifierSuffix),
		Endpoint:   xviper.GetString(prefix + endpointSuffix),
		Secret:     xviper.GetString(prefix + secretSuffix),
		Verified:   xviper.GetInt64(prefix + verifiedSuffix),
		Details:    details,
	}
}

func createEphemeralAccount(parts []string) *account {
	BackgroundMetric("rcc", "rcc.account.ephemeral", common.Version)
	common.NoCache = true
	endpoint := common.DefaultEndpoint
	if len(parts[3]) > 0 {
		endpoint = parts[3]
	}
	return &account{
		Account:    "Ephemeral",
		Identifier: parts[1],
		Secret:     parts[2],
		Endpoint:   endpoint,
		Verified:   0,
		Default:    false,
		Details:    make(map[string]interface{}),
	}
}

func AccountByName(label string) *account {
	dynamic := dynamicAccountPattern.FindStringSubmatch(label)
	if dynamic != nil {
		return createEphemeralAccount(dynamic)
	}
	if len(label) == 0 {
		label = DefaultAccountName()
	}
	found := loadAccount(label)
	if found.Account == label {
		return found
	}
	return nil
}

func findAccounts() accountList {
	accounts := make([]string, 0, 10)
	for _, key := range xviper.AllKeys() {
		if strings.HasPrefix(key, accountsPrefix) && strings.HasSuffix(key, ".label") {
			accounts = append(accounts, xviper.GetString(key))
		}
	}
	sort.Strings(accounts)
	result := make(accountList, 0, len(accounts))
	defaultAccount := DefaultAccountName()
	for _, key := range accounts {
		here := loadAccount(key)
		if here.Account == defaultAccount {
			here.Default = true
		}
		result = append(result, here)
	}
	return result
}
