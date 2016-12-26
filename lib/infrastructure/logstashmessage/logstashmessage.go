package logstashmessage

type LogstashMessage struct {
	logstash_version   string
	logstash_tags      []string
	logstash_host      string
	logstash_file      string
	logstash_type      string
	logstash_timestamp string
	logstash_message   string
}
