package main

import (
	_ "github.com/dcodix/grigori/lib/infrastructure/processor/processor_passthrough"
	_ "github.com/dcodix/grigori/lib/infrastructure/processor/processor_plain_to_json"
	_ "github.com/dcodix/grigori/lib/infrastructure/processor/processor_plain_to_logstash"
)
