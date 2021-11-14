# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

- Core implementation
  - Create Role, RoleBinding, ServiceAccount
    - [x] Single Namespace
    - [ ] Multi Namespaces
  - [x] Deletion expired resources
  - [x] Generate kubeconfig & stored to secret
- [x] kubeconfig Zip encryption
- [ ] kubeconfig PGP encryption
  - [ ] encryption from github public key
  - [ ] fetch cinderella public key API
- Audit log system
  - [x] tee channel
  - event destination
    - [x] log, io.Writer
    - [ ] k8s event
    - [x] slack event
    - [ ] datadog event
- slack integration
  - App home
    - [ ] Admin user view
      - [ ] Pending claim list
      - [ ] Accept/Reject claim button
      - [ ] audit event
    - [ ] General user view
  - Claim modal view
    - [x] Claim submit
    - [ ] already claimed message
    - [x] radio button default
    - [ ] claim period
  - [x] Send encrypted kubeconfig
  - [x] audit log
- logging
  - structured log
- Persistent data store
  - DB?
