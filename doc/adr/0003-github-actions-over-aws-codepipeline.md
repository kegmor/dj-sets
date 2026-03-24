# 3. GitHub Actions over AWS CodePipeline

Date: 2026-03-20

## Status

Accepted

## Context

We need a CI/CD pipeline to build, test, and deploy. The two realistic options are GitHub Actions and AWS CodePipeline.

## Considerations

While CodePipeline would have demonstrated deeper AWS ecosystem knowledge, GitHub Actions is free for public repositories. The code already exists on GitHub and the workflows live right next to the code. CodePipeline is a separate AWS service that would require additional managing while making viewing the results less seamless.

## Decision

Use GitHub Actions

## Consequences

CodePipeline would have been better to show a deeper understanding of AWS tools and IAM roles are used natively. OIDC federation setup was required to connect to actions and it's an extra layer to configure and maintain. Deployments will stop if OIDC breaks.
