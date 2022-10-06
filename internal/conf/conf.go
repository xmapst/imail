package conf

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"

	"github.com/midoks/imail/internal/assets/conf"
	"github.com/midoks/imail/internal/tools"
)

// Asset is a wrapper for getting conf assets.
func Asset(name string) ([]byte, error) {
	return conf.Asset(name)
}

// AssetDir is a wrapper for getting conf assets.
func AssetDir(name string) ([]string, error) {
	return conf.AssetDir(name)
}

// MustAsset is a wrapper for getting conf assets.
func MustAsset(name string) []byte {
	return conf.MustAsset(name)
}

// File is the configuration object.
var File *ini.File

func Init(customConf string) error {

	File, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, conf.MustAsset("conf/app.conf"))
	if err != nil {
		return errors.Wrap(err, "parse 'conf/app.conf'")
	}

	File.NameMapper = ini.TitleUnderscore

	if customConf == "" {
		customConf = filepath.Join(CustomDir(), "conf", "app.conf")
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return errors.Wrap(err, "get absolute path")
		}
	}
	CustomConf = customConf

	if tools.IsFile(customConf) {
		if err = File.Append(customConf); err != nil {
			return errors.Wrapf(err, "append %q", customConf)
		}
	} else {
		log.Println("Custom config ", customConf, " not found. Ignore this warning if you're running for the first time")
	}

	if err = File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return errors.Wrap(err, "mapping default section")
	}

	// ***************************
	// ----- Log settings -----
	// ***************************
	if err = File.Section("log").MapTo(&Log); err != nil {
		return errors.Wrap(err, "mapping [log] section")
	}

	// ***************************
	// ----- Database settings -----
	// ***************************
	if err = File.Section("database").MapTo(&Database); err != nil {
		return errors.Wrap(err, "mapping [log] section")
	}

	// ****************************
	// ----- Web settings -----
	// ****************************

	if err = File.Section("web").MapTo(&Web); err != nil {
		return errors.Wrap(err, "mapping [web] section")
	}

	Web.AppDataPath = ensureAbs(Web.AppDataPath)

	if !strings.HasSuffix(Web.ExternalURL, "/") {
		Web.ExternalURL += "/"
	}
	Web.URL, err = url.Parse(Web.ExternalURL)
	if err != nil {
		return errors.Wrapf(err, "parse '[server] EXTERNAL_URL' %q", err)
	}

	// Subpath should start with '/' and end without '/', i.e. '/{subpath}'.
	Web.Subpath = strings.TrimRight(Web.URL.Path, "/")
	Web.SubpathDepth = strings.Count(Web.Subpath, "/")

	unixSocketMode, err := strconv.ParseUint(Web.UnixSocketPermission, 8, 32)
	if err != nil {
		return errors.Wrapf(err, "parse '[server] unix_socket_permission' %q", Web.UnixSocketPermission)
	}
	if unixSocketMode > 0777 {
		unixSocketMode = 0666
	}
	Web.UnixSocketMode = os.FileMode(unixSocketMode)

	// ****************************
	// ----- Session settings -----
	// ****************************

	if err = File.Section("session").MapTo(&Session); err != nil {
		return errors.Wrap(err, "mapping [session] section")
	}

	// ***********************************
	// ----- Authentication settings -----
	// ***********************************

	if err = File.Section("auth").MapTo(&Auth); err != nil {
		return errors.Wrap(err, "mapping [auth] section")
	}

	// ***************************
	// ----- SMTP settings -----
	// ***************************
	if err = File.Section("smtp").MapTo(&Smtp); err != nil {
		return errors.Wrap(err, "mapping [smtp] section")
	}

	// ***************************
	// ----- Pop3 settings -----
	// ***************************
	if err = File.Section("pop3").MapTo(&Pop3); err != nil {
		return errors.Wrap(err, "mapping [pop] section")
	}

	// ***************************
	// ----- Imap settings -----
	// ***************************
	if err = File.Section("imap").MapTo(&Imap); err != nil {
		return errors.Wrap(err, "mapping [imap] section")
	}

	// ***************************
	// ----- Rspamd settings -----
	// ***************************
	if err = File.Section("rspamd").MapTo(&Rspamd); err != nil {
		return errors.Wrap(err, "mapping [rspamd] section")
	}

	// ***************************
	// ----- Ssl settings -----
	// ***************************
	if err = File.Section("ssl").MapTo(&Ssl); err != nil {
		return errors.Wrap(err, "mapping [ssl] section")
	}

	// *****************************
	// ----- Security settings -----
	// *****************************

	if err = File.Section("security").MapTo(&Security); err != nil {
		return errors.Wrap(err, "mapping [security] section")
	}

	// ***************************
	// ----- i18n settings -----
	// ***************************
	I18n = new(i18nConf)
	if err = File.Section("i18n").MapTo(&I18n); err != nil {
		return errors.Wrap(err, "mapping [i18n] section")
	}

	if err = File.Section("cache").MapTo(&Cache); err != nil {
		return errors.Wrap(err, "mapping [cache] section")
	} else if err = File.Section("other").MapTo(&Other); err != nil {
		return errors.Wrap(err, "mapping [other] section")
	}

	// Check run user when the install is locked.
	if Security.InstallLock {
		currentUser, match := CheckRunUser(App.RunUser)
		if !match {
			return fmt.Errorf("user configured to run imail is %q, but the current user is %q", App.RunUser, currentUser)
		}
	}

	return nil
}
