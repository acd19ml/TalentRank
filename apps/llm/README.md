
# 生成 serivce pb编译文件,acd19ml/TalentRank/apps/llm
```
protoc -I=. --go_out=. --go_opt=module="github.com/acd19ml/TalentRank/apps/llm" pb/llm.proto
```


# 补充rpc 接口定义protobuf的代码生成,acd19ml/TalentRank/apps/llm
```
protoc -I=. --go_out=. --go_opt=module="github.com/acd19ml/TalentRank/apps/llm" --go-grpc_out=. --go-grpc_opt=module="github.com/acd19ml/TalentRank/apps/llm" pb/llm.proto
```