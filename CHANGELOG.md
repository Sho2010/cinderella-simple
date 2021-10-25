# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

- Core implementation
  - Create Role, RoleBinding, ServiceAccount
  - Deletion expired resources
  - Generate kubeconfig & stored to secret
- kubeconfig Zip encryption
- kubeconfig PGP encryption
  - encryption from github public key
  - fetch cinderella public key API
- Audit log system
  - k8s event
- slack integration
  - interactive message
  - slash command
  - audit log

