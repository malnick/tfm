# Genesis
Genesis simplifies Terraform deployments by instituting a well-defined pattern around deployments, the components that make up those deployments, and the common workflows around them. 

## tfm Command
`tfm` uses the genesis library to make a simple command line wrapper for CRUD operations on deployments and their components. 

```
NAME:
   tfm - Terraform wrapper for implementing deployment and component abstraction

USAGE:
   tfm-darwin-.44b5060 [global options] command [command options] [arguments...]

VERSION:
   UNSET

AUTHOR:
   Jeff Malnick <malnick at google mail>

COMMANDS:
     create   Make a new component for a deployment
     plan     Run a terraform plan for a component
     apply    Run a terraform apply for a component
     list     List all components
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        Run in debug mode
   --help, -h     show help
   --version, -v  print the version
```

## Deployments and Components Pattern
Terraform makes it easy to codify interactions with cloud and other API interfaces. Most people use Terraform for codifying their infrastructure on a cloud platform. Genesis takes Terraform and institutes common patterns that make it easier to organize, maintain and scale your Terraform codebases. The top level abstractions Genesis relies are `deployment` and `component`. A `deployment` abstraction is the business description of a set of components. A component is a specific subset of a deployment that is independent.

### Deployment
For some companies, a deployment could be a VPC on AWS. A VPC is good example of a deployment since it encompasses many independent components. However, a deployment does not need to be a single VPC, it could be many, it's up to the operator. The idea is to have a top level description that identifies a logical center to store components. 

### Component
A component is a individual set of cloud resources that are independent from their siblingins inside a deployment. An example component may be a Cassandra cluster. This consists of many EC2 resources, security groups, etc. Components are identified by their need to be separated from their sibling components with terraform remote state. This reduces the chances of one change adversely affecting other components. 
