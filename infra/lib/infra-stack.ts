import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins';
import * as s3deploy from 'aws-cdk-lib/aws-s3-deployment';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigwv2 from 'aws-cdk-lib/aws-apigatewayv2';
import { HttpLambdaIntegration } from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import * as s3n from 'aws-cdk-lib/aws-s3-notifications';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';

export class InfraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // 1. S3 バケットの作成

    // S3 バケット削除設定
    const removalPolicy = cdk.RemovalPolicy.DESTROY;
    const autoDeleteObjects = true;

    // Inputバケット（画像アップロード用：1日で自動削除）
    const inputBucket = new s3.Bucket(this, 'InputBucket', {
      removalPolicy,
      autoDeleteObjects,
      lifecycleRules: [{ expiration: cdk.Duration.days(1) }],
      cors: [{
        allowedMethods: [s3.HttpMethods.PUT],
        allowedOrigins: ['*'],
        allowedHeaders: ['*'],
      }],
    });

    // Outputバケット（生成したExcel保存用：1日で自動削除）
    const outputBucket = new s3.Bucket(this, 'OutputBucket', {
      removalPolicy,
      autoDeleteObjects,
      lifecycleRules: [{ expiration: cdk.Duration.days(1) }],
    });

    // Websiteバケット（フロントエンドの静的ファイルホスティング用）
    const websiteBucket = new s3.Bucket(this, 'WebsiteBucket', {
      removalPolicy,
      autoDeleteObjects,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL, 
    });

    // Slack OAuthトークン保存用テーブル
    const slackTokenTable = new dynamodb.Table(this, 'SlackTokenTable', {
      partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
      removalPolicy,
    });

    // 2. CloudFront の作成

    // CloudFront
    const distribution = new cloudfront.Distribution(this, 'WebsiteDistribution', {
      defaultBehavior: {
        origin: origins.S3BucketOrigin.withOriginAccessControl(websiteBucket),
        viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
      },
      defaultRootObject: 'index.html',
    });

    // デプロイ時は CloudFront のキャッシュを最新に更新
    new s3deploy.BucketDeployment(this, 'DeployWebsite', {
      sources: [s3deploy.Source.asset('./test-assets')],
      destinationBucket: websiteBucket,
      distribution: distribution,
      distributionPaths: ['/*'],
    });

    // 3. Lambda 関数の作成（Goランタイム）

    // API Handler（HTTP API）
    const apiHandler = new lambda.Function(this, 'ApiHandler', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      architecture: lambda.Architecture.ARM_64,
      code: lambda.Code.fromAsset('../backend/bin/api'),
      environment: {
        ALLOWED_ORIGINS: '*',
        INPUT_BUCKET_NAME: inputBucket.bucketName,
        SLACK_TOKEN_TABLE_NAME: slackTokenTable.tableName,
      },
    });

    // Processor Handler（Bedrock解析・JSON生成・Slack通知）
    const processorHandler = new lambda.Function(this, 'ProcessorHandler', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      architecture: lambda.Architecture.ARM_64,
      code: lambda.Code.fromAsset('./test-assets/dummy-lambda'),
      timeout: cdk.Duration.seconds(30),
      environment: {
        ALLOWED_ORIGINS: '*',
        OUTPUT_BUCKET_NAME: outputBucket.bucketName,
        OUTPUT_FORMAT: 'json', 
        SLACK_WEBHOOK_URL: process.env.SLACK_WEBHOOK_URL || '',
      },
    });

    // 4. 権限（IAM）と トリガー（Event）の設定

    // API Handlerには、Inputバケットへの「書き込む」権限と Slack Token Tableへの「読み込む」・「書き込む」権限を付与
    inputBucket.grantWrite(apiHandler);
    slackTokenTable.grantReadWriteData(apiHandler);

    // ProcessorHandlerには、Inputから「読み込む」権限、Outputへ「書き込む」権限を付与
    inputBucket.grantRead(processorHandler);
    outputBucket.grantWrite(processorHandler);

    // ProcessorHandlerにTextractの実行権限を付与
    processorHandler.addToRolePolicy(new iam.PolicyStatement({
      actions: ['textract:AnalyzeDocument', 'textract:DetectDocumentText'],
      resources: ['*'],
    }));

    // Inputバケットに画像が入ったら ProcessorHandler を自動起動
    inputBucket.addEventNotification(
      s3.EventType.OBJECT_CREATED,
      new s3n.LambdaDestination(processorHandler)
    );

    // 5. API Gateway の構築 (HTTP API)

    // HTTP API
    const api = new apigwv2.HttpApi(this, 'JpegToJsonHttpApi', {
      apiName: 'Jpeg To Json HTTP API',
      corsPreflight: {
        allowOrigins: ['*'],
        allowMethods: [apigwv2.CorsHttpMethod.ANY],
        allowHeaders: ['*'],
      },
    });

    // POST /presign で API Handler を呼び出す
    const apiIntegration = new HttpLambdaIntegration('ApiIntegration', apiHandler);
    api.addRoutes({
      path: '/presign',
      methods: [apigwv2.HttpMethod.POST],
      integration: apiIntegration,
    });

    // 6. ログでURLを出力
    new cdk.CfnOutput(this, 'CloudFrontURL', { value: `https://${distribution.distributionDomainName}`});
  }
}
