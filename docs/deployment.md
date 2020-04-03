# Keymaster Deployment for CI

This document assumes you are deploying the full workflow-based
system for use with CI.

## Software environments

The change target of a CI deployment operation will be referred 
to as an "environment". Ideally each environment will have it's own
unique cloud account. This is not required, however. For example,
"test" and "staging" might be tenants in the same cloud account.

## Workflow engine 

The workflow engine is a Java (Spring Boot) application, not open
source at this time. It is packaged as a Docker image and is designed
to run in AWS. Kubernetes is the primary target for development but 
any runtime that supports Docker containers will work (e.g. Fargate).

Keymaster is designed such that the workflow engine is not fully 
trusted. This means that the issuing lambdas do their own validation
 of SAML assertions and do not fully rely on the security of workflow.
 
Nevertheless, the workflow engine should be deployed and managed in 
a way commensurate with the compliance requirements of tenants. In
practice this means it must meet the union of the requirements of 
all tenants.

It is advisable to deploy two instances of the workflow engine,
for production vs non-production workflows. This will allow testing
of upgrades and changes to the workflow engine with less chance of 
disrupting production deployments.

The workflow engine requires:

* A docker
* An HTTPS ingress with an associated domain name
* DynamoDB access (for managing approval state)

It does not require:

* Direct connectivity to any target software environments.

Terraform code will be provided to deploy the workflow engine.

At this time the workflow engine does not require any environment
specific configuration - workflow is driven by configuration
provided by the issuing lambda (passed via the client).

## Identity Provider (IDP)

You will need an identity provider. At this stage only SAML IDPs
are supported, although OpenID connect is on our roadmap.

One federation will be required between the IDP and the workflow
engine. Items to be deployed are:

* One SAML federation between IDP and workflow engine

## Deployment roles, keys and resources

IAM roles which support the required CI changes will need to be 
created for each target software environment.

The keymaster issuing lambda for the target environment will need
to be granted permission to assume the relevant roles.

In high security environments issued credentials maybe wrapped
with a KMS key. 

The keymaster terraform code provides an example.

Items to be deployed are:

* Low-privilege deployment roles for the CI runner(s)
* High-privilege deployment roles for the CI runner(s)
* One or more "credential wrapping keys"
* Permission for the low-privilege deployment role to Decrypt
  using the credential wrapping key
* Permission for the km issuing lambda to Encrypt using the 
  credential wrapping key
* Assume-role policies allowing km issuance to assume the roles

## CI Runners

In situations with relaxed security requirements, shared or
pooled CI runners may use km authentication to get deployment
credentials either via workflow or non-interactively (effectively,
with automatic approval).

In environments with specific compliance requirements, a dedicated
runner should be deployed to make changes. One logical place to
put that runner is inside the target environment.  This is the 
simplest option in the case that deployments need to interact not 
only with cloud resources, but also resources (kube masters, 
internal services, data stores in the environment) that can not
be accessed from outside the target environment itself. 

Items to be deployed are:

* One CI runner per target environment
  * Configured to run in the low-privilege CI runner role

## Keymaster issuing lambda

We recommend to deploy one keymaster instance for each unique
software environment that requires credential issuance. So for
example, if there are two environments (say, "test" and "staging")
then an issuing lambda should be deployed for each. 

This is not a requirement. It is possible to deploy a
single issuing lambda that supports an arbitrarily large number
of roles, workflows and environments. This may simplify some 
administration but the "blast radius" for changes will be larger 
and the overall system may be more difficult to secure.

For the highest levels of security, deploy the issuing lamdbda 
into a dedicated account which is not accessible to the deployment
roles of the target environment. 

The keymaster issuing lambda deployment includes the following
resources:

* The keymaster issuing lambda
* Configuration for the keymaster issuing lambda
  * This may include access to SSH CA keys if ssh CA issuance
    is used.

Provisioning code is provided in the km terraform folder.

## IP Oracle lamdba

If IP whitelisting is configured on the km issuing lambda, you
will also need to deploy the IP Oracle lambda. 

This is a small, stateless lambda that looks at the caller's 
IP address and issues a signed token containing that in a claim.

This allows a lambda called via "invoke" to 

The IP oracle deployment includes the following resources:

* An API gateway with an associated domain name
* A KMS key to use for "IP token" signing
* A lambda deployment

Provisioning code is provided in the km terraform folder.