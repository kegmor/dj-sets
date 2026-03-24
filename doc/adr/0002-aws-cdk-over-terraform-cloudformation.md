# 2. AWS CDK over Terraform/CloudFormation

Date: 2026-03-20

## Status

Accepted

## Context

The project is exclusively AWS. The purpose of the portfolio project is to demonstrate knowledge in the AWS ecosystem. The three options were AWS CloudFormation, AWS CDK, and Terraform.

## Considerations

Because AWS CDK generates AWS CloudFormation under the hood and TypeScript is what will be used for the frontend, CDK being TypeScript, makes more sense to use since it will allow to both write the infrastructure and the front end code in one language.

Terraform uses its own language called HCL and it would be an additional burden to learn another language for such a small project.

## Decision

To use AWS CDK

## Consequences

If there eventually becomes a need to be cloud agnostic choosing CDK won't directly transfer to Terraform.
