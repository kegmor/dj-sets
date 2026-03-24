# 4. Go/TypeScript over single-language

Date: 2026-03-20

## Status

Accepted

## Context

A full TypeScript stack is a valid approach. Sharing types between frontend and backend only requiring knee deep knowledge of one ecosystem makes sense. So why not just TypeScript?

## Considerations

The backend will be running on AWS Lambda in a cost-conscious serverless setup. Go compiles to a single binary with significantly faster cold start times than Node.js.
The backend will eventually handle concurrent tasks. Go was specifically designed for concurrency. 

## Decision

Use Go for backend and TypeScript for frontend and infrastructure.

## Consequences

Downside would be needing two testing frameworks, two dependency management systems, and development velocity might be slower initially than having used a single-language stack.
