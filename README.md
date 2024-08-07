# About

This client is meant to add and remove several concurrent versions of the same service, corresponding to various ongoing feature branches, in a gitops manner.

Currently, both kubernetes fluxcd and terraform orchestration are supported across several repos.

# Usage

## Playground Repo

As you are reading this documentation, refer to the followin resources for a playground that can be experimented with when discovering this project:
- **Playground repo**: https://github.com/Ferlab-Ste-Justine/ferlease-playground
- **Configuration example**: https://github.com/Ferlab-Ste-Justine/ferlease/tree/main/usage-example

You will have to clone the playground repo and adjust the configuration accordingly.

From there, you should be able to play with the **release** and **teardown** commands.

## Assumption

### Fluxcd

This project assumes that any fluxcd configuration for the service to be released will be broken down into the following parts:
  - An **apps** directory containing the kubernetes resources of the service to deploy
  - An orchestration file containing all fluxcd resources to manage the app (minimally a **fluxcd kustomization** resource pointing to the **apps** directory)
  - A **kustomize kustomization** file that manages all the fluxcd resource files of various releases to be deployed

### Terraform

This project assumes that any terraform configuration for the service to be released will be broken down into the following parts:
  - An entrypoint terraform orchestration file
  - An optional terraform module inlined in a subdirectory of the directory containing the entrypoint file which the entrypoint file can refer to

## Templatization Support

This project will orchestrate releases based on a provided template that follows the golang templating format.

Currently, the following variables are supported in the template:
- **Service**: Name of the service to deploy (ex: **web-portal**)
- **Release**: Name of the service's release to deploy, usually associated with some feature (ex: **accessible-search-form**)
- **Environment**: Environment to deploy the release in (ex: **qa**)
- **CustomParams**: A key/value list of custom parameters 

### Fluxcd

The template directory should contain the following components:
- **app**: Directory that should contain the templatized orchestration files for your service
- **fluxcd.yml**: Templatized orchestration file for the fluxcd resources that will manage your service
- **filesystem-conventions.yml**: Templatized convention file that tells ferlease how to name the files it will generate and where to place them. It contains the following keys:
  - **naming**: What name to give to the **app** directory and to the fluxcd orchestration file.
  - **fluxcd_directory**: Directory where it should place the fluxcd orchestration file. The file will be added to the resources of a pre-existing **kustomization.yaml** file in that directory.
  - **apps_directory**: Directory containing apps under which it will place the **app** directory

### Terraform

The template directory should contain the following components:
  - **entrypoint.tf**: Templatized terraform entrypoint file that contains a template for the top-level terraform file that will be generated.
  - **module**: A directory containing a templatized terraform module that will be generated in the same directory as the entrypoint file. Note that this is optional and for simple cases, just the entrypoint file may be sufficient.
  - **filesystem-conventions.yml**: Templatized convention file that tells ferlease how to name the files it will generate and where to place them. It contains the following keys:
    - **naming**: What name to give to the **entrypoint** file (with the **.tf** suffix appended to it) and to the optional terraform inline module directory.
    - **directory**: Directory where it should place both the entrypoint file and the optional inline module directory.

## Commands

This project has two commands:
- **release**: Performs a release, adding the populated template files to the repo
- **teardown**: Removes a previous release, removing the files from the repo

Additionally, the location of the configuration file can be specified using a **config** argument.

## Configuration File

ferlease expects a configuration file that specifies how it should behave.

Some configuration properties can be templatized (as specified in the properties' description) with the same parameters the template supports (**Service**, **Release**, **Environment** and **CustomParams**) in addition to the following parameters:
- **RepoDir**: Path of the directory where ferlease will have cloned the git repo to operate in, determined at runtime
- **Operation**: Operation that ferlease is performing. Can be either **release** or **teardown**.

The configuration file should be in yaml and is expected to have the following properties:
- **service**: Name of the service to manage
- **release**: Name of the service's release to manage
- **environment**: Environment to manage the release in
- **custom_parameters**: Key/value list of additional custom parameters
- **author**: Author information for the git commit. It has the following keys:
  - **name**: Name of the author
  - **email**: Email of the author
- **commit_message**: Default commit message for all orchestrations. This property can be templatized.
- **push_retries**: In the unlikely even that a barrage of upstream commits keep blocking a gitops operation, how many times to retry before giving up.
- **push_retry_interval**: If a gitop operation is blocked by an upstream commit, how long to wait before re-cloning, re-commiting and re-attempting the push. Should be a string in golang duration format.
- **orchestrations**: List of orchestrations to apply or remove. Each entry contains the following parameters:
  - **type**: Type of the orchestration. Valid values are **fluxcd** or **terraform**.
  - **repo**: Url of the git repo to operate on
  - **ref**: Branch of the git repo to operate on
    - **git_auth**: Git authentication. It has the following keys:
      - **ssh_key**: Ssh key to use to authentify with the git server. It should be the path to a file containing the key, not the key itself.
      - **known_key**: Path to a file containing the git server's ssh fingerprint. Used to authentify the server.
  - **commit_signature**: Path to gpg private key and its passphrase to sign commits. It has the following keys:
    - **key**: Path to a file containing the private key that will sign the commit
    - **passphrase**: Path to a file container the secret passphrase to decrypt the private key
  - **commit_message**: Commit message for this orchestration. If unspecified, the default one will be used. This property can be templatized. 
  - **accepted_signatures**: Path to a directory containing the gpg public keys of all authorized signers for the repo
  - **template_directory**: Path to the directory containing the release's template. This property can be templatized.







