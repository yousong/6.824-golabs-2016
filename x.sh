#
# How to use this
#
# ➜  ~/dev/6.824-golabs-2016/src/main git:(yousong) sh -x ../../x.sh sh ./test-mr.sh
# + GOPATH=/Users/yousong/dev/6.824-golabs-2016
# + sh ./test-mr.sh
# 
# ==> Part I
# ok      mapreduce       3.953s
# 
# ==> Part II
# Passed test
# 
# ==> Part III
# ok      mapreduce       2.421s
# 
# ==> Part IV
# ok      mapreduce       11.942s
# 
# ==> Part V (challenge)
# Passed test
# ➜  ~/dev/6.824-golabs-2016/src/mapreduce git:(yousong) sh -x ../../x.sh go test -run TestBasic mapreduce
# + GOPATH=/Users/yousong/dev/6.824-golabs-2016
# + go test -run TestBasic mapreduce
# ok      mapreduce       4.192s
# ➜  ~/dev/6.824-golabs-2016/src/mapreduce git:(yousong) sh -x ../../x.sh go test  mapreduce
# + GOPATH=/Users/yousong/dev/6.824-golabs-2016
# + go test mapreduce
# ok      mapreduce       18.589s
#
GOPATH="$HOME/dev/6.824-golabs-2016" "$@"
