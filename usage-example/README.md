# About

This is an environment to troubleshoot  ferlease using several sources.

It uses a ferlease binary compiled from the code in this repo and a configuration file also in this repo.

The target should be a fork of this repo which will be both the target of the releases and also contain the templates instructing ferlease what to push on the repo: https://github.com/Ferlab-Ste-Justine/ferlease-playground

Note that for this example, there are 3 sources that are different directories in the same repository, but the sources could also be distinct repositories.

# Requirements

If you run the example directly on your host, you will need golang verson 1.23.0 or later installed.

If you run the example in docker, you will need to have docker installed.

The example expects that you have an ssh key on your host that should have access to a fork of the playground repo in your account at the following path: **~/.ssh/id_rsa**

# Instruction

1. Fork the following repo under your account: https://github.com/Ferlab-Ste-Justine/ferlease-playground

2. Edit either the **config.yml** (if you want to run the example directly on your host) or the **config-docker.yml** (if you want to run the example in docker) file and change the **git@github.com:Ferlab-Ste-Justine/ferlease-playground.git** repo references to your fork of the repo.

3. Run **run_release.sh** and **run_teardown.sh** to do a release and teardown directly from your host. Run **run_release_docker.sh** and **run_teardown_docker.sh** to do a release and teardown in docker.