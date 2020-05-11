package cmd

import "time"

// Conf for rundeck job manifest Conf data type
type Conf struct {
	Elasticsearch string `yaml:"elasticsearch"`
	Kibana        string `yaml:"kibana"`
	IndexPattern  string `yaml:"indexPattern"`
	KibanaMaxPage string `yaml:"kibanaMaxPage"`
	Interval      string `yaml:"interval"`
}

// Attributes for attributes body
type Attributes struct {
	Title         string `json:"title"`
	TimeFieldName string `json:"timeFieldName"`
}

// IndexPatternBody for cretate new index - POST request
type IndexPatternBody struct {
	Attributes Attributes `json:"attributes"`
}

// EsIndex for elasticsearch index list response
type EsIndex []struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}

// EsIndexPattern for kibana index pattern
type EsIndexPattern struct {
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	Total        int `json:"total"`
	SavedObjects []struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			Title         string `json:"title"`
			TimeFieldName string `json:"timeFieldName"`
			Fields        string `json:"fields"`
		} `json:"attributes,omitempty"`
		References       []interface{} `json:"references"`
		MigrationVersion struct {
			IndexPattern string `json:"index-pattern"`
		} `json:"migrationVersion"`
		UpdatedAt time.Time `json:"updated_at"`
		Version   string    `json:"version"`
	} `json:"saved_objects"`
}

// Banner beautiful
const Banner string = `
      ___                       ___           ___           ___           ___           ___     
     /\__\                     /\  \         /\__\         /\  \         /\  \         /\  \    
    /:/ _/_         ___        \:\  \       /:/  /         \:\  \       /::\  \       /::\  \   
   /:/ /\  \       /|  |        \:\  \     /:/  /           \:\  \     /:/\:\__\     /:/\:\  \  
  /:/ /::\  \     |:|  |    _____\:\  \   /:/  /  ___   ___ /::\  \   /:/ /:/  /    /:/  \:\  \ 
 /:/_/:/\:\__\    |:|  |   /::::::::\__\ /:/__/  /\__\ /\  /:/\:\__\ /:/_/:/__/___ /:/__/ \:\__\
 \:\/:/ /:/  /  __|:|__|   \:\~~\~~\/__/ \:\  \ /:/  / \:\/:/  \/__/ \:\/:::::/  / \:\  \ /:/  /
  \::/ /:/  /  /::::\  \    \:\  \        \:\  /:/  /   \::/__/       \::/~~/~~~~   \:\  /:/  / 
   \/_/:/  /   ~~~~\:\  \    \:\  \        \:\/:/  /     \:\  \        \:\~~\        \:\/:/  /  
     /:/  /         \:\__\    \:\__\        \::/  /       \:\__\        \:\__\        \::/  /   
     \/__/           \/__/     \/__/         \/__/         \/__/         \/__/         \/__/    

`
