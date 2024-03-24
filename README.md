# Amazon Kinesis Client Library for Go
[![Integration test](https://github.com/arthurbailao/aws-kcl/actions/workflows/integration_test.yml/badge.svg?branch=main)](https://github.com/arthurbailao/aws-kcl/actions/workflows/integration_test.yml)

This package provides an interface to the [Amazon Kinesis Client Library][amazon-kcl] (KCL) [MultiLangDaemon][multi-lang-daemon] for Golang.

Developers can use the KCL to build distributed applications that process streaming data reliably at scale. The KCL takes care of many of the complex tasks associated with distributed computing, such as load-balancing across multiple instances, responding to instance failures, checkpointing processed records and reacting to changes in stream volume.

This package wraps and manages the interaction with the [MultiLangDaemon][multi-lang-daemon], which is provided as part of the [Amazon KCL for Java][amazon-kcl-github] so that developers can focus on implementing their record processing logic.

A record processor in Go must implement the [RecordProcessor][record-processor-interface] interface and call the function [Run][function-run].

[amazon-kcl]: https://docs.aws.amazon.com/streams/latest/dev/developing-consumers-with-kcl-v2.html
[multi-lang-daemon]: https://github.com/awslabs/amazon-kinesis-client/blob/master/amazon-kinesis-client-multilang/src/main/java/software/amazon/kinesis/multilang/package-info.java
[amazon-kcl-github]: https://github.com/awslabs/amazon-kinesis-client
[record-processor-interface]: https://godoc.org/github.com/arthurbailao/aws-kcl#RecordProcessor
[function-run]: https://godoc.org/github.com/arthurbailao/aws-kcl#Run
