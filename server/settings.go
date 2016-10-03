package server

import "github.com/naoina/toml"
import "os"
import "io/ioutil"
import "crypto/tls"

type ParadiseSettings struct {
	Host           string
	Port           int
	MaxConnections int
	MaxPassive     int
	Exec           string
	Pem            string
	Key            string
}

func Load509Config() *tls.Config {
	// use https://letsencrypt.org to get the pem and key files
	cert, cerr := tls.LoadX509KeyPair(Settings.Pem, Settings.Key)
	if cerr != nil {
		return nil
	}

	config := &tls.Config{}
	if config.NextProtos == nil {
		config.NextProtos = []string{"ftp"}
	}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0] = cert

	return config
}

func ReadSettings() ParadiseSettings {
	f, err := os.Open("conf/settings.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config ParadiseSettings
	if err := toml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}
	return config
}
