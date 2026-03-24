#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib/core';
import { DatabaseStack } from '../lib/database-stack';
import { ApiStack } from '../lib/api-stack';

const app = new cdk.App();
const databaseStack = new DatabaseStack(app, 'DatabaseStack', {
  env: { account: process.env.CDK_DEFAULT_ACCOUNT, region: process.env.CDK_DEFAULT_REGION },
});

new ApiStack(app, 'ApiStack', {
  env: { account: process.env.CDK_DEFAULT_ACCOUNT, region: process.env.CDK_DEFAULT_REGION },
  vpc: databaseStack.vpc,
  dbSecurityGroup: databaseStack.dbSecurityGroup,
  dbSecret: databaseStack.dbSecret,
});
