package mapreduce

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	//"sort"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	var kvs []*KeyValue
	var kvmap = make(map[string][]string)

	// decode content from each map task
	for i := 0; i < nMap; i++ {
		fn := reduceName(jobName, i, reduceTaskNumber)
		rfile, err := os.Open(fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s failed: %s\n", fn, err)
			return
		}
		defer rfile.Close()
		dec := json.NewDecoder(rfile)
		for {
			var kv KeyValue
			err = dec.Decode(&kv)
			if err == nil {
				kvs = append(kvs, &kv)
			} else if err == io.EOF {
				break
			} else {
				fmt.Fprintf(os.Stderr, "decode %s failed: %s\n", fn, err)
				return
			}
		}
	}

	for _, kv := range kvs {
		if kvmap[kv.Key] != nil {
			kvmap[kv.Key] = append(kvmap[kv.Key], kv.Value)
		} else {
			kvmap[kv.Key] = []string{kv.Value}
		}
	}

	// make mergeFile
	fn := mergeName(jobName, reduceTaskNumber)
	mergeFile, err := os.Create(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create %s failed: %s\n", fn, err)
		return
	}
	defer mergeFile.Close()
	enc := json.NewEncoder(mergeFile)

	// call reduce and write result
	for key, values := range kvmap {
		r := reduceF(key, values)
		enc.Encode(KeyValue{key, r})
	}
}
