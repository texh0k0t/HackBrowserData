package core

import (
	"encoding/base64"
	"errors"

	"hack-browser-data/core/decrypt"
	"hack-browser-data/utils"

	"github.com/tidwall/gjson"
)

const (
	chromeProfilePath     = "/AppData/Local/Google/Chrome/User Data/*/"
	chromeKeyPath         = "/AppData/Local/Google/Chrome/User Data/Local State"
	chromeBetaProfilePath = "/AppData/Local/Google/Chrome Beta/User Data/*/"
	chromeBetaKeyPath     = "/AppData/Local/Google/Chrome Beta/User Data/Local State"
	edgeProfilePath       = "/AppData/Local/Microsoft/Edge/User Data/*/"
	edgeKeyPath           = "/AppData/Local/Microsoft/Edge/User Data/Local State"
	braveProfilePath      = "/AppData/Local/BraveSoftware/Brave-Browser/User Data/*/"
	braveKeyPath          = "/AppData/Local/BraveSoftware/Brave-Browser/User Data/Local State"
	speed360ProfilePath   = "/AppData/Local/360chrome/Chrome/User Data/*/"
	qqBrowserProfilePath  = "/AppData/Local/Tencent/QQBrowser/User Data/*/"
	firefoxProfilePath    = "/AppData/Roaming/Mozilla/Firefox/Profiles/*.default-release/"
)

var (
	browserList = map[string]struct {
		ProfilePath string
		Name        string
		KeyPath     string
		Storage     string
		New         func(profile, key, name, storage string) (Browser, error)
	}{
		"chrome": {
			ProfilePath: userProfile + chromeProfilePath,
			KeyPath:     userProfile + chromeKeyPath,
			Name:        chromeName,
			New:         NewChromium,
		},
		"chrome-beta": {
			ProfilePath: userProfile + chromeBetaProfilePath,
			KeyPath:     userProfile + chromeBetaKeyPath,
			Name:        chromeBetaName,
			New:         NewChromium,
		},
		"edge": {
			ProfilePath: userProfile + edgeProfilePath,
			KeyPath:     userProfile + edgeKeyPath,
			Name:        edgeName,
			New:         NewChromium,
		},
		"360": {
			ProfilePath: userProfile + speed360ProfilePath,
			Name:        speed360Name,
			New:         NewChromium,
		},
		"qq": {
			ProfilePath: userProfile + qqBrowserProfilePath,
			Name:        qqBrowserName,
			New:         NewChromium,
		},
		"firefox": {
			ProfilePath: userProfile + firefoxProfilePath,
			Name:        firefoxName,
			New:         NewFirefox,
		},
		"brave": {
			ProfilePath: userProfile + braveProfilePath,
			KeyPath:     userProfile + braveKeyPath,
			Name:        braveName,
			New:         NewChromium,
		},
	}
)

var (
	errBase64DecodeFailed = errors.New("decode base64 failed")
)

// InitSecretKey on windows with win32 DPAPI
// conference from @https://gist.github.com/akamajoris/ed2f14d817d5514e7548
func (c *Chromium) InitSecretKey() error {
	if c.keyPath == "" {
		return nil
	}
	keyFile, err := utils.ReadFile(c.keyPath)
	if err != nil {
		return err
	}
	encryptedKey := gjson.Get(keyFile, "os_crypt.encrypted_key")
	if encryptedKey.Exists() {
		pureKey, err := base64.StdEncoding.DecodeString(encryptedKey.String())
		if err != nil {
			return errBase64DecodeFailed
		}
		c.secretKey, err = decrypt.DPApi(pureKey[5:])
		return err
	}
	return nil
}
