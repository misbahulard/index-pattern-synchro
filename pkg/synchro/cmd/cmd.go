package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var cfg Conf

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

// Run index-pattern synchro runtime
func Run() {
	fmt.Println(Banner)

	Config(&cfg)
	fmt.Println("Config")
	fmt.Println("----------------")
	fmt.Printf("elasticsearch \t: %s\n", cfg.Elasticsearch.Host)
	fmt.Printf("kibana \t\t: %s\n", cfg.Kibana.Host)
	fmt.Printf("index-pattern \t: %s\n", cfg.Elasticsearch.IndexPattern)
	fmt.Printf("kibana max page : %s\n", cfg.Kibana.MaxPage)
	fmt.Printf("interval \t: %s minutes\n\n", cfg.Interval)

	intInterval, _ := strconv.ParseInt(cfg.Interval, 10, 64)
	interval := "@every " + cfg.Interval + "m"

	schedule := time.Now().Add(time.Minute * time.Duration(intInterval))
	log.Info("next execute: " + schedule.Format("2006-01-02 15:04:05"))

	c := cron.New()
	c.AddFunc(interval, watchJob)
	c.Run()
}

func watchJob() {
	log.Info("synchronizing...")
	listIndices := getIndices()
	listIndexPatterns := getIndexPatterns()
	finalIndices := findDifferences(listIndices, listIndexPatterns)

	if len(finalIndices) == 0 {
		log.Info("all clean captain")
	} else {
		log.Info("new indices:")

		for _, index := range finalIndices {
			log.Info("- " + index)
		}

		createIndexPattern(finalIndices)
	}

	intInterval, _ := strconv.ParseInt(cfg.Interval, 10, 64)

	schedule := time.Now().Add(time.Minute * time.Duration(intInterval))
	log.Info("next execute: " + schedule.Format("2006-01-02 15:04:05"))
}

func createIndexPattern(slice []string) {
	log.Info("create index pattern: ")
	client := resty.New()

	if cfg.Kibana.Auth {
		client.SetBasicAuth(cfg.Kibana.User, cfg.Kibana.Password)
	}

	for _, item := range slice {
		url := cfg.Kibana.Host + "/api/saved_objects/index-pattern/" + item

		attributes := Attributes{
			Title:         item + "-*",
			TimeFieldName: "@timestamp",
		}

		indexPatternBody := IndexPatternBody{
			Attributes: attributes,
		}

		resp, err := client.R().
			SetHeader("kbn-xsrf", "true").
			SetHeader("Content-Type", "application/json").
			SetBody(indexPatternBody).
			Post(url)

		if err != nil {
			log.Fatalln(err)
		}

		log.Infof("- %s --- %s", item, resp.Body())
	}
}

func getIndices() []string {
	url := cfg.Elasticsearch.Host + "/_cat/indices?format=json"

	client := resty.New()

	if cfg.Elasticsearch.Auth {
		client.SetBasicAuth(cfg.Elasticsearch.User, cfg.Elasticsearch.Password)
	}

	resp, err := client.R().Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	var esIndex EsIndex
	json.Unmarshal(resp.Body(), &esIndex)

	list := []string{}
	re, _ := regexp.Compile("^" + cfg.Elasticsearch.IndexPattern)

	for _, index := range esIndex {
		matched := re.MatchString(index.Index)

		if matched {
			chunks := strings.Split(index.Index, "-")
			indexName := strings.Join(chunks[:len(chunks)-1], "-")
			list = append(list, indexName)
		}
	}

	return unique(list)
}

func getIndexPatterns() []string {
	url := cfg.Kibana.Host + "/api/saved_objects/_find?type=index-pattern&per_page=" + cfg.Kibana.MaxPage

	client := resty.New()

	if cfg.Kibana.Auth {
		client.SetBasicAuth(cfg.Kibana.User, cfg.Kibana.Password)
	}

	resp, err := client.R().Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	var esIndexPattern EsIndexPattern
	json.Unmarshal(resp.Body(), &esIndexPattern)

	list := []string{}

	ip := strings.Split(cfg.Elasticsearch.IndexPattern, "-")
	ipClean := strings.Join(ip[:len(ip)-1], "-")

	re, _ := regexp.Compile(cfg.Elasticsearch.IndexPattern)

	for _, object := range esIndexPattern.SavedObjects {
		matched := re.MatchString(object.Attributes.Title)

		if matched {
			chunks := strings.Split(object.Attributes.Title, "-")
			indexName := strings.Join(chunks[:len(chunks)-1], "-")
			if indexName != ipClean {
				list = append(list, indexName)
			}
		}
	}

	return unique(list)
}

func findDifferences(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func restructureIndex(slice []string) []string {
	list := []string{}

	for _, entry := range slice {
		list = append(list, entry+"-*")
	}

	return list
}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
