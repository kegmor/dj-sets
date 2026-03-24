# 5. modular monolith evolving to microservices

Date: 2026-03-20

## Status

Accepted

## Context

Create microservices from the beginning or start with a modular monolith and eventually split off when scaling and dependency profiles call for it?

## Context

A single deployable Go application with clean internal package boundaries that mirror what microservices would be seems the most prudent. Extracting, as the code develops, the services that have a clear reason to be independent shows restraint and not reaching for complexity without justification.

## Decision

To Start modular monolith then switch microservices as needed.

## Consequences

For a single-user application, microservices is a bit overkill. A well-structured monolith demonstrating clear internal boundaries can be just as impressive. In the end it's more risk of time wasted without much reward.
