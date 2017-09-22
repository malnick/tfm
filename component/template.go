package component

var (
	mainTemplate = `terraform {
  backend "s3" {
    bucket  = "{{ .RemoteState.S3Bucket }}"
    key     = "{{ .Name }}.tfstate"
    region  = "{{ .RemoteState.AWSRegion }}"
    profile = "{{ .RemoteState.AWSProfile }}"
  }
}`

	varsTemplate = `variable "deployment" {
	default = "{{ .NameSpace }}"
}`
)
