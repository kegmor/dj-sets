import {
  CfnOutput,
  Duration,
  RemovalPolicy,
  Stack, 
  StackProps,  
 } from 'aws-cdk-lib/core';
import {
  InstanceType, 
  InstanceClass, 
  InstanceSize,
  Peer,
  Port,
  SubnetType, 
  SecurityGroup, 
  Vpc, 
  InterfaceVpcEndpointAwsService 
} from 'aws-cdk-lib/aws-ec2'
import {
  Credentials,
  DatabaseInstance, 
  DatabaseInstanceEngine, 
  PostgresEngineVersion,
  StorageType, 
} from 'aws-cdk-lib/aws-rds'
import { Construct } from 'constructs';
import { ISecret } from 'aws-cdk-lib/aws-secretsmanager';

export class DatabaseStack extends Stack {
  public readonly vpc: Vpc;
  public readonly dbSecurityGroup: SecurityGroup;
  public readonly dbSecret: ISecret;
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const vpc = new Vpc(this, 'DatabaseVpc', {
      maxAzs: 2,
      natGateways: 0, //small project doesn't need outbound internet access
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: 'Public',
          subnetType: SubnetType.PUBLIC,
        },
        {
          cidrMask: 24,
          name: 'Isolated',
          subnetType: SubnetType.PRIVATE_ISOLATED,
        },
      ]
    });

    vpc.addInterfaceEndpoint('SecretsManagerEndpoint', {
      service: InterfaceVpcEndpointAwsService.SECRETS_MANAGER,
    });

    vpc.addInterfaceEndpoint('LambdaEndpoint', {
      service: InterfaceVpcEndpointAwsService.LAMBDA,
    });
    
    const dbSecurityGroup = new SecurityGroup(this, 'DbSecurityGroup', {
      vpc,
      description: "Security group for RDS PostgreSQL",
      allowAllOutbound: false,
    })

    dbSecurityGroup.addIngressRule(
      Peer.ipv4(vpc.vpcCidrBlock),
      Port.tcp(5432),
    )

    dbSecurityGroup.addEgressRule(
      Peer.ipv4(vpc.vpcCidrBlock),
      Port.tcp(443),
    );

    dbSecurityGroup.addEgressRule(
      Peer.ipv4(vpc.vpcCidrBlock),
      Port.tcp(5432),
    );

    // Create the RDS PostgreSQL instance
    const database = new DatabaseInstance(this, 'PostgresDb', {
      engine: DatabaseInstanceEngine.postgres({
        version: PostgresEngineVersion.VER_16_12,
      }),
      instanceType: InstanceType.of(
        InstanceClass.T3,
        InstanceSize.MICRO,
      ),
      vpc,
      vpcSubnets: {
        subnetType: SubnetType.PRIVATE_ISOLATED,
      },
      securityGroups: [dbSecurityGroup],

      databaseName: 'djsetdb',
      port: 5432,
      credentials: Credentials.fromGeneratedSecret('dbadmin', {
        secretName: 'rds-credentials',
      }),

      // Storage settings
      allocatedStorage: 20, // Start with 20 GB only storing metadata
      storageType: StorageType.GP2,
      storageEncrypted: true,

      // Availability and durability
      multiAz: false,
      deletionProtection: false,

      // Backup configuration
      backupRetention: Duration.days(7),
      preferredBackupWindow: '3:00-4:00', // 3 AM UTC
      preferredMaintenanceWindow: 'sun:04:00-sun:05:00',

      // Removal behavior
      removalPolicy: RemovalPolicy.DESTROY,
    });

    this.vpc = vpc;
    this.dbSecurityGroup = dbSecurityGroup;
    this.dbSecret = database.secret!;
  }
}
