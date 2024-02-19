# example

<div style="text-align:center"><img src="example.drawio.png" /></div>

- [diagram.drawio.xml](diagram.drawio.xml)
- [structure.yaml](structure.yaml)

## Commands
```bash
$ aws-terraform-generator diagram -s mystack -i diagram.drawio.xml -o mystack.yaml
$ aws-terraform-generator structure -i structure.yaml -o ./output
$ aws-terraform-generator apigateway -i mystack.yaml -o ./output
$ aws-terraform-generator lambda -i mystack.yaml -o ./output/mystack
$ aws-terraform-generator sqs -i mystack.yaml -o output/mystack/mod/sqs.tf
$ aws-terraform-generator s3 -i mystack.yaml -o output/mystack/mod/s3.tf
```