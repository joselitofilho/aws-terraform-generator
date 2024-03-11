# example

<div style="text-align:center"><img src="example.drawio.svg" /></div>

- [diagram.xml](diagram.xml)
- [structure.config.yaml](structure.config.yaml)

## Commands
```bash
$ aws-terraform-generator diagram -s mystack -c diagram.config.yaml -d diagram.xml -o diagram.yaml
$ aws-terraform-generator structure -c structure.config.yaml -o ./output
$ aws-terraform-generator apigateway -c diagram.yaml -o ./output
$ aws-terraform-generator lambda -c diagram.yaml -o ./output/mystack
$ aws-terraform-generator kinesis -c diagram.yaml -o output/mystack
$ aws-terraform-generator sqs -c diagram.yaml -o output/mystack
$ aws-terraform-generator s3 -c diagram.yaml -o output/mystack
```