package main

type Tree struct {
	Next  *Tree
	Value int
}

var shardData = map[string]interface{}{}

var shardData_sync = make(chan *data_register_block)

type data_register_block struct {
	location map[string]interface{}
	data     interface{}
	key      string
}

func register_data(loc map[string]interface{}, key string, data interface{}) { // send data to channel
	shardData_sync <- &data_register_block{location: loc, data: data, key: key}
}
func data_sync() {
	for {
		block := <-shardData_sync
		block.location[block.key] = block.data
		/*
			switch v := block.data.(type) {
			case string:
				block.location[block.key] = v
			case int:
				block.location[block.key] = v
			case int8:
				block.location[block.key] = v
			case int16:
				block.location[block.key] = v
			case int32:
				block.location[block.key] = v
			case int64:
				block.location[block.key] = v
			case int128:
				block.location[block.key] = v
			case float32:
				block.location[block.key] = v
			case float64:
				block.location[block.key] = v
			case float128:
				block.location[block.key] = v
			case decimal64:
				block.location[block.key] = v
			case decimal128:
				block.location[block.key] = v
			case decimal192:
				block.location[block.key] = v
			case bool:
				block.location[block.key] = v
			case byte:
				block.location[block.key] = v
			case []byte:
				block.location[block.key] = v
			case [][]byte:
				block.location[block.key] = v
			case []string:
				block.location[block.key] = v
			case []int:
				block.location[block.key] = v
			case []int8:
				block.location[block.key] = v
			case []int16:
				block.location[block.key] = v
			case []int32:
				block.location[block.key] = v
			case []int64:
				block.location[block.key] = v
			case []float32:
				block.location[block.key] = v
			case []float64:
				block.location[block.key] = v
			case []bool:
				block.location[block.key] = v
			case map[string]interface{}:
				block.location[block.key] = v
			}
		*/
	}

}
