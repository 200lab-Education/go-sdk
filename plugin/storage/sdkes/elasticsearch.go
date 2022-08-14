/*
 * @author           Viet Tran <viettranx@gmail.com>
 * @copyright       2019 Viet Tran <viettranx@gmail.com>
 * @license           Apache-2.0
 */

// github: https://github.com/olivere/elastic

package sdkes

import (
	"flag"
	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

const (
	DefaultElasticSearchName = "DefaultES"
)

type ESConfig struct {
	Prefix         string
	URL            string
	HasSniff       bool
	HasHealthCheck bool
	Index          string
	Username       string
	Password       string
}

type es struct {
	name   string
	client *elastic.Client
	logger logger.Logger
	*ESConfig
}

func NewES(name, flagPrefix string) *es {
	return &es{
		name: name,
		ESConfig: &ESConfig{
			Prefix: flagPrefix,
		},
	}
}

func (es *es) GetDefaultES() *es {
	return NewES(DefaultElasticSearchName, "")
}

func (es *es) GetPrefix() string {
	return es.Prefix
}

func (es *es) GetIndex() string {
	return es.Index
}

func (es *es) isDisabled() bool {
	return es.URL == ""
}

func (es *es) InitFlags() {
	prefix := es.Prefix
	if es.Prefix != "" {
		prefix += "-"
	}
	flag.StringVar(&es.Index, prefix+"es-index", "", "elastic search index")
	flag.StringVar(&es.URL, prefix+"es-url", "", "elastic search connection-string. ex: http://localhost:9200")
	flag.BoolVar(&es.HasSniff, prefix+"es-has-sniff", false, "elastic search sniffing mode. default: false")
	flag.BoolVar(&es.HasHealthCheck, prefix+"es-has-health-check", false, "elastic search health-check mode. default: false")
	flag.StringVar(&es.Username, prefix+"es-username", "", "elasticsearch username")
	flag.StringVar(&es.Password, prefix+"es-password", "", "elasticsearch password")
	flag.Parse()
}

func (es *es) Configure() error {
	if es.isDisabled() {
		return nil
	}
	es.logger = logger.GetCurrent().GetLogger(es.name)
	es.logger.Info("connecting to elastic search at ", es.URL, "...")

	var client *elastic.Client
	var err error
	if es.Username != "" && es.Password != "" {
		client, err = elastic.NewClient(
			elastic.SetURL(es.URL),
			elastic.SetBasicAuth(es.Username, es.Password),
			elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC ", log.LstdFlags)),
			elastic.SetSniff(es.HasSniff),
			elastic.SetHealthcheck(es.HasHealthCheck))
	} else {
		client, err = elastic.NewClient(
			elastic.SetURL(es.URL),
			elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC ", log.LstdFlags)),
			elastic.SetSniff(es.HasSniff),
			elastic.SetHealthcheck(es.HasHealthCheck))
	}

	if err != nil {
		es.logger.Error("cannot connect to elastic search. ", err.Error())
		return err
	}

	// Connect successfully, assign client
	es.client = client
	return nil
}

func (es *es) Name() string {
	return es.name
}

func (es *es) Get() interface{} {
	return es.client
}

func (es *es) Run() error {
	return es.Configure()
}

func (es *es) Stop() <-chan bool {

	c := make(chan bool)
	go func() {
		if es.client != nil {
			es.client.Stop()
		}
		c <- true
	}()
	return c
}
