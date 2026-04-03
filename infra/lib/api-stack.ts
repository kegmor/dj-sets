import { Stack, StackProps, Duration } from 'aws-cdk-lib/core';
import { Construct } from 'constructs';
import { Vpc, SecurityGroup, SubnetType } from 'aws-cdk-lib/aws-ec2';
import { ISecret } from 'aws-cdk-lib/aws-secretsmanager';
import { Function, Runtime, Code } from 'aws-cdk-lib/aws-lambda';
import { LambdaRestApi, ApiKeySourceType } from 'aws-cdk-lib/aws-apigateway';
import { PolicyStatement } from 'aws-cdk-lib/aws-iam';

interface ApiStackProps extends StackProps {
  vpc: Vpc;
  dbSecurityGroup: SecurityGroup;
  dbSecret: ISecret;
}

export class ApiStack extends Stack {
  constructor(scope: Construct, id: string, props: ApiStackProps) {
    super(scope, id, props);

     // YouTube Lambda - outside VPC for internet access
    const youtubeFunction = new Function(this, 'YouTubeHandler', {
      runtime: Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      code: Code.fromAsset('../backend/cmd/youtube'),
      timeout: Duration.seconds(15),
    });

     // Grant YouTube Lambda access to the YouTube API key secret
    youtubeFunction.addToRolePolicy(new PolicyStatement({
      actions: ['secretsmanager:GetSecretValue'],
      resources: [
        `arn:aws:secretsmanager:${this.region}:${this.account}:secret:youtube-api-key-*`,
      ],
    }));

    const apiFunction = new Function(this, 'ApiHandler', {
        runtime: Runtime.PROVIDED_AL2023,
        handler: 'bootstrap',
        code: Code.fromAsset('../backend/cmd/api'),
        vpc: props.vpc,
        vpcSubnets: {
            subnetType: SubnetType.PRIVATE_ISOLATED,
        },
        securityGroups: [props.dbSecurityGroup],
        timeout: Duration.seconds(30),
        environment: {
            YOUTUBE_LAMBDA: youtubeFunction.functionName,
        },
    });

    const migrateFunction = new Function(this, 'MigrateHandler', {
    runtime: Runtime.PROVIDED_AL2023,
    handler: 'bootstrap',
    code: Code.fromAsset('../backend/cmd/migrate'),
    vpc: props.vpc,
    vpcSubnets: {
        subnetType: SubnetType.PRIVATE_ISOLATED,
    },
    securityGroups: [props.dbSecurityGroup],
    timeout: Duration.seconds(30),
    });

    migrateFunction.addToRolePolicy(new PolicyStatement({
        actions: ['secretsmanager:GetSecretValue'],
        resources: [
            props.dbSecret.secretArn,
        ],
    }));

    // Grant API Lambda access to invoke YouTube Lambda
    youtubeFunction.grantInvoke(apiFunction);

    apiFunction.addToRolePolicy(new PolicyStatement({
        actions: ['secretsmanager:GetSecretValue'],
        resources: [
            props.dbSecret.secretArn,
        ],
    }));

    const api = new LambdaRestApi(this, 'DjSetsApi', {
    handler: apiFunction,
    apiKeySourceType: ApiKeySourceType.HEADER,
    defaultCorsPreflightOptions: {
        allowOrigins: ['*'],
        allowMethods: ['GET', 'POST', 'DELETE', 'OPTIONS'],
        allowHeaders: ['Content-Type', 'x-api-key'],
    },
    defaultMethodOptions: {
        apiKeyRequired: true,
    },
});

    const apiKey = api.addApiKey('DjSetsApiKey');
    const plan = api.addUsagePlan('UsagePlan', {
        throttle: {
            rateLimit: 10,
            burstLimit: 5,
        },
    });
    
    plan.addApiKey(apiKey);
    plan.addApiStage({ stage: api.deploymentStage });
  }
}