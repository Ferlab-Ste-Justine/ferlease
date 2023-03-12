# About

This client is meant to add and remove several concurrent versions of the same fluxcd orchestrated kubernetes service, corresponding to various ongoing feature branches, in a gitops manner.

# Usage

## Playground Repo

As you are reading this documentation, refer to the followin resources for a playground that can be experimented with when discovering this project:
- **Playground repo**: https://github.com/Ferlab-Ste-Justine/ferlease-playground
- **Configuration example**: https://github.com/Ferlab-Ste-Justine/ferlease/tree/main/usage-example

You will have to clone the playground repo and adjust the configuration accordingly.

From there, you should be able to play with the **release** and **teardown** commands.

## Assumption

This project assumes that the orchestration of the service to be released will be broken down into the following parts:
- An **apps** directory containing the kubernetes resources of the service to deploy
- An orchestration file containing all fluxcd resources to manage the app (minimally a **fluxcd kustomization** resource pointing to the **apps** directory)
- A **kustomize kustomization** file that manages all the fluxcd resource files of various releases to be deployed

## Templatization Support

This project will orchestrate releases based on a provided template that follows the golang templating format.

Currently, the following variables are supported in the template:
- **Service**: Name of the service to deploy (ex: **web-portal**)
- **Release**: Name of the service's release to deploy, usually associated with some feature (ex: **accessible-search-form**)
- **Environment**: Environment to deploy the release in (ex: **qa**)

The template directory should contain the following components:
- **app**: Directory that should contain the templatized orchestration files for your service
- **fluxcd.yml**: Templatized orchestration file for the fluxcd resources that will manage your service
- **filesystem-conventions.yml**: Templatized convention file that tells ferlease how to name the files it will generate and where to place them. It contains the following keys:
  - **naming**: What name to give to the **app** directory and to the fluxcd orchestration file.
  - **fluxcd_directory**: Directory where it should place the fluxcd orchestration file. The file will be added to the resources of a pre-existing **kustomization.yaml** file in that directory.
  - **apps_directory**: Directory containing apps under which it will place the **app** directory

## Commands

This project has two commands:
- **release**: Performs a release, adding the populated template files to the repo
- **teardown**: Removes a previous release, removing the files from the repo

Additionally, the location of the configuration file can be specified using a **config** argument.

## Configuration File

ferlease expects a configuration file that specifies how it should behave.

Some configuration properties can be templatized (as specified in the properties' description) with the same parameters the template supports (**Service**, **Release**, **Environment**) in addition to the following parameters:
- **RepoDir**: Path of the directory where ferlease will have cloned the git repo to operate in, determined at runtime
- **Operation**: Operation that ferlease is performing. Can be either **release** or **teardown**.

The configuration file should be in yaml and is expected to have the following properties:
- **service**: Name of the service to manage
- **release**: Name of the service's release to manage
- **environment**: Environment to manage the release in
- **repo**: Url of the git repo to operate on
- **ref**: Branch of the git repo to operate on
- **git_ssh_key**: Ssh key to use to authentify with the git server. It should be the path to a file containing the key, not the key itself.
- **git_known_key**: Path to a file containing the git server's ssh fingerprint. Used to authentify the server.
- **template_directory**: Path to the directory containing the release's template. This property can be templatized.
- **commit_message**: Content of the commit message. This property can be templatized.
- **push_retries**: In the unlikely even that a barrage of upstream commits keep blocking a gitops operation, how many times to retry before giving up.
- **push_retry_interval**: If a gitop operation is blocked by an upstream commit, how long to wait before re-cloning, re-commiting and re-attempting the push. Should be a string in golang duration format.